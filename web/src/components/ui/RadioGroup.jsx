export default function RadioGroup({
  id,
  label,
  options = [],
  value,
  onChange,
  className = "",
  ...props
}) {
  return (
    <div className="flex flex-col gap-1">
      <label htmlFor={id} className="text-gray-200">
        {label}
      </label>
      <section
        className={`flex flex-col gap-2 border border-[hsl(217,33%,20%)] bg-[#0f1729] p-2 rounded-lg focus-within:ring-2 focus-within:ring-blue-500 ${className}`}
      >
        {options.map((option) => (
          <label
            key={option.value}
            className="flex items-center gap-2 text-gray-200 cursor-pointer"
          >
            <input
              type="radio"
              name={id}
              value={option.value}
              checked={value === option.value}
              onChange={(e) => onChange(e.target.value)}
              className="accent-blue-500 focus:ring-2 focus:ring-blue-500"
              {...props}
            />
            <span>{option.label}</span>
          </label>
        ))}
      </section>
    </div>
  );
}

