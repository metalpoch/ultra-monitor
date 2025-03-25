import NoticeModalComponent from '../modal/notice.tsx';
import { TypeReports } from '../../constant/reports.ts';
import { ReportController } from '../../controllers/report';  
import type { User } from '../../schemas/user.ts';
import React, { useState} from 'react';

interface FileProps {
    user: User;
    type: string;
}


export default function FileInputComponent(props: FileProps) {
    const [file, setFile] = useState<any | null>(null);
    const [uploadSuccess, SetUploadSuccess] = useState(false);
    const [uploadError, SetUploadError] = useState(false);

    const sendReport = async () => {
        const response = await ReportController.uploadFile(file, props.user.p00);
        if (response) SetUploadSuccess(true);
        else SetUploadError(true);
    }

    const handlerFile = (event: any) => {
        let file = event.target.files[0];
        if (file) setFile(file);
    }

    const handlerNoticeModal = () => {
        SetUploadSuccess(false);
        SetUploadError(false);
    }

    return (
        <section className='w-full'>
            {uploadSuccess && <NoticeModalComponent showModal={uploadSuccess} title="Archivo subido" message="El archivo se ha subido correctamente." onClick={handlerNoticeModal} />}
            {uploadError && <NoticeModalComponent showModal={uploadError} title="Error" message="Ha ocurrido un error al subir el archivo." onClick={handlerNoticeModal} />}
            <div className="min-w-fit w-full h-fit bg-blue-900 flex flex-col justify-center items-center py-8 max-md:rounded-b-xl lg:rounded-br-xl">
                <div className="w-full h-fit flex flex-col justify-center items-center gap-3 py-2 lg:flex-row">
                    <h2 className="text-xl text-white font-semibold">Subir archivo</h2>
                    <input id={props.type} onChange={handlerFile} type="file" className="w-48 text-md text-white" />
                </div>
                <button 
                    id={`upload-${props.type}`}
                    onClick={sendReport}
                    className="w-fit h-fit bg-blue-700 text-white font-semibold rounded-full px-10 py-2 transition-all ease-linear duration-300 hover:bg-green-700"
                >
                    Subir
                </button>
            </div>
        </section>
    );
}