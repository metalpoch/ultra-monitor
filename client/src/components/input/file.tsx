import { TypeReports } from '../../constant/reports';
import React from 'react';

export default function InputFileComponent() {
    return (
        <div className="min-w-fit w-full h-fit bg-blue-900 flex flex-col justify-center items-center py-8 max-md:rounded-b-xl lg:rounded-br-xl">
            <div className="w-full h-fit flex flex-col justify-center items-center gap-3 py-2 lg:flex-row">
                <h2 className="text-xl text-white font-semibold">Subir archivo</h2>
                <input id={`file-${TypeReports.CATASTRE}`} type="file" className="w-48 text-md text-white" />
            </div>
            <button id={`upload-${TypeReports.CATASTRE}`} className="w-fit h-fit bg-blue-700 text-white font-semibold rounded-full px-10 py-2 transition-all ease-linear duration-300 hover:bg-green-700">Subir</button>
        </div>
    );
}