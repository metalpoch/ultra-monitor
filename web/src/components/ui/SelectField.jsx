export default function SelectField({
    id,
    label,
    options,
    value,
    onChange,
    className = "",
    ...props
}) {
    return (
        <div className="flex flex-col gap-1">
            <label htmlFor={id} className="text-gray-200">{label}</label>
            <select
                id={id}
                value={value}
                onChange={onChange}
                className={`w-full border border-[hsl(217,33%,20%)] bg-[#0f1729] text-gray-200 p-2 rounded-xs focus:outline-none focus:ring-2 focus:ring-blue-500 ${className}`}
                {...props}
            >
                {options.map((option) => (
                    <option key={option.value} value={option.value}>
                        {option.label}
                    </option>
                ))}
            </select>
        </div>
    );
}