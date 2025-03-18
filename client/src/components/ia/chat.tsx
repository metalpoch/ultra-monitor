import { answerIA } from '../../controllers/ia';
import { questionIA } from '../../controllers/ia';
import { useStore } from '@nanostores/react';
import React, { useState, useEffect } from 'react';

/**
 * @interface Data required for the IA chat.
 * 
 * @param {string | null} question Question of the IA.
 * @param {string | null} answer Answer of the IA.
 */
interface DataIAProps {
    question: string | null;
    answer: string | null;
}

export default function ChatIAComponent() {
    const $answerIA = useStore(answerIA);
    const $questionIA = useStore(questionIA);

    const [data, setData] = useState<DataIAProps[]>([]);

    useEffect(() => {
        let newData: DataIAProps = { question: $questionIA, answer: $answerIA };
        setData([...data, newData]);
    }, [$answerIA]);

    return (<>
        {data && data.length <= 1 && 
            <div className="w-full h-full bg-gray-55 rounded-xl shadow-xl flex flex-col items-center justify-center">
                <h2 className="text-2xl text-blue-700 font-bold">Hola, usuario</h2>
            </div>
        }
        {data && data.length > 1 &&
            <div className="w-full h-full max-h-full overflow-y-auto bg-gray-55 rounded-xl shadow-xl flex flex-col px-6 py-8">
                {data.map((consult: DataIAProps, index: number) => (
                    <section key={index}>
                        {consult.question && <h2 className="text-2xl text-gray-700 font-bold">{consult.question}</h2>}
                        {consult.question && consult.answer && <p className='text-base text-gray-700 px-4'>{consult.answer}</p>}
                        {consult.question && !consult.answer && <p className='text-base text-gray-700 px-4'>Oh no! No se pudo obtener respuesta a tu pregunta</p>}
                    </section>
                ))}
            </div>
        }
    </>);
}