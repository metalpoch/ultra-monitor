import type { TrendSchema} from "../schemas/trend";
import { ErrorHandler } from "../lib/errors";

/**
 * @class Handler of the IA requests for the API.
 */
export class TrendService {
    private static url: string = process.env.PUBLIC_API_SMART || import.meta.env.PUBLIC_API_SMART;

    /**
     * Request API to get trend data of an OLT.
     * 
     * @param {string} sysname OLT acronym.
     * @param {number} month Month to consult.
     */
    static async getTrend(sysname: string, month: string): Promise<{ status: (number | null), data: (TrendSchema | null), err: (ErrorHandler | null) }> {
        try {
            const response = await fetch(`${this.url}/trend?sysname=${sysname}&future_month=${month}`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                }
            });
            if (response.ok) return { status: response.status, data: await response.json() as TrendSchema, err: null };
            else return { status: response.status, data: null, err: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, data: null, err: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }
}
