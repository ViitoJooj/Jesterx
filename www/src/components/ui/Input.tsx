import { forwardRef, type InputHTMLAttributes } from "react";
import { Icon, type IconName } from "../icons/Icon";
import "./Input.scss";

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
    label?: string;
    icon?: IconName;
    error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
    ({ label, icon, error, className = "", id, ...rest }, ref) => {
        const inputId = id ?? label?.toLowerCase().replace(/\s+/g, "-");
        return (
            <div className={`jx-field ${error ? "jx-field--error" : ""} ${className}`}>
                {label && (
                    <label className="jx-field__label" htmlFor={inputId}>
                        {label}
                    </label>
                )}
                <div className="jx-input">
                    {icon && (
                        <span className="jx-input__icon">
                            <Icon name={icon} size={16} />
                        </span>
                    )}
                    <input
                        ref={ref}
                        id={inputId}
                        className={`jx-input__control ${icon ? "jx-input__control--icon" : ""}`}
                        {...rest}
                    />
                </div>
                {error && <span className="jx-field__error">{error}</span>}
            </div>
        );
    }
);

Input.displayName = "Input";
