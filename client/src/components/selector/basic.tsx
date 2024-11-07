import React, { useEffect } from "react";

interface Props {
    id: any;
    label: string;
    options: any[];
    onChange?: (newValue: any) => void;
}

export default function BasicSelectorComponent(content: Props) {

    const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
        if (content.onChange && event.target.value) content.onChange(event.target.value);
    }

    useEffect(() => {
        const select = document.getElementById(content.id) as HTMLSelectElement;
        content.options?.length > 0 ? select.disabled = false : select.disabled = true;
    }, [content.options]);

    return (
        <div className="w-full h-8 flex flex-row min-w-fit max-md:h-fit max-md:flex-col">
            <label className="w-fit h-full px-6 py-1 bg-blue-700 text-white font-bold max-md:w-full md:rounded-l-full max-md:rounded-t-xl">{content.label}</label>
            <select id={content.id} onChange={handleChange} className={`w-full h-full ${content.options?.length > 0 ? "bg-gray-100" : "bg-gray-300"} px-4 outline-none md:rounded-r-full max-sm:rounded-b-xl`} >
                <option>---</option>
                {content.options?.length > 0 && content.options.map((option, index) => (
                    <option key={index} value={option}>{option}</option>
                ))}
            </select>
        </div>
    );
}