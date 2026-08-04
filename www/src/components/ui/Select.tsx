import { forwardRef, type SelectHTMLAttributes } from "react";
import { Icon } from "../icons/Icon";
import "./Select.scss";

export interface SelectOption {
    value: string;
    label: string;
}

export interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
    label?: string;
    options: SelectOption[];
    error?: string;
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(
    ({ label, options, error, className = "", id, ...rest }, ref) => {
        const selectId = id ?? label?.toLowerCase().replace(/\s+/g, "-");
        return (
            <div className={`jx-field ${error ? "jx-field--error" : ""} ${className}`}>
                {label && (
                    <label className="jx-field__label" htmlFor={selectId}>
                        {label}
                    </label>
                )}
                <div className="jx-select">
                    <select ref={ref} id={selectId} className="jx-select__control" {...rest}>
                        {options.map((opt) => (
                            <option key={opt.value} value={opt.value}>
                                {opt.label}
                            </option>
                        ))}
                    </select>
                    <span className="jx-select__chevron">
                        <Icon name="chevron-down" size={16} />
                    </span>
                </div>
                {error && <span className="jx-field__error">{error}</span>}
            </div>
        );
    }
);

Select.displayName = "Select";
