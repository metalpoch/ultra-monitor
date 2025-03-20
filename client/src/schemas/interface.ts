/**
 * Schema of an interface.
 */
export interface InterfaceSchema {
    id: number,
    ifindex: number,
    ifname: string,
    ifDescr: string,
    ifAlias: string,
    created_at: Date,
    updated_at: Date,
}