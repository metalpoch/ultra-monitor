import { ErrorHandler } from "../lib/errors";

export class OdnService {
    private static url: string = import.meta.env.PUBLIC_API_CORE;

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
