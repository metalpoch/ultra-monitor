import { ErrorHandler } from "../lib/errors";

/**
 * @class Handler of all traffic/measurament requests for the API.
 */
export class TrafficService {
    private static url: string = import.meta.env.PUBLIC_API_CORE;

    /**
     * Request API to get traffic of an OLT device.
     * 
     * @param {number} deviceID OLT device ID.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     */
    static async getDevice(deviceID: number, initialDate: string, endDate: string, initialTime: string, endTime: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/device/${deviceID}?init_date=${initialDate}T${initialTime}:00-04:00&end_date=${endDate}T${endTime}:59-04:00`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }


    /**
     * Request API to get traffic of an interface.
     * 
     * @param {number} interfaceID Interface ID.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     */
    static async getInterface(interfaceID: number, initialDate: string, endDate: string, initialTime: string, endTime: string) {
        try {
            const response = await fetch(`${this.url}/traffic/interface/${interfaceID}?init_date=${initialDate}T${initialTime}:00-04:00&end_date=${endDate}T${endTime}:59-04:00`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }


    /**
     * Request API to get traffic of a state.
     * 
     * @param {string} state Name of the state.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     */
    static async getState(state: string, initialDate: string, endDate: string, initialTime: string, endTime: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/location/${state}?init_date=${initialDate}T${initialTime}:00-04:00&end_date=${endDate}T${endTime}:59-04:00`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    
    /**
     * Request API to get traffic of a county.
     * 
     * @param {string} state Name of the state.
     * @param {string} county Name of the county.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     */
    static async getCounty(state: string, county: string, initialDate: string, endDate: string, initialTime: string, endTime: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/location/${state}/${county}?init_date=${initialDate}T${initialTime}:00-04:00&end_date=${endDate}T${endTime}:59-04:00`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }


    /**
     * Request API to get traffic of a municipality.
     * 
     * @param {string} state Name of the state.
     * @param {string} county Name of the county.
     * @param {string} municipality Name of the municipality.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     */
    static async getMunicipality(state: string, county: string, municipality: string, initialDate: string, endDate: string, initialTime: string, endTime: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/location/${state}/${county}/${municipality}?init_date=${initialDate}T${initialTime}:00-04:00&end_date=${endDate}T${endTime}:59-04:00`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }


    /**
     * Request API to get traffic of an ODN.
     * 
     * @param {string} state Name of the ODN.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     */
    static async getOdn(odn: string, initialDate: string, endDate: string, initialTime: string, endTime: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/odn/${odn}?init_date=${initialDate}T${initialTime}:00-04:00&end_date=${endDate}T${endTime}:59-04:00`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

}