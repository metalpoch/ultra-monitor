import type { DeviceSchema } from "./device";

/**
 * Schema information general of an OLT device or a location.
 */
export interface InfoSchema {
    device?: DeviceSchema;
    card?: number;
    port?: number;
    state?: string;
    county?: string;
    municipality?: string;
    otherDevices?: DeviceSchema[];
    odn?: string;
}