import type { ButtonHTMLAttributes } from 'react';

type ButtonVariant = 'primary' | 'secondary' | 'accent' | 'ghost' | 'outline' | 'glass';
type ButtonSize = 'sm' | 'md' | 'lg';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  label?: string;
  icon?: React.ReactNode;
  isFullWidth?: boolean;
}

export function Button({
  variant = 'primary',
  size = 'md',
  label,
  icon,
  isFullWidth,
  className = '',
  children,
  ...props
}: ButtonProps) {
  const baseClasses = `btn btn-${variant} btn-${size}`;
  const widthClass = isFullWidth ? 'w-full' : '';

  return (
    <button className={`${baseClasses} ${widthClass} ${className}`} {...props}>
      {icon && <span className="mr-2">{icon}</span>}
      {label || children}
    </button>
  );
}
