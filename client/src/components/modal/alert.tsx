import React, { useEffect } from 'react';

interface Props {
    showModal: boolean;
    title: string;
    message: string;

}

export default function AlertModalComponent({ showModal, title, message }: Props) {

    const handlerAccept = () => {
        document.getElementById('modal-state')?.classList.add('hidden');
    }

    const handlerCancel = () => {
        document.getElementById('modal-state')?.classList.add('hidden');
    }

    useEffect(() => {
        const modalAccept = document.getElementById('modal-accept');
        const modalCancel = document.getElementById('modal-cancel');
        const modalState = document.getElementById('modal-state');

        if (modalAccept && modalCancel && modalState) {
            modalAccept.addEventListener('click', handlerAccept);
            modalCancel.addEventListener('click', handlerCancel);

            if (showModal) modalState.classList.remove('hidden');
            else modalState.classList.add('hidden');
            
        }

    }, [showModal]);

    return(
        <div id="modal-state" className="absolute z-10" aria-labelledby="modal-title" role="dialog" aria-modal="true">
            <aside className="fixed inset-0 bg-black bg-opacity-55 transition-opacity" aria-hidden="true"></aside>
            <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
                <div id="modal-panel" className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <div className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg">
                        <section className="bg-gray-50 px-4 pb-4 pt-5 sm:p-6 sm:pb-4">
                            <div className="sm:flex sm:items-start">
                                <div className="mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left">
                                    <h3 className="text-base font-semibold leading-6 text-gray-900" id="modal-title">{title}</h3>
                                    <div className="mt-2">
                                        <p className="text-sm text-gray-500">{message}</p>
                                    </div>
                                </div>
                            </div>
                        </section>
                        <section className="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6">
                            <button id="modal-accept" type="button" className="inline-flex w-full justify-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm transition-all ease-linear duration-200 hover:bg-blue-700 sm:ml-3 sm:w-auto">Aceptar</button>
                            <button id="modal-cancel" type="button" className="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 transition-all ease-linear duration-200 hover:bg-gray-500 hover:text-white sm:mt-0 sm:w-auto">Cancel</button>
                        </section>
                    </div>
                </div>
            </div>
        </div>
    );
}