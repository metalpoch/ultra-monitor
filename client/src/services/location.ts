import { ErrorHandler } from "../lib/errors";

/**
 * @class Handler of all location requests for the API.
 */
export class LocationService {
    private static url: string = import.meta.env.PUBLIC_API_CORE;

    /**
     * Request API to get all states avaliable.
     */
    static async getStates(): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/location/state`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }


    /**
     * Request API to get all counties avaliable by a state.
     * 
     * @param {string} state State to be searched.
     */
    static async getCounties(state: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/location/${state}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }


    /**
     * Request API to get all municipalities avaliable by a state and a county.
     * 
     * @param {string} state State to be searched.
     * @param {string} county County to be searched.
     */
    static async getMunicipalities(state: string, county: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/location/${state}/${county}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }
}
