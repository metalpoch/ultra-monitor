import { ErrorHandler } from "../lib/errors";

export class DeviceService {
    private static url: string = import.meta.env.PUBLIC_API_CORE;

    static async getDevices(): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/device`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

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
