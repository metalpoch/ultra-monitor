import type { DeviceSchema } from "./device";

/**
 * Schema of traffic filter options.
 */
export interface FilterOptionSchema {
    optionFilter: string;
    fromDate: string | undefined;
    toDate: string | undefined;
    device: DeviceSchema | undefined;
    card: number | undefined;
    port: number | undefined;
    state: string | undefined;
    county: string | undefined;
    municipality: string | undefined;
    odn: string | undefined;
}