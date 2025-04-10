/**
 * Schema of a measurement of an interface.
 */
export interface MeasurementSchema {
    date: string;
    bandwidth_bps: number;
    in_bps: number;
    out_bps: number;
}

/**
 * Schema of a measurement of traffic to a state.
 */
export interface StateTrafficSchema {
    state: string;
    bandwidth_bps: number;
    in_bps: number;
    out_bps: number;
}

export interface OdnTrafficSchema {
    odn: string;
    bandwidth_bps: number;
    in_bps: number;
    out_bps: number;
}