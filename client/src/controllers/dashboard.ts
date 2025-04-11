import { getPreviousMonth } from "../utils/date";
import type { StateTrafficSchema } from "../schemas/measurement";
import type { OdnTrafficSchema } from "../schemas/measurement";
import { DashboardService } from "../services/dashboard";

export class DashboardController {
    /**
     * Get traffic to state to dashboard.
     */
    static async getTrafficState(): Promise<StateTrafficSchema[]> {
        const monthConsult = getPreviousMonth();
        const response = await DashboardService.getTrafficByState(monthConsult);
        console.log(response);
        if (response.status === 200 ) {           
            return response.info as StateTrafficSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    /**
     * Get traffic top N to state to dashboard. 
     * @param {number} n Top N to get. Default N is 5.
     */
    static async getTrafficByStateTopN(n: number = 5): Promise<StateTrafficSchema[]> {
        const monthConsult = getPreviousMonth();
        const response =  await DashboardService.getTrafficByStateTopN(monthConsult, n);
        if (response.status === 200 ) {           
            return response.info as StateTrafficSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    /**
     * Get traffic to ODN to dashboard.
     */
    static async getTrafficByOdn(): Promise<OdnTrafficSchema[]> {
        const monthConsult = getPreviousMonth();
        const response =  await DashboardService.getTrafficByOdn(monthConsult);
        if (response.status === 200 ) {           
            return response.info as OdnTrafficSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

}