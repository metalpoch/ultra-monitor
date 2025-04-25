import { ReportService } from "../services/report";

export class ReportController {
    /**
     * Upload a file to the report.
     * 
     * @param {file} file File to be uploaded.
     * @param {id} id User ID.
     * @param {category} category Category of the file.
     * @returns {boolean} True if the file was uploaded successfully, false otherwise.
     */
    static async uploadFile(file: any, id: number, category: string = "Reporte Regional"): Promise<boolean> {
        const response = await ReportService.uploadFile(file, id, category);
        if ((response.status === 201) || (response.status === 200)) {
            return true;
        } else {
            console.error(response.info!.message);
            return false;
        }
    }
}