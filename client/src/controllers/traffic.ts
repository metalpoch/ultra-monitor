import { TrafficService } from "../services/traffic";
import type { MeasurementSchema } from "../schemas/measurement";

/**
 * @class Controller for all requests to the traffic/measurements.
 */
export class TrafficController {

    /**
     * Get traffic of an OLT device.
     * 
     * @param {number} deviceID OLT device ID.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     * @returns {MeasurementSchema[]} List of traffic of the OLT device.
     */
    static async getDevice(deviceID: number, initialDate: string, endDate: string, initialTime: string = "00:00", endTime: string = "23:59"): Promise<MeasurementSchema[]> {
        const response = await TrafficService.getDevice(deviceID, initialDate, endDate, initialTime, endTime);
        if (response.status === 200) {
            let newDataTraffic : MeasurementSchema[] = [];
            let dataTraffic = response.info as MeasurementSchema[];
            dataTraffic.map((traffic: MeasurementSchema) => {
                let newTraffic: MeasurementSchema = {
                    date: traffic.date,
                    bandwidth_bps: Number(traffic.bandwidth_bps),
                    in_bps: Number(traffic.in_bps),
                    out_bps: Number(traffic.out_bps)
                }
                newDataTraffic.push(newTraffic);
            });
            return newDataTraffic;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get traffic of an interface.
     * 
     * @param {number} interfaceID Interface ID.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     * @returns {MeasurementSchema[]} List of traffic of the interface.
     */
    static async getInterface(interfaceID: number, initialDate: string, endDate: string, initialTime: string = "00:00", endTime: string = "23:59"): Promise<MeasurementSchema[]> {
        const response = await TrafficService.getInterface(interfaceID, initialDate, endDate, initialTime, endTime);
        if (response.status === 200) {
            let newDataTraffic : MeasurementSchema[] = [];
            let dataTraffic = response.info as MeasurementSchema[];
            dataTraffic.map((traffic: MeasurementSchema) => {
                let newTraffic: MeasurementSchema = {
                    date: traffic.date,
                    bandwidth_bps: Number(traffic.bandwidth_bps),
                    in_bps: Number(traffic.in_bps),
                    out_bps: Number(traffic.out_bps)
                }
                newDataTraffic.push(newTraffic);
            });
            return newDataTraffic;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get traffic of a state.
     * 
     * @param {string} state Name of the state.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     * @returns {MeasurementSchema[]} List of traffic of the state.
     */
    static async getState(state: string, initialDate: string, endDate: string, initialTime: string = "00:00", endTime: string = "23:59"): Promise<MeasurementSchema[]> {
        const response = await TrafficService.getState(state, initialDate, endDate, initialTime, endTime);
        if (response.status === 200) {
            let newDataTraffic : MeasurementSchema[] = [];
            let dataTraffic = response.info as MeasurementSchema[];
            dataTraffic.map((traffic: MeasurementSchema) => {
                let newTraffic: MeasurementSchema = {
                    date: traffic.date,
                    bandwidth_bps: Number(traffic.bandwidth_bps),
                    in_bps: Number(traffic.in_bps),
                    out_bps: Number(traffic.out_bps)
                }
                newDataTraffic.push(newTraffic);
            });
            return newDataTraffic;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get traffic of a county.
     * 
     * @param {string} state Name of the state.
     * @param {string} county Name of the county.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     * @returns {MeasurementSchema[]} List of traffic of the county.
     */
    static async getCounty(state: string, county: string, initialDate: string, endDate: string, initialTime: string = "00:00", endTime: string = "23:59"): Promise<MeasurementSchema[]> {
        const response = await TrafficService.getCounty(state, county, initialDate, endDate, initialTime, endTime);
        if (response.status === 200) {
            let newDataTraffic : MeasurementSchema[] = [];
            let dataTraffic = response.info as MeasurementSchema[];
            dataTraffic.map((traffic: MeasurementSchema) => {
                let newTraffic: MeasurementSchema = {
                    date: traffic.date,
                    bandwidth_bps: Number(traffic.bandwidth_bps),
                    in_bps: Number(traffic.in_bps),
                    out_bps: Number(traffic.out_bps)
                }
                newDataTraffic.push(newTraffic);
            });
            return newDataTraffic;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get traffic of a municipality.
     * 
     * @param {string} state Name of the state.
     * @param {string} county Name of the county.
     * @param {string} municipality Name of the municipality.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     * @returns {MeasurementSchema[]} List of traffic of the municipality.
     */
    static async getMunicipality(state: string, county: string, municipality: string, initialDate: string, endDate: string, initialTime: string = "00:00", endTime: string = "23:59"): Promise<MeasurementSchema[]> {
        const response = await TrafficService.getMunicipality(state, county, municipality, initialDate, endDate, initialTime, endTime);
        if (response.status === 200) {
            let newDataTraffic : MeasurementSchema[] = [];
            let dataTraffic = response.info as MeasurementSchema[];
            dataTraffic.map((traffic: MeasurementSchema) => {
                let newTraffic: MeasurementSchema = {
                    date: traffic.date,
                    bandwidth_bps: Number(traffic.bandwidth_bps),
                    in_bps: Number(traffic.in_bps),
                    out_bps: Number(traffic.out_bps)
                }
                newDataTraffic.push(newTraffic);
            });
            return newDataTraffic;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }


    /**
     * Get traffic of an ODN.
     * 
     * @param {string} odn Name of the ODN.
     * @param {string} initialDate Initial date of traffic.
     * @param {string} endDate End date of traffic.
     * @param {string} initialTime Initial hour of traffic.
     * @param {string} endTime End hour of traffic.
     * @returns {MeasurementSchema[]} List of traffic of the ODN.
     */
    static async getOdn(odn: string, initialDate: string, endDate: string, initialTime: string = "00:00", endTime: string = "23:59"): Promise<MeasurementSchema[]> {
        const response = await TrafficService.getOdn(odn, initialDate, endDate, initialTime, endTime);
        if (response.status === 200) {
            let newDataTraffic : MeasurementSchema[] = [];
            let dataTraffic = response.info as MeasurementSchema[];
            dataTraffic.map((traffic: MeasurementSchema) => {
                let newTraffic: MeasurementSchema = {
                    date: traffic.date,
                    bandwidth_bps: Number(traffic.bandwidth_bps),
                    in_bps: Number(traffic.in_bps),
                    out_bps: Number(traffic.out_bps)
                }
                newDataTraffic.push(newTraffic);
            });
            return newDataTraffic;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

}