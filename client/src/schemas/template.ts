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