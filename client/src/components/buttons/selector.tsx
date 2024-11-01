import React from "react";

interface Props {
    id: any;
    label: string;
    // options: string[];
}

export default function Selector(content: Props) {
    return (
        <div className="w-full h-8 flex flex-row min-w-fit max-md:h-fit max-md:flex-col">
            <label id={content.id} className="w-fit h-full px-6 py-1 bg-blue-700 text-white font-bold max-md:w-full md:rounded-l-full max-md:rounded-t-xl">{content.label}</label>
            <select id={content.id} className="w-full h-full bg-gray-55 px-4 md:rounded-r-full max-sm:rounded-b-xl">
                <option value="">---</option>
            </select>
        </div>
    );
}