import type { Measurement } from "../models/measurement";

export const sortInterfacesByDate = (measurements: Measurement[]): Measurement[] => {
    return measurements.sort((a, b) => {
      return a.date.getTime() - b.date.getTime();
    });
}

export const sortInterfacesByBandwidth = (measurements: Measurement[]): Measurement[] => {
    return measurements.sort((a, b) => {
      return a.bandwidth_bps - b.bandwidth_bps;
    });
}

export const sortInterfacesByIn = (measurements: Measurement[]): Measurement[] => {
    return measurements.sort((a, b) => {
      return a.in_bps - b.in_bps;
    });
}

export const sortInterfacesByOut = (measurements: Measurement[]): Measurement[] => {
    return measurements.sort((a, b) => {
      return a.out_bps - b.out_bps;
    });
}