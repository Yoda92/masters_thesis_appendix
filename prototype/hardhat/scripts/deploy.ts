import { ethers } from "hardhat";
import * as fs from "fs";
import { EnergyCommunity } from "../typechain-types";
import { TestUtility } from "./test.utility";

function getRunnerAddress(energyCommunityContract: EnergyCommunity): string {
  if (
    energyCommunityContract.runner &&
    "address" in energyCommunityContract.runner
  ) {
    return String(energyCommunityContract.runner.address);
  }

  return "Unknown address";
}

async function main() {
  const signers = await ethers.getSigners();
  const smartMeters = signers.slice(3);
  const testAddress = signers[0];
  const energyCommunityContract = await ethers.deployContract("EnergyCommunity", [
    testAddress.address,
    testAddress.address,
    testAddress.address,
  ]);
  await energyCommunityContract.waitForDeployment();

  for (const smartMeter of smartMeters) {
    // For testing purposes
    await TestUtility.approveSmartMeter(energyCommunityContract, smartMeter.address);
  }

  console.log(`EnergyCommunity deployed to ${energyCommunityContract.target}`);
  const output = {
    contractAddress: energyCommunityContract.target,
    contractCreator: getRunnerAddress(energyCommunityContract),
  };
  fs.writeFileSync(
    __dirname + "/out/latest_deploy_address.json",
    JSON.stringify(output)
  );
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
