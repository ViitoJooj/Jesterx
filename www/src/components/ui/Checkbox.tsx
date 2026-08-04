import { Icon } from "../icons/Icon";
import "./Checkbox.scss";

export interface CheckboxProps {
    checked: boolean;
    onChange: (checked: boolean) => void;
    label?: string;
    disabled?: boolean;
}

export function Checkbox({ checked, onChange, label, disabled }: CheckboxProps) {
    return (
        <label className={`jx-checkbox ${disabled ? "jx-checkbox--disabled" : ""}`}>
            <input
                type="checkbox"
                className="jx-checkbox__input"
                checked={checked}
                disabled={disabled}
                onChange={(e) => onChange(e.target.checked)}
            />
            <span className="jx-checkbox__box">
                {checked && <Icon name="check" size={14} strokeWidth={3} />}
            </span>
            {label && <span className="jx-checkbox__label">{label}</span>}
        </label>
    );
}

export interface RadioProps {
    checked: boolean;
    onChange: () => void;
    label?: string;
    disabled?: boolean;
    name?: string;
}

export function Radio({ checked, onChange, label, disabled, name }: RadioProps) {
    return (
        <label className={`jx-radio ${disabled ? "jx-radio--disabled" : ""}`}>
            <input
                type="radio"
                className="jx-radio__input"
                checked={checked}
                disabled={disabled}
                name={name}
                onChange={onChange}
            />
            <span className="jx-radio__circle">{checked && <span className="jx-radio__dot" />}</span>
            {label && <span className="jx-radio__label">{label}</span>}
        </label>
    );
}
