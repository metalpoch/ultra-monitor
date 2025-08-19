export default function DatalistField({
  id,
  label,
  options,
  value,
  onChange,
  placeholder,
  className = "",
  ...props
}) {
  return (
    <div className="flex flex-col gap-1">
      <label htmlFor={id} className="text-gray-200">
        {label}
      </label>
      <section
        className={`flex items-center gap-2 border border-[hsl(217,33%,20%)] bg-[#0f1729] p-2 rounded-xs focus-within:ring-2 focus-within:ring-blue-500 ${className}`}
      >
        <input
          className="w-full bg-[#0f1729] text-gray-200 focus:outline-none"
          id={id}
          list={`${id}-datalist`}
          value={value}
          onChange={onChange}
          placeholder={placeholder}
          autoComplete="off"
          {...props}
        />
        <datalist id={`${id}-datalist`}>
          {options.map((option) =>
            typeof option === "string" ? (
              <option key={option} value={option} />
            ) : (
              <option key={`${option.label}-${option.value}`} value={option.value}>
                {option.label}
              </option>
            )
          )}
        </datalist>
      </section>
    </div>
  );
}
