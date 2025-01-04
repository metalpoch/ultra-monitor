import { OdnService } from "../services/odn";

export class OdnController {

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