import NoticeModalComponent from '../modal/notice';
import { ReportController } from '../../controllers/report';
import { TypeReports } from '../../constant/reports';
import React, { useState } from 'react';

interface Props {
    idUser: number;
}

export default function InputFileComponent(content: Props) {

    const [file, setFile] = useState<File | null>(null);
    const [uploadError, setUploadError] = useState<boolean>(false);
    const [uploadSuccess, setUploadSuccess] = useState<boolean>(false);

    const handlerGetFile = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (!file) return;
        setFile(file);
    }

    const handlerUpload = async () => {
        if (!file) return;
        const status = await ReportController.uploadFile(file, content.idUser);
        if (status) setUploadSuccess(true);
        else setUploadError(true);
    }

    return (
        <div className="min-w-fit w-full h-fit bg-blue-900 flex flex-col justify-center items-center py-8 max-md:rounded-b-xl lg:rounded-br-xl">
            {uploadError && <NoticeModalComponent showModal={true} title="Error al subir el archivo" message="Ha ocurrido un error al intentar subir el archivo. Por favor, inténtelo de nuevo más tarde." />}
            {uploadSuccess && <NoticeModalComponent showModal={true} title="Archivo subido" message="El archivo se ha subido correctamente." />}
            <div className="w-full h-fit flex flex-col justify-center items-center gap-3 py-2 lg:flex-row">
                <h2 className="text-xl text-white font-semibold">Subir archivo</h2>
                <input 
                    onChange={handlerGetFile}
                    id={`file-${TypeReports.CATASTRE}`} 
                    type="file" className="w-48 text-md text-white" 
                />
            </div>
            <button
                onClick={handlerUpload}
                id={`upload-${TypeReports.CATASTRE}`} 
                className="w-fit h-fit bg-blue-700 text-white font-semibold rounded-full px-10 py-2 transition-all ease-linear duration-300 hover:bg-green-700"
            >
                Subir
            </button>
        </div>
    );
}