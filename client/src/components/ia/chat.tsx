import React, { useState, useEffect } from 'react';
import { useStore } from '@nanostores/react';
import { answerIA } from '../../controllers/ia';
import { questionIA } from '../../controllers/ia';

interface DataIA {
    question: string | null;
    answer: string | null;
}

export default function ChatIAComponent() {
    const $answerIA = useStore(answerIA);
    const $questionIA = useStore(questionIA);
    const [data, setData] = useState<DataIA[]>([]);

    useEffect(() => {
        let newData: DataIA = { question: $questionIA, answer: $answerIA };
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
                {data.map((item: DataIA, index: number) => (
                    <section key={index}>
                        {item.question && <h2 className="text-2xl text-gray-700 font-bold">{item.question}</h2>}
                        {item.question && item.answer && <p className='text-base text-gray-700 px-4'>{item.answer}</p>}
                        {item.question && !item.answer && <p className='text-base text-gray-700 px-4'>Oh no! No se pudo obtener respuesta a tu pregunta</p>}
                    </section>
                ))}
            </div>
        }
    </>);
}