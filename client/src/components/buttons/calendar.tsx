import React from "react";

interface Props {
    id: any;
    label: string;
}

export default function CalendarSelector(content: Props) {
    return (
        <div className="h-8 flex flex-row max-sm:w-full max-sm:h-fit max-sm:flex-col">
            <label id={content.id} className="w-fit h-full px-6 py-1 bg-blue-700 text-white font-bold max-sm:w-full md:rounded-l-full max-sm:rounded-t-xl">{content.label}</label>
            <input id={content.id} type="date" className="w-full h-full bg-gray-55 px-2 max-sm:w-full md:rounded-r-full max-sm:rounded-b-xl" />
        </div>
    );
}