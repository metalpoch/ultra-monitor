import React from "react";
import Selector from "../buttons/selector";
import CalendarSelector from "../buttons/calendar";



export default function TrafficFilter() {
    return (
        <aside className="w-fit h-full min-w-fit px-4 py-6 bg-white shadow-2xl flex flex-col items-center justify-center gap-4 rounded-xl md:w-1/2 max-sm:w-full max-sm:gap-3">
            <h3 className="text-2xl text-blue-800 text-center font-bold underline underline-offset-8 decoration-4 rounded-full">Búsqueda por</h3>
            <section className="w-full flex flex-col justify-start gap-3">
                <div className="w-full flex flex-row gap-3 flex-wrap lg:flex-nowrap max-sm:flex-col max-sm:gap-3">
                    <CalendarSelector id="fromDate" label="Desde" />
                    <CalendarSelector id="toDate" label="Hasta"/>
                </div>
                <Selector id="state" label="Estado" />
            </section>
            <section className="w-full flex flex-col justify-start gap-3">
                <h4 className="text-lg text-blue-800 font-bold underline underline-offset-8 decoration-4 rounded-full">Opcional: Ubicación</h4>
                <Selector id="county" label="Municipio" />
                <Selector id="municipality" label="Parroquia" />
            </section>
            <section className="w-full flex flex-col justify-start gap-3">
                <h4 className="text-lg text-blue-800 font-bold underline underline-offset-8 decoration-4 rounded-full">Opcional: Equipo</h4>
                <Selector id="olt" label="OLT" />
                <div className="w-full flex flex-row gap-3 flex-wrap lg:flex-nowrap">
                    <Selector id="card" label="Tarjeta" />
                    <Selector id="port" label="Puerto" />
                </div>
            </section>
            <button type="button" className="w-fit h-fit px-8 py-2 bg-blue-700 text-white font-bold rounded-full transition-all duration-300 ease-linear hover:bg-blue-900">Buscar</button>
        </aside>
    );
}