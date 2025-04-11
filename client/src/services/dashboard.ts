import { ErrorHandler } from "../lib/errors";


export class DashboardService {
    private static url: string = import.meta.env.PUBLIC_API_CORE;


    /**
     * Get traffic to state to dashboard.
     * 
     * @param {string} month Month to get the traffic.
     */
    static async getTrafficByState(month: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/state/${month}`, {
                method: 'GET',
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: 500, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    /**
     * Get traffic top N to state to dashboard.
     * 
     * @param {string} month Month to get the traffic.
     * @param {number} n Number of top states.
     */
    static async getTrafficByStateTopN(month: string, n: number): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/state_n/${month}/${n}`,{
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    /**
     * Get traffic to ODN to dashboard.
     * 
     * @param {string} month Month to get the traffic.
     */
    static async getTrafficByOdn(month: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/odn_d/${month}`,{
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

}