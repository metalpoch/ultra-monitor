import type { Device } from "./device";

export interface FilterOptions {
    optionFilter: string;
    fromDate: string | undefined;
    toDate: string | undefined;
    device: Device | undefined;
    card: number | undefined;
    port: number | undefined;
    state: string | undefined;
    county: string | undefined;
    municipality: string | undefined;
}