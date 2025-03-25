import React, { useEffect } from 'react';

/**
 * @interface Data required for the notice modal.
 * 
 * @param {boolean} showModal if the modal is visible.
 * @param {string} title Title of the modal.
 * @param {string} message Message of the modal.
 */
interface ModalProps {
    showModal: boolean;
    title: string;
    message: string;
    onClick?: () => void;
}

export default function NoticeModalComponent(props: ModalProps) {

    const handlerAccept = () => {
        document.getElementById('modal-state')?.classList.add('hidden');
        if (props.onClick) props.onClick();
    }

    useEffect(() => {
        const modalAccept = document.getElementById('modal-accept');
        const modalState = document.getElementById('modal-state');

        if (modalAccept && modalState) {
            modalAccept.addEventListener('click', handlerAccept);

            if (props.showModal) modalState.classList.remove('hidden');
            else modalState.classList.add('hidden');
        }

    }, [props.showModal]);

    return(
        <div id="modal-state" className="absolute z-10" aria-labelledby="modal-title" role="dialog" aria-modal="true">
            <aside className="fixed inset-0 bg-black bg-opacity-55 transition-opacity" aria-hidden="true"></aside>
            <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
                <div id="modal-panel" className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <div className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg">
                        <section className="bg-gray-50 px-4 pb-4 pt-5 sm:p-6 sm:pb-4">
                            <div className="sm:flex sm:items-start">
                                <div className="mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left">
                                    <h3 className="text-base font-semibold leading-6 text-gray-900" id="modal-title">{props.title}</h3>
                                    <div className="mt-2">
                                        <p className="text-sm text-gray-500">{props.message}</p>
                                    </div>
                                </div>
                            </div>
                        </section>
                        <section className="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6">
                            <button 
                                id="modal-accept" 
                                type="button" 
                                className="inline-flex w-full justify-center rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm transition-all ease-linear duration-200 hover:bg-blue-700 sm:ml-3 sm:w-auto"
                            >
                                Aceptar
                            </button>
                        </section>
                    </div>
                </div>
            </div>
        </div>
    );
}