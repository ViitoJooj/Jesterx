import type { ButtonHTMLAttributes, ReactNode } from "react";
import { Icon, type IconName } from "../icons/Icon";
import "./Button.scss";

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: "primary" | "glass" | "ghost";
    size?: "sm" | "md" | "lg";
    icon?: IconName;
    iconRight?: IconName;
    children?: ReactNode;
}

export function Button({
    variant = "glass",
    size = "md",
    icon,
    iconRight,
    children,
    className = "",
    ...rest
}: ButtonProps) {
    return (
        <button
            className={`jx-btn jx-btn--${variant} jx-btn--${size} ${className}`}
            {...rest}
        >
            {icon && <Icon name={icon} size={size === "sm" ? 14 : size === "lg" ? 20 : 16} />}
            {children && <span>{children}</span>}
            {iconRight && <Icon name={iconRight} size={size === "sm" ? 14 : size === "lg" ? 20 : 16} />}
        </button>
    );
}

export interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    icon: IconName;
    size?: "sm" | "md" | "lg";
    label: string;
}

export function IconButton({ icon, size = "md", label, className = "", ...rest }: IconButtonProps) {
    return (
        <button
            className={`jx-icon-btn jx-icon-btn--${size} ${className}`}
            aria-label={label}
            title={label}
            {...rest}
        >
            <Icon name={icon} size={size === "sm" ? 14 : size === "lg" ? 22 : 18} />
        </button>
    );
}
