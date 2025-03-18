import { ErrorHandler } from "../lib/errors";

/**
 * @class Handler of all ODN requests for the API.
 */
export class OdnService {
    private static url: string = import.meta.env.PUBLIC_API_CORE;

    /**
     * Request API to get all ODN by an OLT device ID.
     * 
     * @param {number} id OLT device ID.
     */
    static async getOdnByDevice(id: number): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/odn/device/${id}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

}
