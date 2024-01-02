import * as fs from "fs";

export class CSVUtility {
  static toCSVFormat(data: unknown[][]): string {
    return data.map((row) => row.join(",")).join("\n");
  }

  static toCSVFile(filePath: string, data: any[][]) {
    const csvContent = CSVUtility.toCSVFormat(data);

    fs.writeFileSync(filePath, csvContent, "utf-8");
  }
}
