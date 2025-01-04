import type { Device } from "./device";

export interface Info {
    device?: Device;
    card?: number;
    port?: number;
    state?: string;
    county?: string;
    municipality?: string;
    otherDevices?: Device[];
    odn?: string;
}