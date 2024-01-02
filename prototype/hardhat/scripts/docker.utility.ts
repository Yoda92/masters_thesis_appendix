import Docker from "dockerode";

export class DockerUtility {
  private static docker: Docker = new Docker();

  private constructor() {}

  static calculateContainerCpuUsage(statistics: any): number {
    const totalUsage = statistics.cpu_stats.cpu_usage.total_usage;
    const totalUsagePrevious = statistics.precpu_stats.cpu_usage.total_usage;
    const totalUsageDelta = totalUsage - totalUsagePrevious;

    const systemUsage = statistics.cpu_stats.system_cpu_usage;
    const systemUsagePrevious = statistics.precpu_stats.system_cpu_usage;
    const systemUsageDelta = systemUsage - systemUsagePrevious;

    const cpuUsagePercentage =
      (totalUsageDelta / systemUsageDelta) *
      statistics.cpu_stats.online_cpus *
      100;

    return cpuUsagePercentage;
  }

  static async getContainerStatistics(
    containerId: string
  ): Promise<NodeJS.ReadableStream> {
    const container = this.docker.getContainer(containerId);
    const statsStream = await container.stats({ stream: true });
    return statsStream;
  }
}
