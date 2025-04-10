import { getMonth } from "../utils/date";
import type { StateTrafficSchema } from "../schemas/measurement";
import type { OdnTrafficSchema } from "../schemas/measurement";
import { DashboardService } from "../services/dashboard";

export class DashboardController {
    /**
     * Get traffic to state to dashboard.
     */
    static async getTrafficState(): Promise<StateTrafficSchema[]> {
        const monthConsult = getMonth();
        const response =  await DashboardService.gettrafficbystate(monthConsult);
        if (response.status === 200 ) {           
            return response.info as StateTrafficSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    static async gettrafficbystatetopN(): Promise<StateTrafficSchema[]> {
        const monthConsult = getMonth();
        var n = 5;
        const response =  await DashboardService.gettrafficbystatetopN(monthConsult, n);
        if (response.status === 200 ) {           
            return response.info as StateTrafficSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    static async gettrafficbyodn(): Promise<OdnTrafficSchema[]> {
        const monthConsult = getMonth();
        const response =  await DashboardService.gettrafficbyodn(monthConsult);
        if (response.status === 200 ) {           
            return response.info as OdnTrafficSchema[];
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

}