import { getMonth } from "../utils/date";
import type { StateTrafficSchema } from "../schemas/measurement";

export class DashboardController {
    /**
     * Get traffic to state to dashboard.
     */
    static async getTrafficState(): Promise<StateTrafficSchema[]> {
        const monthConsult = getMonth();
        const response = { status: 200, data: [] as StateTrafficSchema[], err: { message: null } };
        if (response.status === 200 && response.data) {           
            return response.data as StateTrafficSchema[];
        } else {
            console.error(response.err!.message);
            return [];
        }
}
}