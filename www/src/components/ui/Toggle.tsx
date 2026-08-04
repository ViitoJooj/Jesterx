import "./Toggle.scss";

export interface ToggleProps {
    checked: boolean;
    onChange: (checked: boolean) => void;
    disabled?: boolean;
    label?: string;
}

export function Toggle({ checked, onChange, disabled, label }: ToggleProps) {
    return (
        <button
            type="button"
            role="switch"
            aria-checked={checked}
            aria-label={label}
            disabled={disabled}
            className={`jx-toggle ${checked ? "jx-toggle--on" : ""}`}
            onClick={() => onChange(!checked)}
        >
            <span className="jx-toggle__thumb" />
        </button>
    );
}
