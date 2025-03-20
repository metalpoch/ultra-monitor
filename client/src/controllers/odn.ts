import { OdnService } from "../services/odn";

/**
 * @class Controller for all requests to the ODN.
 */
export class OdnController {

    /**
     * Get all ODN by an OLT device ID.
     * 
     * @param {number} id OLT device ID.
     * @returns {string[]} List of all ODN avaliable.
     */
    static async getOdnByDevice(id: number): Promise<string[]> {
        const response = await OdnService.getOdnByDevice(id);
        console.log(response);
        if (response.status === 200) {
            let odns = response.info as string[];
            return odns;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }
}