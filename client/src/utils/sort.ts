import type { MeasurementSchema } from "../schemas/measurement";

/**
 * Sort a list of measurament by 'date'.
 * 
 * @param {MeasurementSchema[]} measurements List of measuraments of interfaces.
 * @returns {MeasurementSchema[]} Sorted list of measuraments of interfaces.
 */
export const sortInterfacesByDate = (measurements: MeasurementSchema[]): MeasurementSchema[] => {
    let data = measurements.sort((a, b) => {
      return new Date(a.date).getTime() - new Date(b.date).getTime();
    });
    console.log(data);
    return data;
}


/**
 * Sort a list of measurament by 'bandwith'.
 * 
 * @param {MeasurementSchema[]} measurements List of measuraments of interfaces.
 * @returns {MeasurementSchema[]} Sorted list of measuraments of interfaces.
 */
export const sortInterfacesByBandwidth = (measurements: MeasurementSchema[]): MeasurementSchema[] => {
    return measurements.sort((a, b) => {
      return a.bandwidth_bps - b.bandwidth_bps;
    });
}


/**
 * Sort a list of measurament by 'In'.
 * 
 * @param {MeasurementSchema[]} measurements List of measuraments of interfaces.
 * @returns {MeasurementSchema[]} Sorted list of measuraments of interfaces.
 */
export const sortInterfacesByIn = (measurements: MeasurementSchema[]): MeasurementSchema[] => {
    return measurements.sort((a, b) => {
      return a.in_bps - b.in_bps;
    });
}


/**
 * Sort a list of measurament by 'Out'.
 * 
 * @param {MeasurementSchema[]} measurements List of measuraments of interfaces.
 * @returns {MeasurementSchema[]} Sorted list of measuraments of interfaces.
 */
export const sortInterfacesByOut = (measurements: MeasurementSchema[]): MeasurementSchema[] => {
    return measurements.sort((a, b) => {
      return a.out_bps - b.out_bps;
    });
}