import type { Template } from "./template";

/**
 * Schema of a OLT device.
 */
export interface DeviceSchema {
    id: number,
    ip: string,
    community: string,
    sysname: string, 
    syslocation: string,
    is_alive: boolean,
    template: Template,
    last_check: Date,
    created_at: Date,
    updated_at: Date,
}