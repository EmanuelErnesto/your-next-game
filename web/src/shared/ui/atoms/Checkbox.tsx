import type { InputHTMLAttributes } from 'react';

interface CheckboxProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
}

export function Checkbox({ label, className = '', ...props }: CheckboxProps) {
  return (
    <label className={`flex items-center gap-3 cursor-pointer group ${className}`}>
      <input
        type="checkbox"
        className="checkbox checkbox-sm checkbox-primary rounded-sm"
        {...props}
      />
      <span className="text-sm text-gray-400 group-hover:text-gray-200">{label}</span>
    </label>
  );
}
