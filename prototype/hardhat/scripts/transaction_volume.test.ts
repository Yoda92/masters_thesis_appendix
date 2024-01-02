import { ethers } from "hardhat";
import latestDeployAddress from "./out/latest_deploy_address.json";
import { TestUtility } from "./test.utility";
import { CSVUtility } from "./csv.utility";

const MIN_SMART_METER_PUBLISH_PR_HOUR = 1;
const MAX_SMART_METER_PUBLISH_PR_HOUR = 60;
const HOURS_PR_DAY = 24;
const COMMUNITY_SIZE = 33;
// 1ETH = 2019.32 USD
// https://coinmarketcap.com/
// 22-11-2023
const CONVERSION_RATE = 2019;
const CONVERSION_FACTOR = 1000000000;
const SETTLEMENT_PERIOD = 104;

const HEADERS = [
  "tph",
  "transactions",
  "gas_usage",
  "estimated_ethereum_price",
];

async function runTransactionSpreadTest() {
  const accounts = await ethers.getSigners();
  const smartMeters = accounts.slice(3);
  let energyCommunityContract = await ethers.getContractAt(
    "EnergyCommunity",
    latestDeployAddress.contractAddress
  );

  energyCommunityContract.addListener(
    energyCommunityContract.getEvent("SettlementCompleted"),
    (event, ...args) => {
      console.log(args);
    }
  );

  // Each smart meter publishes data
  var publishSmartMeterPeriodicDataResponse = null;
  for (const smartMeter of smartMeters) {
    const connection = energyCommunityContract.connect(smartMeter);
    await TestUtility.registerSmartMeter(connection);
    publishSmartMeterPeriodicDataResponse = await connection
      .publishSmartMeterPeriodicData(1, 0, SETTLEMENT_PERIOD)
      .then((response) => response.wait());
    if (!publishSmartMeterPeriodicDataResponse) {
      throw new Error("Unconfirmed transaction.");
    }
  }
  if (!publishSmartMeterPeriodicDataResponse) {
    throw new Error("Unconfirmed transaction.");
  }

  // Aggregator publishes BESS data
  const connection = energyCommunityContract.connect(accounts[0]);

  const publishBESSPeriodicDataResponse = await connection
    .publishBESSPeriodicData(0, 2, 1, SETTLEMENT_PERIOD)
    .then((response) => response.wait());

  if (!publishBESSPeriodicDataResponse) {
    throw new Error("Unconfirmed transaction.");
  }

  console.log(
    "BESS periodic data GAS usage: " + publishBESSPeriodicDataResponse.gasUsed
  );

  // Electricity provider triggers settlement
  const settlementResponse = await connection
    .performSettlement(SETTLEMENT_PERIOD)
    .then((response) => response.wait());

  if (!settlementResponse) {
    throw new Error("Unconfirmed transaction.");
  }

  console.log("Settlement GAS usage: " + settlementResponse.gasUsed);

  const PublishPrHourList = Array.from(
    {
      length:
        MAX_SMART_METER_PUBLISH_PR_HOUR - MIN_SMART_METER_PUBLISH_PR_HOUR + 1,
    },
    (v, k) => k + 1
  );

  const rows = new Array();
  rows.push(HEADERS);
  rows.push([0, 0, 0, 0]);

  for (const publishPrHour of PublishPrHourList) {
    const totalSmartMeterPublishPrHour = COMMUNITY_SIZE * publishPrHour;
    const totalBESSPeriodicDataPublishPrHour = 1 * publishPrHour;
    const totalSettlementsPrHour = 1 * publishPrHour;
    const totalTransactionsPrHour =
      totalSmartMeterPublishPrHour +
      totalBESSPeriodicDataPublishPrHour +
      totalSettlementsPrHour;
    const totalGasUsage =
      BigInt(publishSmartMeterPeriodicDataResponse.gasUsed) *
        BigInt(totalSmartMeterPublishPrHour) +
      BigInt(publishBESSPeriodicDataResponse.gasUsed) *
        BigInt(totalBESSPeriodicDataPublishPrHour) +
      BigInt(settlementResponse.gasUsed) * BigInt(totalSettlementsPrHour);
    const totalGasUagePrice =
      (BigInt(CONVERSION_RATE) * totalGasUsage) / BigInt(CONVERSION_FACTOR);
    console.log(
      publishPrHour +
        ", GAS: " +
        totalGasUsage +
        ", Price: " +
        totalGasUagePrice +
        ", Transactions: " +
        totalTransactionsPrHour
    );

    rows.push([
      publishPrHour,
      totalTransactionsPrHour * HOURS_PR_DAY,
      totalGasUsage * BigInt(HOURS_PR_DAY),
      totalGasUagePrice * BigInt(HOURS_PR_DAY),
    ]);

    CSVUtility.toCSVFile(__dirname + "/out/transaction_volume.csv", rows);
  }
}

async function main() {
  await runTransactionSpreadTest();
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
