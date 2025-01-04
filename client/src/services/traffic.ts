import { ErrorHandler } from "../lib/errors";

export class TrafficService {
    private static url: string = import.meta.env.PUBLIC_API_CORE;

    static async getDevice(deviceID: number, initialDate: string, endDate: string, initialTime: string, endTime: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/device/${deviceID}?init_date=${initialDate}T${initialTime}:00Z&end_date=${endDate}T${endTime}:59Z`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    static async getInterface(interfaceID: number, initialDate: string, endDate: string, initialTime: string, endTime: string) {
        try {
            const response = await fetch(`${this.url}/traffic/interface/${interfaceID}?init_date=${initialDate}T${initialTime}:00Z&end_date=${endDate}T${endTime}:59Z`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    static async getState(state: string, initialDate: string, endDate: string, initialTime: string, endTime: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/location/${state}?init_date=${initialDate}T${initialTime}:00Z&end_date=${endDate}T${endTime}:59Z`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    static async getCounty(state: string, county: string, initialDate: string, endDate: string, initialTime: string, endTime: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/location/${state}/${county}?init_date=${initialDate}T${initialTime}:00Z&end_date=${endDate}T${endTime}:59Z`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    static async getMunicipality(state: string, county: string, municipality: string, initialDate: string, endDate: string, initialTime: string, endTime: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/location/${state}/${county}/${municipality}?init_date=${initialDate}T${initialTime}:00Z&end_date=${endDate}T${endTime}:59Z`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    static async getOdn(odn: string, initialDate: string, endDate: string, initialTime: string, endTime: string): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/traffic/odn/${odn}?init_date=${initialDate}T${initialTime}:00Z&end_date=${endDate}T${endTime}:59Z`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    
}