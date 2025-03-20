import { DeviceService } from "../services/device";
import { InterfaceService } from "../services/interface";
import { Strings } from "../constant/strings";
import type { InterfaceSchema } from "../schemas/interface";
import type { DeviceSchema } from "../schemas/device";

/**
 * @class Controller for all requests to the OLT device.
 */
export class DeviceController {

    /**
     * Get all OLT device. Bring only basic OLT status information (ip, community, sysname, syslocation, is_alive, etc...).
     * 
     * @returns {DeviceSchema[]} List of OLT devices available.
     */
    static async getAllDevices(): Promise<DeviceSchema[]> {
        const response = await DeviceService.getDevices();
        if (response.status === 200) {
            return response.info as DeviceSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get all OLT device sysname/acronyms.
     * 
     * @returns {string[]} List of all acronyms available.
     */
    static async getAllAcronyms(): Promise<string[]> {
        const response = await DeviceService.getDevices();
        if (response.status === 200) {
            let deviceNames: string[] = [];
            response.info.map((device: DeviceSchema) => {
                deviceNames.push(device.sysname);
            });
            deviceNames = deviceNames.sort();
            return deviceNames;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get all OLT devices by a state.
     * 
     * @param {string} state Name the state.
     * @returns {DeviceSchema[]} List of OLT devices available.
     */
    static async getAllDevicesByState(state: string): Promise<DeviceSchema[]> {
        const response = await DeviceService.getDevicesByState(state);
        if (response.status === 200) {
            return response.info as DeviceSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get all OLT devices by a county and a state.
     * 
     * @param {string} state Name of the state.
     * @param {string} county Name of the county.
     * @returns {DeviceSchema[]} List of OLT devices available.
     */
    static async getAllDevicesByCounty(state: string, county: string): Promise<DeviceSchema[]> {
        const response = await DeviceService.getDevicesByCounty(state, county);
        if (response.status === 200) {
            return response.info as DeviceSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get all OLT devices by a county and a state.
     * 
     * @param {string} state Name of the state.
     * @param {string} county Name of the county.
     * @param {string} municipality Name of the municipality.
     * @returns {DeviceSchema[]} List of OLT devices available.
     */
    static async getAllDevicesByMunicipality(state: string, county: string, municipality: string): Promise<DeviceSchema[]> {
        const response = await DeviceService.getDevicesByMunicipality(state, county, municipality);
        if (response.status === 200) {
            return response.info as DeviceSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get all interfaces of an OLT device.
     * 
     * @param {number} deviceID OLT device ID.
     * @returns {InterfaceSchema[]} List of interfaces of the OLT device.
     */
    static async getAllCards(deviceID: number): Promise<InterfaceSchema[]> {
        const response = await InterfaceService.getInterfacesByShell(deviceID);
        if (response.status === 200) {
            return response.info as InterfaceSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get all cards available of an OLT device.
     * 
     * @param {number} id OLT device ID.
     * @returns {number[]} List of all cards available.
     */
    static async getAllCardNumbers(id: number): Promise<number[]> {
        const response = await InterfaceService.getInterfacesByShell(id);
        if (response.status === 200) {
            let cards: number[] = [];
            response.info.map((interface_: InterfaceSchema) => {
                let ifName = interface_.ifname;
                if (ifName.includes(Strings.GPON)) {
                    let card = parseInt(ifName.split("/")[1]);
                    if (!cards.includes(card)) cards.push(card);
                }
            });
            cards = cards.sort((a, b) => a - b);
            return cards;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get all interfaces of an OLT device by a card.
     * 
     * @param {number} deviceID OLT device ID.
     * @param {number} card OLT device card.
     * @returns {InterfaceSchema[]} List of interfaces of the OLT device.
     */
    static async getAllPorts(id: number, card: number): Promise<InterfaceSchema[]> {
        const response = await InterfaceService.getInterfacesByCard(id, card);
        if (response.status === 200) {
            return response.info as InterfaceSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get all ports available of an OLT device.
     * 
     * @param {number} id OLT device ID.
     * @param {number} card OLT device card.
     * @returns {number[]} List of all ports available.
     */    
    static async getAllPortNumbers(id: number, card: number): Promise<number[]> {
        const response = await InterfaceService.getInterfacesByCard(id, card);
        if (response.status === 200) {
            let cards: number[] = [];
            response.info.map((interface_: InterfaceSchema) => {
                let ifName = interface_.ifname;
                if (ifName.includes(Strings.GPON)) {
                    let card = parseInt(ifName.split("/")[2]);
                    if (!cards.includes(card)) cards.push(card);
                }
            });
            cards = cards.sort((a, b) => a - b);
            return cards;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get information about an interface of an OLT device.
     * 
     * @param {number} id OLT device ID.
     * @param {number} card OLT device card.
     * @param {number} port OLT device port.
     * @returns {InterfaceSchema | null} Information about the interface.
     */
    static async getInterface(id: number, card: number, port: number): Promise<InterfaceSchema | null> {
        const response = await InterfaceService.getInterfacesByPort(id, card, port);
        if (response.status === 200) {
            return response.info?.[0] as InterfaceSchema;
        } else {
            console.error(response.info!.message);
            return null;
        }
    }


    /**
     * Get information about an OLT device by its ID.
     * 
     * @param {number} id OLT device ID.
     * @returns {DeviceSchema | null} Information about the OLT device.
     */
    static async getDeviceByID(id: number): Promise<DeviceSchema | null> {
        const response = await DeviceService.getDeviceByID(id);
        if (response.status === 200) {
            return response.info as DeviceSchema;
        } else {
            console.error(response.info!.message);
            return null;
        }
    }


    /**
     * Get information about an OLT device by its sysname.
     * 
     * @param {string} sysname OLT device sysname.
     * @returns 
     */
    static async getDeviceBySysname(sysname: string): Promise<DeviceSchema | null> {
        const response = await DeviceService.getDeviceBySysname(sysname);
        if (response.status === 200) {
            return response.info as DeviceSchema;
        } else {
            console.error(response.info!.message);
            return null;
        }
    }


    /**
     * Get information about an interface of an OLT device by its ID.
     * 
     * @param {number} interfaceID Interface ID.
     * @returns {InterfaceSchema | null} Information about the interface.
     */
    static async getInterfaceByID(interfaceID: number): Promise<InterfaceSchema | null> {
        const response = await InterfaceService.getInterfaceByID(interfaceID);
        if (response.status === 200) {
            return response.info as InterfaceSchema;
        } else {
            console.error(response.info!.message);
            return null;
        }
    }
}