import { DeviceService } from "../services/device";
import { InterfaceService } from "../services/interface";
import { Strings } from "../constant/strings";
import type { Interface } from "../models/interface";
import type { Device } from "../models/device";

export class DeviceController {
    static async getAllDevices(): Promise<Device[]> {
        const response = await DeviceService.getDevices();
        if (response.status === 200) {
            return response.info as Device[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    static async getAllDevicesNames(): Promise<string[]> {
        const response = await DeviceService.getDevices();
        if (response.status === 200) {
            let deviceNames: string[] = [];
            response.info.map((device: Device) => {
                deviceNames.push(device.sysname);
            });
            deviceNames = deviceNames.sort();
            return deviceNames;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    static async getAllDevicesByState(state: string): Promise<Device[]> {
        const response = await DeviceService.getDevicesByState(state);
        if (response.status === 200) {
            return response.info as Device[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    static async getAllDevicesByCounty(state: string, county: string): Promise<Device[]> {
        const response = await DeviceService.getDevicesByCounty(state, county);
        if (response.status === 200) {
            return response.info as Device[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    static async getAllDevicesByMunicipality(state: string, county: string, municipality: string): Promise<Device[]> {
        const response = await DeviceService.getDevicesByMunicipality(state, county, municipality);
        if (response.status === 200) {
            return response.info as Device[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    static async getAllCards(deviceID: number): Promise<Interface[]> {
        const response = await InterfaceService.getInterfacesByCard(deviceID);
        if (response.status === 200) {
            return response.info as Interface[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    static async getAllCardNumbers(deviceID: number): Promise<number[]> {
        const response = await InterfaceService.getInterfacesByCard(deviceID);
        if (response.status === 200) {
            let cards: number[] = [];
            response.info.map((interface_: Interface) => {
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

    static async getAllPorts(id: number, card: number): Promise<Interface[]> {
        const response = await InterfaceService.getInterfacesByPort(id, card);
        if (response.status === 200) {
            return response.info as Interface[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    static async getAllPortNumbers(id: number, card: number): Promise<number[]> {
        const response = await InterfaceService.getInterfacesByPort(id, card);
        if (response.status === 200) {
            let cards: number[] = [];
            response.info.map((interface_: Interface) => {
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

    static async getInterface(deviceID: number, card: number, port: number): Promise<Interface | null> {
        const response = await InterfaceService.getInterfaces(deviceID, card, port);
        if (response.status === 200) {
            return response.info?.[0] as Interface;
        } else {
            console.error(response.info!.message);
            return null;
        }
    }

    static async getDeviceByID(id: number): Promise<Device | null> {
        const response = await DeviceService.getDeviceByID(id);
        if (response.status === 200) {
            return response.info as Device;
        } else {
            console.error(response.info!.message);
            return null;
        }
    }

    static async getDeviceBySysname(sysname: string): Promise<Device | null> {
        const response = await DeviceService.getDeviceBySysname(sysname);
        if (response.status === 200) {
            return response.info as Device;
        } else {
            console.error(response.info!.message);
            return null;
        }
    }

    static async getInterfaceByID(interfaceID: number): Promise<Interface | null> {
        const response = await InterfaceService.getInterfaceByID(interfaceID);
        if (response.status === 200) {
            return response.info as Interface;
        } else {
            console.error(response.info!.message);
            return null;
        }
    }
}