import React from "react";
import SpinnerBasicComponent from "../spinner/basic";
import { LoadStatus } from "../../constant/loadStatus";
import type { Info } from "../../models/info";
import type { LoadingStateValue } from "../../types/loadingState";

interface Props {
    loading: LoadingStateValue;
    info?: Info;
}

export default function EquipmentInfoComponent(content: Props) {
    return(<>
        {content.loading === LoadStatus.EMPTY && 
            <div className="w-full min-h-52 bg-white p-6 rounded-md flex flex-col justify-center items-center">
                <h4 className="text-xl font-semibold text-gray-300 text-center">Sin búsqueda de equipo</h4>
            </div>
        }
        {content.loading === LoadStatus.LOADING && 
            <div className="w-full min-h-52 bg-white p-6 rounded-md flex flex-col justify-center items-center">
                <SpinnerBasicComponent />
            </div>
        }
        {content.loading === LoadStatus.LOADED && !content.info && 
            <div className="w-full min-h-52 bg-white p-6 rounded-md flex flex-col justify-center items-center">
                <h4 className="text-xl font-semibold text-gray-300 text-center">No se encontró información del equipo</h4>
            </div>
        }
        {content.loading === LoadStatus.LOADED && content.info &&
            <div className="w-full h-fit bg-white p-6 rounded-md flex flex-col gap-2">
                {content.info.device?.sysname && 
                    <h2 className="text-3xl text-blue-700 font-bold">OLT: 
                        <span className="text-3xl text-blue-700"> {content.info.device.sysname}</span>
                    </h2>
                }
                <div className="w-full h-fit flex flex-row gap-16">
                    <section className="w-fit flex flex-col">
                        {content.info.device?.ip &&
                            <h2 className="text-xl text-blue-800 font-semibold">IP:
                                <span className="text-xl text-gray-700"> {content.info.device.ip}</span>
                            </h2>
                        }
                        {content.info.device?.community &&
                            <h2 className="text-xl text-blue-800 font-semibold">Comunidad: 
                                <span className="text-xl text-gray-700"> {content.info.device.community}</span>
                            </h2>
                        }
                    </section>
                    {content.info.card && content.info.port && 
                        <section className="w-fit flex flex-col">
                            <h2 className="text-xl text-blue-800 font-semibold">Shell: 
                                <span className="text-xl text-gray-700"> 0</span>
                            </h2>
                            <h2 className="text-xl text-blue-800 font-semibold">Tarjeta: 
                                <span className="text-xl text-gray-700"> {content.info.card}</span>
                            </h2>
                            <h2 className="text-xl text-blue-800 font-semibold">Puerto: 
                                <span className="text-xl text-gray-700"> {content.info.port}</span>
                            </h2>
                        </section>
                    }
                    <section className="w-fit flex flex-col">
                        {content.info.device?.syslocation && 
                            <h2 className="text-xl text-blue-800 font-semibold">Localidad: 
                                <span className="text-xl text-gray-700"> {content.info.device.syslocation}</span>
                            </h2>
                        }
                        <h2 className={`text-xl text-blue-800 font-semibold`}>Estatus: 
                            <span className={`${content.info.device?.is_alive ? "text-green-600" : "text-red-600"}`}> {content.info.device?.is_alive ? "Activo" : "Inactivo"}</span>
                        </h2>
                        {content.info.device?.template && 
                            <h2 className="text-xl text-blue-800 font-semibold">Equipo: 
                                <span className="text-xl text-gray-700"> {content.info.device.template.Name}</span>
                            </h2>
                        }
                    </section>
                    <section className="w-fit flex flex-col">
                        {content.info.state && 
                            <h2 className="text-xl text-blue-800 font-semibold">Estado: 
                                <span className="text-xl text-gray-700"> {content.info.state}</span>
                            </h2>
                        }
                        {content.info.municipality && 
                            <h2 className="text-xl text-blue-800 font-semibold">Municipio: 
                                <span className="text-xl text-gray-700"> {content.info.municipality}</span>
                            </h2>
                        }
                        {content.info.county && 
                            <h2 className="text-xl text-blue-800 font-semibold">Parroquia: 
                                <span className="text-xl text-gray-700"> {content.info.county}</span>
                            </h2>
                        }
                    </section>
                </div>
                {content.info.device?.last_check && <h2 className="text-xs text-gray-900 font-semibold">Última Revisión del Equipo: {content.info.device?.last_check.toString().split("T")[0]} {content.info.device?.last_check.toString().split("T")[1].split(".")[0]}</h2>}
            </div>
        }
    </>);
}