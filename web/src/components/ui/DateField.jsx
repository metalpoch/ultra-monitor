import dayjs from "dayjs";
import { useState } from "react";

export default function DateField({
  id,
  label,
  value,
  onChange,
  endOfDay = false,
  min = "2025-07-23",
  max = new Date().toISOString().split("T")[0],
  className = "",
  ...props
}) {
  const [dateValue, setDateValue] = useState(value ? value.split("T")[0] : "");

  const handleChange = (e) => {
    const selectedDate = e.target.value;
    setDateValue(selectedDate);
    onChange && onChange(dayjs(selectedDate).toJSON());
  };

  return (
    <div className="flex flex-col gap-1">
      <label htmlFor={id} className="text-gray-200">{label}</label>
      <input
        id={id}
        type="date"
        value={dateValue}
        onChange={handleChange}
        min={min}
        max={max}
        className={`w-full border border-[hsl(217,33%,20%)] bg-[#0f1729] text-gray-200 p-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${className}`}
        {...props}
      />
    </div>
  );
}
