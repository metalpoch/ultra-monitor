import React from "react";

export default function InputField({
  id,
  label,
  icon: Icon,
  type = "text",
  value,
  onChange,
  placeholder,
  className = "",
  ...props
}) {
  return (
    <div className="flex flex-col gap-1">
      <label htmlFor={id} className="text-gray-200">{label}</label>
      <section className={`flex items-center gap-2 border border-[hsl(217,33%,20%)] bg-[#0f1729] p-2 rounded-xs focus-within:ring-2 focus-within:ring-blue-500 ${className}`}>
        {Icon && <Icon className="w-5" />}
        <input
          className="w-full bg-[#f1729] focus:outline-none"
          id={id}
          type={type}
          value={value}
          onChange={onChange}
          placeholder={placeholder}
          required
          {...props}
        />
      </section>
    </div>
  );
}