import { LocationService } from "../services/location";

/**
 * @class Controller for all requests to the location.
 */
export class LocationController {

    /**
     * Get all states avaliable.
     * 
     * @returns {string[]} List of all states avaliable.
     */
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


    /**
     * Get all county avaliable by a state.
     * 
     * @param state Name of the state.
     * @returns {string[]} List of all counties avaliable.
     */
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


    /**
     * Get all municipality avaliable by a state and a county.
     * 
     * @param {string} state Name of the state.
     * @param {string} county Name of the county.
     * @returns {string[]} List of all municipalities avaliable.
     */
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