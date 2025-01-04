import { ErrorHandler } from "../lib/errors";

export class InterfaceService {
    private static url: string = import.meta.env.PUBLIC_API_CORE;

    static async getInterfacesByCard(deviceID: number, shell: number = 0): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/interface/device/${deviceID}/find?shell=${shell}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    static async getInterfacesByPort(deviceID: number, card: number, shell: number = 0): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/interface/device/${deviceID}/find?shell=${shell}&card=${card}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    static async getInterfaces(deviceID: number, card: number, port: number, shell: number = 0): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/interface/device/${deviceID}/find?shell=${shell}&card=${card}&port=${port}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }

    static async getInterfaceByID(id: number): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/info/interface/${id}`, {
                method: 'GET'
            });
            if (response.ok) return { status: response.status, info: await response.json() };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }
}