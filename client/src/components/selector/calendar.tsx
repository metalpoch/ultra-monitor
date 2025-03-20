import React from "react";

/**
 * @interface Data required for the calendar selector.
 * 
 * @param {any} id ID of the selector.
 * @param {string} label Label of the selector.
 * @param {(newDate: string) => void} onChange Callback to handle the change of the selector.
 */
interface CalendarProps {
    id: any;
    label: string;
    onChange?: (newDate: string) => void;
}

export default function CalendarSelectorComponent(content: CalendarProps) {

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        if (content.onChange) content.onChange(event.target.value);
    }

    return (
        <div className="h-8 flex flex-row min-w-fit max-md:w-full max-md:h-fit max-md:flex-col">
            <label className="w-fit h-full px-4 py-1 bg-blue-700 text-white font-bold max-md:w-full md:rounded-l-full max-md:rounded-t-xl">{content.label}</label>
            <input id={content.id} type="date" onChange={handleChange} className="w-full h-full bg-gray-100 px-2 outline-none max-sm:w-full md:rounded-r-full max-sm:rounded-b-xl" />
        </div>
    );
}