import { ErrorHandler } from "../lib/errors";

/**
 * @class Handler of all OLT device requests for the API.
 */
export class DeviceService {
    private static url: string = import.meta.env.PUBLIC_API_CORE;


    /**
     * Request API to get all OLT devices.
     */
    static async getDevices(): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            console.log(this.url);
            const response = await fetch(`${this.url}/info/device`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }


    /**
     * Request API to get all OLT devices by a state.
     * 
     * @param {string} state State to be searched.
     */
    static async getDevicesByState(state: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/device/location/${state}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }


    /**
     * Request API to get all OLT devices by a county of a state.
     * 
     * @param {string} state State to be searched.
     * @param {string} county County to be searched.
     */
    static async getDevicesByCounty(state: string, county: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/device/location/${state}/${county}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }


    /**
     * Request API to get all OLT devices by municipality of a county and state.
     * 
     * @param {string} state State to be searched.
     * @param {string} county County to be searched.
     * @param {string} municipality Municipality to be searched.
     */
    static async getDevicesByMunicipality(state: string, county: string, municipality: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {  
        try {
            const response = await fetch(`${this.url}/info/device/location/${state}/${county}/${municipality}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }


    /**
     * Request API to get information about an OLT device by its ID.
     * 
     * @param {number} idDevice OLD device ID.
     */
    static async getDeviceByID(idDevice: number): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/device/${idDevice}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }


    /**
     * Request API to get information about an OLT device by its sysname.
     * 
     * @param {string} sysname OLD device sysname.
     */
    static async getDeviceBySysname(sysname: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/device/sysname/${sysname}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

}
