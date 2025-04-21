import { ErrorHandler } from "../lib/errors";

/**
 * @class Handler of all inteface requests for the API.
 */
export class InterfaceService {
    private static url: string = process.env.PUBLIC_API_CORE ?? import.meta.env.PUBLIC_API_CORE;

    /**
     * Request API to get all interfaces by a shell of an OLT device.
     * 
     * @param {number} deviceID OLT device ID.
     * @param {number} shell OLT device shell.
     */
    static async getInterfacesByShell(deviceID: number, shell: number = 0): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
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


    /**
     * Request API to get all interfaces by card of an OLT device.
     * 
     * @param {number} deviceID OLT device ID.
     * @param {number} card OLT device card.
     * @param {number} shell OLT device shell.
     */
    static async getInterfacesByCard(deviceID: number, card: number, shell: number = 0): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
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


    /**
     * Request API to get all interfaces by port of an OLT device.
     * 
     * @param {number} deviceID OLT device ID.
     * @param {number} card OLT device card.
     * @param {number} port OLT device port.
     * @param {number} shell OLT device shell.
     */
    static async getInterfacesByPort(deviceID: number, card: number, port: number, shell: number = 0): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
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


    /**
     * Request API to get all information about an interface by its ID.
     * 
     * @param {number} id Interface ID.
     */
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
