/**
 * Schema of template of an OLT device.
 */
export interface Template {
    ID: number,
    Name: string,
    OidBw: string,
    OidIn: string,
    OidOut: string,
    CreatedAt: Date,
    UpdatedAt: Date,
}

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

