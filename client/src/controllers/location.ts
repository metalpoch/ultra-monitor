import { LocationService } from "../services/location";

export class LocationController {

    static async getStates(): Promise<string[]> {
        const response = await LocationService.getStates();
        if (response.status === 200) {
            let counties = response.info.sort();
            return counties;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    static async getCounties(state: string): Promise<string[]> {
        const response = await LocationService.getCounties(state);
        if (response.status === 200) {
            let counties = response.info.sort();
            return counties;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }

    static async getMunicipalities(state: string, county: string): Promise<string[]> {
        const response = await LocationService.getMunicipalities(state, county);
        if (response.status === 200) {
            let municipality = response.info.sort();
            return municipality;
        } else {
            console.error(response.info!.message);
            return [];
        }
    }
}