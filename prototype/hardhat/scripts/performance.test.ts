import { ethers } from "hardhat";
import latestDeployAddress from "../scripts/out/latest_deploy_address.json";
import { DockerUtility } from "./docker.utility";
import { TestUtility } from "./test.utility";
import { CSVUtility } from "./csv.utility";

const MIN_TPS = 1;
const MAX_TPS = 60;
const TPS_TEST_LENGTH = 5;
const CONTAINER_ID = "fec195b506c2";
const COLUMNS = [
  "tps",
  "mean_confirmation_time",
  "max_confirmation_time",
  "cpu_consumption",
];
const SETTLEMENT_PERIOD = 102;

function mean(values: Array<number>): number {
  const sum = values.reduce((sum, current) => sum + current);
  return sum / values.length;
}

async function runPerformanceTest() {
  const accounts = await ethers.getSigners();
  const smartMeters = accounts.slice(3);
  
  let energyCommunityContract = await ethers.getContractAt(
    "EnergyCommunity",
    latestDeployAddress.contractAddress
  );

  for (const smartMeter of smartMeters) {
    const connection = energyCommunityContract.connect(smartMeter);
    await TestUtility.registerSmartMeter(connection);
  }

  const rows: Array<Array<string>> = new Array();
  rows.push(COLUMNS);

  const TPSList = Array.from(
    { length: MAX_TPS - MIN_TPS + 1 },
    (v, k) => k + 1
  );

  let smartMeterIndex = 0;
  const containerStatistics = await DockerUtility.getContainerStatistics(
    CONTAINER_ID
  );

  await TestUtility.delay(2000);

  const baselineContainerCPUMeasurements = new Array();

  TestUtility.addPeriodicCPUUsageListener(
    containerStatistics,
    baselineContainerCPUMeasurements
  );
  console.log("Testing baseline");
  for (let i = 0; i < TPS_TEST_LENGTH; i++) {
    await TestUtility.delay(1000);
  }
  containerStatistics.removeAllListeners();
  const meanBaselineContainerCPUMeasurements = mean(
    baselineContainerCPUMeasurements
  );
  rows.push(["0", "0", "0", String(meanBaselineContainerCPUMeasurements)]);
  console.log(
    "CPU Measurements average: " + meanBaselineContainerCPUMeasurements + "%"
  );

  for (const tps of TPSList) {
    const TPSTestResults = [];
    const TPSContainerCPUMeasurements: Array<number> = [];
    TestUtility.addPeriodicCPUUsageListener(
      containerStatistics,
      TPSContainerCPUMeasurements
    );
    console.log("Testing with TPS: " + tps);
    for (let y = 0; y < TPS_TEST_LENGTH; y++) {
      for (let i = 0; i < tps; i++) {
        const index = smartMeterIndex % smartMeters.length;
        const connection = energyCommunityContract.connect(smartMeters[index]);
        const response = TestUtility.measureExecutionTime(
          connection
            .publishSmartMeterPeriodicData(1, 1, SETTLEMENT_PERIOD)
            .then((result) => result.wait())
        );
        TPSTestResults.push(response);

        smartMeterIndex++;
        await TestUtility.delay(10);
      }
      await TestUtility.delay(1000 - (10 * tps));

    }
    console.log("Finished test with TPS: " + tps);
    const finished_test_result = await Promise.all(TPSTestResults);
    const finished_test_result_execution_times = finished_test_result.map(
      (result) => result[1]
    );
    const execution_times_average = mean(finished_test_result_execution_times);
    const execution_times_max = Math.max(
      ...finished_test_result_execution_times
    );
    const cpu_measurements_average = mean(TPSContainerCPUMeasurements);
    console.log("CPU Measurements average: " + cpu_measurements_average + "%");
    console.log("Average: " + execution_times_average + "ms");
    console.log("Max: " + execution_times_max + "ms");
    rows.push([
      String(tps),
      String(execution_times_average),
      String(execution_times_max),
      String(cpu_measurements_average),
    ]);
    containerStatistics.removeAllListeners();
  }

  CSVUtility.toCSVFile(__dirname + "/out/performance_test.csv", rows);
}

async function main() {
  await runPerformanceTest();
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
