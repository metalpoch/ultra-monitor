import SpinnerBasicComponent from '../spinner/basic';
import { answerIA, loadingIA } from '../../controllers/ia';
import { questionIA } from '../../controllers/ia';
import { useStore } from '@nanostores/react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
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
    const $loadingIA = useStore(loadingIA);


    const [loading, setLoading] = useState(false);
    const [data, setData] = useState<DataIAProps[]>([]);

    const handlerLoading = (status: boolean = false) => {
        setLoading(status);
    }

    useEffect(() => {
        let newData: DataIAProps = { question: $questionIA, answer: $answerIA };
        setData([...data, newData]);
    }, [$answerIA, $questionIA]);

    useEffect(() => {
        handlerLoading($loadingIA);
    }, [$loadingIA]);

    return (<>
        {!loading && data && data.length <= 1 && 
            <div className="w-full h-full bg-gray-55 rounded-xl shadow-xl flex flex-col items-center justify-center drop-shadow-[1px_1px_2px_rgba(0,0,0,0.25)]">
                <h2 className="text-2xl text-blue-700 font-bold">Hola, usuario</h2>
            </div>
        }
        {!loading && data && data.length > 1 && 
            <div className="w-full h-full max-h-full overflow-y-auto bg-gray-55 rounded-xl shadow-xl flex flex-col px-6 py-8">
                {data.map((consult: DataIAProps, index: number) => (
                    <section key={index}>
                        {consult.question && consult.answer && <h2 className="text-2xl text-gray-700 font-bold">{consult.question}</h2>}
                        {consult.question && consult.answer && <ReactMarkdown remarkPlugins={[remarkGfm]}>{consult.answer}</ReactMarkdown>}
                    </section>
                ))}
            </div>
        }
        {loading &&
            <div className="w-full h-full bg-gray-55 rounded-xl shadow-xl flex flex-col items-center justify-center">
                <SpinnerBasicComponent />
            </div>
        }
    </>);
}