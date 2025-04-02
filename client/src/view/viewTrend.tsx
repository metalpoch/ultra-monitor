import BasicSelectorComponent from '../components/selector/basic';
import LineGraphComponent from '../components/graphic/line';
import NoticeModalComponent from '../components/modal/notice';
import { DeviceController } from '../controllers/device';
import { TrendController } from '../controllers/trend'; 
import { LoadStatus } from "../constant/loadStatus";
import type { TrendSchema } from '../schemas/trend';
import type { LoadingStateValue } from '../types/loadingState';
import React, { useEffect, useState } from 'react';

export default function ViewTrend() {
    const [loadingData, setLoadingData] = useState<LoadingStateValue>(LoadStatus.EMPTY);
    const [error, setError] = useState<boolean>(false);
    const [olts, setOLTs] = useState<string[]>([]);
    const [trend, setTrend] = useState<TrendSchema | null>(null);
    
    const [oltSelect, setOltSelect] = useState<string>();

    /**
     * Get the data of all OLT availables.
     */
    const getData = async () => {
        let data = await DeviceController.getAllAcronyms();
        if (data && data.length > 0) setOLTs(data);
        else setError(true);
    };

    /**
     * Get traffic trend of OLT selected.
     */
    const getTrend = async () => {
        if (oltSelect) {
            setLoadingData(LoadStatus.LOADED);
            const data = await TrendController.getTrend(oltSelect);
            if (data) {
                setTrend(data);
                console.log(data);
            }
            else {
                setLoadingData(LoadStatus.LOADED);
                setError(true);
            }
        }
    }

    /**
     * Handler to get OLT selected.
     */
    const handlerOLTSelect = (oltSelect: string) => {
        setOltSelect(oltSelect);
    }


    useEffect(() => {
        setLoadingData(LoadStatus.EMPTY);
        getData();
    }, []);

    useEffect(() => {
        getTrend();
    }, [oltSelect])

    return(
        <main className='w-full flex flex-col gap-4 md:flex-row md:gap-1 md:justify-center md:max-h-[800px]'>
            {error && <NoticeModalComponent showModal={true} title='Error al obtener la predicción' message='No se pudo obtener la predicción del tráfico del OLT seleccionado. Por favor, inténtelo de nuevo más tarde.' onClick={() => setError(false)} />}
            <section id='filter' className='w-96 h-fit p-6 bg-gray-50 rounded-md flex flex-col justify-center gap-2'>
                <h1 className='text-xl font-bold text-blue-800'>Tendencia de</h1>
                <BasicSelectorComponent id='filter-selector' label='OLT' options={olts} onChange={handlerOLTSelect}/>
            </section>
            <LineGraphComponent title='Tendencia de Tráfico' canvasID='trend-graph' data={[]} loading={loadingData} />
        </main>
    )
}