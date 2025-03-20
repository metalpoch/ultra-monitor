/**
 * Schema of a measurement of an interface.
 */
export interface MeasurementSchema {
    date: string;
    bandwidth_bps: number;
    in_bps: number;
    out_bps: number;
}