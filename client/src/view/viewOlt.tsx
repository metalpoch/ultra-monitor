import TrafficFilterComponent from "../components/filter/traffic";
import LineGraphicComponent from "../components/graphic/line";
import TrafficTableComponent from '../components/tables/traffic';
import EquipmentInfoComponent from '../components/info/equipment';
import LocationInfoComponent from '../components/info/location';
import React, { useState } from "react";
import { Routes } from '../constant/routes';
import { Strings } from '../constant/strings';
import { LoadStatus } from '../constant/loadStatus';
import { TrafficController } from '../controllers/traffic';
import { DeviceController } from "../controllers/device";
import type { LoadingStateValue } from '../types/loadingState';
import type { Measurement } from '../models/measurement';
import type { FilterOptions } from '../models/filter';
import type { Info } from '../models/info';

export default function ViewOLT() {
    const [loadingData, setLoadingData] = useState<LoadingStateValue>(LoadStatus.EMPTY);
    const [dataTraffic, setDataTraffic] = useState<Measurement[]>([]);
    const [optionFilter, setOptionFilter] = useState<string>(Strings.EQUIPMENT);
    const [infoType, setInfoType] = useState<string>();
    const [info, setInfo] = useState<Info>();

    const handlerTraffic = async (filters: FilterOptions) => {
        setLoadingData(LoadStatus.LOADING);
        setInfoType(filters.optionFilter);
        if (filters.fromDate && filters.toDate) {
            let traffic: Measurement[] = [];
            if (filters.optionFilter === Strings.EQUIPMENT) {
                if (filters.device && filters.card && filters.port) {
                    let info: Info = { device: filters.device, card: filters.card, port: filters.port }
                    setInfo(info);
                    const interface_ = await DeviceController.getInterface(filters.device.id, filters.card, filters.port);
                    if (interface_) {
                        traffic = await TrafficController.getInterface(interface_.id, filters.fromDate, filters.toDate);
                        if (traffic) setDataTraffic(traffic);
                        setLoadingData(LoadStatus.LOADED);
                    }
                } else if (filters.device) {
                    let info: Info = { device: filters.device }
                    setInfo(info);
                    traffic = await TrafficController.getDevice(filters.device.id, filters.fromDate, filters.toDate);
                    if (traffic) setDataTraffic(traffic);
                    setLoadingData(LoadStatus.LOADED);
                }
            } else if (filters.optionFilter === Strings.ODN) {
                if (filters.device && filters.odn) {
                    let info: Info = { device: filters.device, odn: filters.odn }
                    setInfo(info);
                    traffic = await TrafficController.getOdn(filters.odn, filters.fromDate, filters.toDate);
                    if (traffic) setDataTraffic(traffic);
                    setLoadingData(LoadStatus.LOADED);
                }
            } else if (filters.optionFilter === Strings.LOCATION) {
                if (filters.state && filters.county && filters.municipality) {
                    let info: Info
                    traffic = await TrafficController.getMunicipality(filters.state, filters.county, filters.municipality, filters.fromDate, filters.toDate);
                    const devices = await DeviceController.getAllDevicesByMunicipality(filters.state, filters.county, filters.municipality);
                    if (traffic) setDataTraffic(traffic);
                    if (devices) info = { state: filters.state, county: filters.county, municipality: filters.municipality, otherDevices: devices }
                    else info = { state: filters.state, county: filters.county, municipality: filters.municipality }
                    setInfo(info);
                    setLoadingData(LoadStatus.LOADED);
                } else if (filters.state && filters.county) {
                    let info: Info = { state: filters.state, county: filters.county }
                    setInfo(info);
                    traffic = await TrafficController.getCounty(filters.state, filters.county, filters.fromDate, filters.toDate);
                    if (traffic) setDataTraffic(traffic);
                    setLoadingData(LoadStatus.LOADED);
                } else if (filters.state) {
                    let info: Info = { state: filters.state }
                    setInfo(info);
                    traffic = await TrafficController.getState(filters.state, filters.fromDate, filters.toDate);
                    if (traffic) setDataTraffic(traffic);
                    setLoadingData(LoadStatus.LOADED);
                }
            }
        }
    }

    const handlerOptionFilterChange = (newOption: string) => {
        setOptionFilter(newOption);
    }

    return(
        <main className="min-w-80 h-fit flex flex-col gap-2">
            <section className="w-full h-fit flex flex-row gap-4 mb-2">
                <a className="w-fit h-fit px-10 py-1 bg-blue-700 text-white font-bold rounded-full cursor-default">Tráfico</a>
                {/* <a className="w-fit h-fit px-8 py-2 bg-blue-800 text-white font-bold rounded-full transition-all duration-300 ease-linear hover:bg-blue-400" href={Routes.OLT_IP}>IP's Activas</a> */}
            </section>
            <section className="w-full h-fit max-h-fit flex flex-col lg:flex-row flex-nowrap gap-2">
                <TrafficFilterComponent loading={loadingData} onClick={handlerTraffic} onClickOptionFilter={handlerOptionFilterChange} />
                <LineGraphicComponent loading={loadingData} title="Tráfico" canvasID="traffic" data={dataTraffic} />
            </section>
            <section>
                {!infoType && (optionFilter === Strings.EQUIPMENT || optionFilter === Strings.ODN) &&
                    <EquipmentInfoComponent loading={loadingData} info={info} />
                }
                {!infoType && optionFilter === Strings.LOCATION &&
                    <LocationInfoComponent loading={loadingData} info={info} />
                }
                {infoType && (infoType === Strings.EQUIPMENT || optionFilter === Strings.ODN) &&
                    <EquipmentInfoComponent loading={loadingData} info={info} />
                }
                {infoType && infoType === Strings.LOCATION &&
                    <LocationInfoComponent loading={loadingData} info={info} />
                }
            </section>
            <section>
                <TrafficTableComponent loading={loadingData} data={dataTraffic} />
            </section>
        </main>
    )    
}