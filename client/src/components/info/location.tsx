import React, { useState } from "react";
import SpinnerBasicComponent from "../spinner/basic";
import InfoEquipmentModalComponent from "../modal/equipment";
import { Device } from "../../models/device";
import { LoadStatus } from "../../constant/loadStatus";
import type { Info } from "../../models/info";
import type { LoadingStateValue } from "../../types/loadingState";

interface Props {
    loading: LoadingStateValue;
    info?: Info;
}

export default function EquipmentInfoComponent(content: Props) {

    const [showInfoModal, setShowInfoModal] = useState<boolean>(false);
    const [deviceSelected, setDeviceSelected] = useState<Device>();

    const handlerHiddenInfoModal = () => {
        setShowInfoModal(false);
        setDeviceSelected(undefined);
    }

    const handlerShowInfoModal = (device: Device) => {
        setShowInfoModal(true);
        setDeviceSelected(device);
    }

    return(<>
        {deviceSelected && <InfoEquipmentModalComponent showModal={showInfoModal} device={deviceSelected} onClick={handlerHiddenInfoModal} />}
        {content.loading === LoadStatus.EMPTY && 
            <div className="w-full min-h-52 bg-white p-6 rounded-md flex flex-col justify-center items-center">
                <h4 className="text-xl font-semibold text-gray-300 text-center">Sin búsqueda de ubicación</h4>
            </div>
        }
        {content.loading === LoadStatus.LOADING && 
            <div className="w-full min-h-52 bg-white p-6 rounded-md flex flex-col justify-center items-center">
                <SpinnerBasicComponent />
            </div>
        }
        {content.loading === LoadStatus.LOADED && !content.info && 
            <div className="w-full min-h-52 bg-white p-6 rounded-md flex flex-col justify-center items-center">
                <h4 className="text-xl font-semibold text-gray-300 text-center">No se encontró información de la ubicación</h4>
            </div>
        }
        {content.loading === LoadStatus.LOADED && content.info &&
            <div className="min-w-fit w-full h-fit bg-white p-6 rounded-md flex flex-col gap-2">
                {content.info.state && <>
                    <section className="w-full h-fit flex flex-col items-center gap-2 lg:flex-row lg:gap-10">
                        <h2 className="text-3xl text-blue-700 font-bold">Estado: 
                            <span className="text-3xl text-blue-700"> {content.info.state}</span>
                        </h2>
                        {content.info?.county &&
                            <h2 className="text-xl text-blue-800 font-semibold">Municipio: 
                                <span className="text-xl text-gray-700"> {content.info.county}</span>
                            </h2>
                        }
                        {content.info?.municipality &&
                            <h2 className="text-xl text-blue-800 font-semibold">Parroquia: 
                                <span className="text-xl text-gray-700"> {content.info.municipality}</span>
                            </h2>
                        }
                    </section>
                </>}
                {content.info.otherDevices && <>
                    <h2 className="text-2xl text-blue-800 font-semibold">Equipos OLT pertenecientes:</h2>
                    <section className="w-full max-h-80 flex flex-row flex-wrap gap-2 overflow-y-auto">
                        {content.info.otherDevices.map((device: Device, index: number) => (
                            <button type="button" key={index} onClick={() => {handlerShowInfoModal(device)}} className={`w-fit h-fit flex flex-row items-center gap-2 rounded-md ${device.is_alive ? "bg-green-700": "bg-red-700"} px-4 py-1`}>
                                <h3 className="text-xl text-white font-semibold">{device.sysname}</h3>                            
                            </button>
                        ))}
                    </section>
                </>}
            </div>
        }
    </>);
}