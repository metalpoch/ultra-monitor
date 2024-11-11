import CalendarSelectorComponent from "../selector/calendar";
import BasicSelectorComponent from "../selector/basic";
import TimeSelectorComponent from "../selector/time";
import React from "react";
import { useState, useEffect } from "react";
import { DeviceController } from "../../controllers/device";
import { LocationController } from "../../controllers/location";
import { Strings } from "../../constant/strings";
import type { FilterOptions } from "../../models/filter";
import type { Device } from "../../models/device";

interface Props {
    onClick: (filters: FilterOptions) => void;
    onClickOptionFilter?: (option: string) => void;
}

export default function TrafficFilterComponent(content: Props) {
    // OPTIONS (BUTTONS)
    const [optionFilter, setOptionFilter] = useState<string>(Strings.EQUIPMENT);
    const [filterAvailable, setFilterAvailable] = useState<boolean>(false);

    // DATA COLLECTED (FILTER OPTIONS)
    const [devices, setDevices] = useState<string[]>([]);
    const [ODNs, setODNs] = useState<string[]>([]);
    const [cards, setCards] = useState<number[]>([]);
    const [ports, setPorts] = useState<number[]>([]);
    const [states, setStates] = useState<string[]>([]);
    const [counties, setCounties] = useState<string[]>([]);
    const [municipalities, setMunicipalities] = useState<string[]>([]);

    // DATA SELECTED (FILTER SELECTION)
    const [device, setDevice] = useState<Device>();
    const [fromDate, setFromDate] = useState<string>();
    const [toDate, setToDate] = useState<string>();
    const [state, setState] = useState<string>();
    const [county, setCounty] = useState<string>();
    const [municipality, setMunicipality] = useState<string>();
    const [card, setCard] = useState<number>();
    const [port, setPort] = useState<number>();

    const handlerSearchFilter = async () => {
        if (filterAvailable) {
            let filter: FilterOptions = {
                optionFilter: optionFilter,
                fromDate: fromDate,
                toDate: toDate,
                device: device,
                card: card,
                port: port,
                state: state,
                county: county,
                municipality: municipality
            }
            content.onClick(filter);
        }
    };

    const handlerOptionFilterChange = (newOption: string) => {
        setOptionFilter(newOption);
        if (content.onClickOptionFilter) content.onClickOptionFilter(newOption);
    };

    const handlerFromDateChange = (newDate: string) => {
        setFromDate(newDate);
    };

    const handlerToDateChange = (newDate: string) => {
        setToDate(newDate);
    };

    const handlerStateChange = async (newState: string) => {
        setMunicipalities([]);
        setCounties([]);
        if (newState === Strings.EMPTYSELECTION) setState(undefined);
        else {
            setState(newState);
            let currentCounties = await LocationController.getCounties(newState);
            setCounties(currentCounties);
        }
    };

    const handlerCountyChange = async (newCounty: string) => {
        if (newCounty === Strings.EMPTYSELECTION) {
            setCounty(undefined);
            setMunicipalities([]);
        } else {
            setCounty(newCounty);
            if (state) {
                let currentMunicipalities = await LocationController.getMunicipalities(state, newCounty);
                setMunicipalities(currentMunicipalities);   
            }
        }
    };

    const handlerMunicipalityChange = (newMunicipality: string) => {
        if (newMunicipality === Strings.EMPTYSELECTION) setMunicipality(undefined);
        else setMunicipality(newMunicipality);
    };

    const handlerDeviceChange = async (newAcronym: string) => {
        setCards([]);
        setPorts([]);
        if (newAcronym === Strings.EMPTYSELECTION) setDevice(undefined);
        else {
            let currentDevice = await DeviceController.getDeviceBySysname(newAcronym);
            if (currentDevice) {
                setDevice(currentDevice)
                let currentCards = await DeviceController.getAllCardNumbers(currentDevice.id);
                setCards(currentCards);
            } else setCards([]);
        }
    };

    const handlerCardChange = async (newCard: number) => {
        if (newCard.toString() === Strings.EMPTYSELECTION) {
            setCard(undefined);
            setPorts([]);
        }
        else {
            setCard(newCard);
            if (newCard && device)  {
                let currentPorts = await DeviceController.getAllPortNumbers(device.id, newCard);
                setPorts(currentPorts);
            }
        }
    };

    const handlerPortChange = (newPort: number) => {
        if (newPort.toString() === Strings.EMPTYSELECTION) setPort(undefined);
        else setPort(newPort);
    };

    useEffect(() => {
        const getDevices = async () => {
            const devices = await DeviceController.getAllDevicesNames();
            setDevices(devices);
        }
        const getStates = async () => {
            const states = await LocationController.getStates();
            setStates(states);
        }
        getDevices();
        getStates();
    }, []);

    useEffect(() => {
        if (fromDate === "NaN-NaN-NaN") setFromDate(undefined);
        if (toDate === "NaN-NaN-NaN") setToDate(undefined);
        if (fromDate && toDate) {
            if (optionFilter === Strings.EQUIPMENT) {
                if (device && !card && !port) setFilterAvailable(true);
                else if (!port) setFilterAvailable(false);
                else if (device && card && port) setFilterAvailable(true);
            } else if (optionFilter === Strings.LOCATION) {
                if (!state && !county && !municipality) setFilterAvailable(false);
                else setFilterAvailable(true);
            }
        } else setFilterAvailable(false);
    }, [fromDate, toDate, device, card, port, state, county, municipality]);

    return (
        <div className="min-w-fit w-full h-2/3 max-h-fit p-6 bg-white flex flex-col items-center justify-center gap-4 self-center rounded-xl lg:self-start md:w-1/2 max-sm:w-full max-sm:gap-3">
            <h3 className="text-2xl text-blue-800 text-center font-bold lg:self-start">Búsqueda por</h3>
            <section className="w-full flex flex-row justify-center items-center lg:justify-start gap-3">
                <button type="button" onClick={() => handlerOptionFilterChange(Strings.EQUIPMENT)} className={`w-fit h-fit px-6 py-1 ${optionFilter === Strings.EQUIPMENT ? "bg-blue-600" : "bg-blue-900"} text-white font-bold rounded-full transition-all duration-300 ease-linear ${optionFilter === Strings.EQUIPMENT ? "" : "hover:bg-blue-500"}`}>Equipo</button>
                <button type="button" onClick={() => handlerOptionFilterChange(Strings.LOCATION)} className={`w-fit h-fit px-6 py-1 ${optionFilter === Strings.LOCATION ? "bg-blue-600" : "bg-blue-900"} text-white font-bold rounded-full transition-all duration-300 ease-linear ${optionFilter === Strings.LOCATION ? "" : "hover:bg-blue-500"}`}>Ubicación</button>
            </section>
            <section className="w-full flex flex-col justify-start gap-3">
                <div className="w-full flex flex-row gap-3 flex-wrap lg:flex-nowrap max-sm:flex-col max-sm:gap-3">
                    <CalendarSelectorComponent id="fromDate" label="Desde" onChange={handlerFromDateChange} />
                    <CalendarSelectorComponent id="toDate" label="Hasta" onChange={handlerToDateChange} />
                </div>
            </section>
            {optionFilter === "equipment" && <section className="w-full flex flex-col justify-start gap-3">
                <BasicSelectorComponent id="olt" label="OLT" options={devices} onChange={handlerDeviceChange} />
                <BasicSelectorComponent id="odn" label="ODN" options={ODNs} onChange={handlerDeviceChange} />
                <div className="w-full flex flex-row gap-3 flex-wrap lg:flex-nowrap">
                    <BasicSelectorComponent id="card" label="Tarjeta" options={cards} onChange={handlerCardChange} />
                    <BasicSelectorComponent id="port" label="Puerto" options={ports} onChange={handlerPortChange} />
                </div>
            </section>}
            {optionFilter === "location" && <section className="w-full flex flex-col justify-start gap-3">
                <BasicSelectorComponent id="state" label="Estado" options={states} onChange={handlerStateChange} />
                <BasicSelectorComponent id="county" label="Municipio" options={counties} onChange={handlerCountyChange} />
                <BasicSelectorComponent id="municipality" label="Parroquia" options={municipalities} onChange={handlerMunicipalityChange} />
            </section>}
            <button type="button" id="filterButton" onClick={handlerSearchFilter} className={`w-fit h-fit px-8 py-2 ${filterAvailable ? "bg-blue-700" : "bg-gray-300"} text-white font-bold rounded-full transition-all duration-300 ease-linear ${filterAvailable ? "hover:bg-blue-900" : ""}`}> Filtrar </button>
        </div>
    );
}