import { EnergyCommunity } from "../typechain-types";
import { DockerUtility } from "./docker.utility";

export class TestUtility {
  static async registerSmartMeter(energyCommunityContract: EnergyCommunity) {
    const smartMeterAddress = TestUtility.getRunnerAddress(
      energyCommunityContract
    );
    try {
      await energyCommunityContract.registerSmartMeter();
      console.log(
        "Smart meter with address:" + smartMeterAddress + " registered."
      );
    } catch (error) {
      console.log(
        "Smart meter with address:" + smartMeterAddress + " already exists."
      );
    }
  }

  static async approveSmartMeter(energyCommunityContract: EnergyCommunity, smartMeterAddress: string) {
    await energyCommunityContract.approveSmartMeter(smartMeterAddress);
    console.log(
      "Smart meter with address:" + smartMeterAddress + " approved."
    );
  }

  static getRunnerAddress(energyCommunityContract: EnergyCommunity): string {
    if (
      energyCommunityContract.runner &&
      "address" in energyCommunityContract.runner
    ) {
      return String(energyCommunityContract.runner.address);
    }

    return "Unknown address";
  }

  static async measureExecutionTime<T>(
    promise: Promise<T>
  ): Promise<[T, number]> {
    const before = Date.now();
    const result = await promise;
    const after = Date.now();

    const executionTime = after - before;

    return [result, executionTime];
  }

  static async delay(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  static addPeriodicCPUUsageListener(
    containerStatistics: NodeJS.ReadableStream,
    results: Array<number>
  ) {
    containerStatistics.on("data", (data) => {
      const result = DockerUtility.calculateContainerCpuUsage(
        JSON.parse(String(data))
      );

      if (isNaN(result)) {
        return;
      }

      results.push(result);
    });
  }
}
