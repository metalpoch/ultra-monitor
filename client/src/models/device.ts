import type { Template } from "./template";

export interface Device {
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