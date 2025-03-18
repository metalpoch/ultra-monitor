import SpinnerBasicComponent from "../spinner/basic";
import { LoadStatus } from "../../constant/loadStatus";
import { sortInterfacesByDate, sortInterfacesByBandwidth, sortInterfacesByIn, sortInterfacesByOut } from "../../utils/sort";
import { getUnit, mbpsToGbps } from "../../utils/transform";
import type { MeasurementSchema } from "../../schemas/measurement";
import type { SortTraffic } from "../../types/sortTraffic";
import type { LoadingStateValue } from "../../types/loadingState";
import React, { useState, useEffect } from "react";

/**
 * @interface Data required for the traffic table.
 * 
 * @param {LoadingStateValue} loading Loading state of the table component.
 * @param {MeasurementSchema[]} data Data of the table.
 */
interface TrafficProps {
    loading: LoadingStateValue;
    data: MeasurementSchema[];
}


/**
 * @interface Data required for the traffic table.
 * 
 * @param {string} date Date of the traffic.
 * @param {string} time Time of the traffic.
 * @param {string} bandwidth Bandwidth of the traffic.
 * @param {string} in In of the traffic.
 * @param {string} out Out of the traffic.
 */
interface DataTraffic {
    date: string;
    time: string
    bandwidth: string;
    in: string;
    out: string;
}

export default function TrafficTableComponent(content: TrafficProps) {
    
    const [interfaces, setInterfaces] = useState<DataTraffic[]>([]);

    /**
     * Transform the data of the traffic to data available for the table.
     * 
     * @param {MeasurementSchema[]} data Data of the traffic.
     * @returns {DataTraffic[]} Data of the traffic to the table.
     */
    const transformData = (data: MeasurementSchema[]): DataTraffic[] => {
        let dataTraffic: DataTraffic[] = [];
        data.map((interface_: MeasurementSchema) => {
            let newInterface: DataTraffic = {
                date: `${interface_.date.toString().split("T")[0]}`,
                time: `${interface_.date.toString().split("T")[1].split(".")[0]}`,
                bandwidth: getUnit(interface_.bandwidth_bps),
                in: getUnit(interface_.in_bps),
                out: getUnit(interface_.out_bps)
            }
            dataTraffic.push(newInterface);
        });
        return dataTraffic;
    }


    /**
     * Sort the data of the traffic by bandwidth.
     */
    const sortByBandwidth = () => {
        if (content.data?.length > 0) {
            let sortedInterfaces = sortInterfacesByBandwidth(content.data);
            let dataTraffic = transformData(sortedInterfaces);
            setInterfaces([...dataTraffic]);
        }
    }


    /**
     * Sort the data of the traffic by in.
     */
    const sortByIn = () => {
        if (content.data?.length > 0) {
            let sortedInterfaces = sortInterfacesByIn(content.data);
            let dataTraffic = transformData(sortedInterfaces);
            setInterfaces([...dataTraffic]);
        }
    }


    /**
     * Sort the data of the traffic by out.
     */
    const sortByOut = () => {
        if (content.data?.length > 0) {
            let sortedInterfaces = sortInterfacesByOut(content.data);
            let dataTraffic = transformData(sortedInterfaces);
            setInterfaces([...dataTraffic]);
        }
    }
    

    /**
     * Sort the data of the traffic by date.
     */
    const sortByDate = () => {
        if (content.data?.length > 0) {
            let sortedInterfaces = sortInterfacesByDate(content.data);
            let dataTraffic = transformData(sortedInterfaces);
            setInterfaces([...dataTraffic]);
        }
    }


    /**
     * Handler to sort the data to the table.
     */
    const handlerSort = (event: React.MouseEvent<HTMLSelectElement>) => {
        let option = event.currentTarget.value as SortTraffic;
        if (option === "date") sortByDate();
        else if (option === "bandwidth") sortByBandwidth();
        else if (option === "in") sortByIn();
        else if (option === "out") sortByOut();
    }

    useEffect(() => {
        setInterfaces(transformData(content.data));
    }, [content.data]);

    return(<>
        {content.loading === LoadStatus.EMPTY &&
            <div className="min-w-80 w-full min-h-52 bg-white py-4 px-6 flex flex-col justify-center items-center rounded-md">
                <h4 className="text-xl font-semibold text-gray-300 text-center">Sin búsqueda</h4>
            </div>
        }
        {content.loading === LoadStatus.LOADING &&
            <div className="min-w-80 w-full min-h-52 bg-white py-4 px-6 flex flex-col justify-center items-center rounded-md">
                <SpinnerBasicComponent />
            </div>
        }
        {content.loading === LoadStatus.LOADED && !interfaces &&
            <div className="min-w-80 w-full min-h-52 bg-white py-4 px-6 flex flex-col justify-center items-center rounded-md">
                <h4 className="text-xl font-semibold text-gray-300 text-center">No se encontró información del tráfico</h4>
            </div>
        }
        {content.loading === LoadStatus.LOADED && interfaces && interfaces.length <= 0 &&
            <div className="min-w-80 w-full min-h-52 bg-white py-4 px-6 flex flex-col justify-center items-center rounded-md">
                <h4 className="text-xl font-semibold text-gray-300 text-center">No se encontró información del tráfico</h4>
            </div>
        }
        {content.loading === LoadStatus.LOADED && interfaces && interfaces.length > 0 && 
            <div className="min-w-80 w-full h-fit bg-white py-4 px-6 flex flex-col gap-3 rounded-md">
                <section className="min-w-fit w-full flex flex-row flex-wrap items-center justify-between max-sm:gap-4">
                    <h3 className="w-fit h-8 text-2xl text-blue-800 text-nowrap font-bold underline underline-offset-8 decoration-4">Datos del Tráfico</h3>
                    <div className="w-fit min-w-fit flex flex-row justify-end items-center">
                        <label className="px-4 py-1 bg-blue-800 rounded-l-full text-white">Ordenado Por</label>
                        <select name="" id="" className="h-8 outline-none bg-gray-100 rounded-r-full" onClick={handlerSort}>
                            <option value="date">Fecha</option>
                            <option value="bandwidth">Ancho de Banda</option>
                            <option value="in">In (bps)</option>
                            <option value="out">Out (bps)</option>
                        </select>
                    </div>
                </section>
                <section className="w-full h-96 overflow-y-auto">
                    <table className="w-full">
                        <thead>
                            <tr className="border-b-2 border-blue-800">
                                <th className="py-2 text-md text-gray-500">Ancho de Banda (bps)</th>
                                <th className="py-2 text-md text-gray-500">In (bps)</th>
                                <th className="py-2 text-md text-gray-500">Out (bps)</th>
                                <th className="py-2 text-md text-gray-500">Fecha</th>
                                <th className="py-2 text-md text-gray-500">Hora</th>
                            </tr>
                        </thead>
                        <tbody>
                            {interfaces.map((interface_: DataTraffic, index: number) => (
                                <tr key={index} className="border-b border-gray-100">
                                    <td className="py-2 text-md text-center text-gray-500">{interface_.bandwidth}</td>
                                    <td className="py-2 text-md text-center text-gray-500">{interface_.in}</td>
                                    <td className="py-2 text-md text-center text-gray-500">{interface_.out}</td>
                                    <td className="py-2 text-md text-center text-gray-500">{interface_.date}</td>
                                    <td className="py-2 text-md text-center text-gray-500">{interface_.time}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </section>
            </div>
        }
    </>);
}