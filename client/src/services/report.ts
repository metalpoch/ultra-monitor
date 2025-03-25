import { ErrorHandler } from "../lib/errors";

/**
 * @class Handler of all authentication requests for the API.
 */
export class ReportService {
    private static url: string = import.meta.env.PUBLIC_API_REPORT;


    /**
     * Request API for the upload of a file.
     * 
     * @param {file} file File to be uploaded.
     * @param {id} id ID of the file.
     * @param {category} category Category of the file.
     */
    static async uploadFile(file: any, id: number, category: string = "Reporte Regional"): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            let data = new FormData();
            data.append("file", file);
            data.append("user_id", id.toString());
            data.append("category", category);

            const response = await fetch(`${this.url}/report`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: data
            });
            if (response.status === 200) return { status: response.status, info: null }
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) }
        } catch (err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) }
        }
    }
}