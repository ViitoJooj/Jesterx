import "./SegmentedControl.scss";

export interface SegmentedOption {
    value: string;
    label: string;
}

export interface SegmentedControlProps {
    options: SegmentedOption[];
    value: string;
    onChange: (value: string) => void;
}

export function SegmentedControl({ options, value, onChange }: SegmentedControlProps) {
    return (
        <div className="jx-segmented" role="tablist">
            {options.map((opt) => (
                <button
                    key={opt.value}
                    role="tab"
                    aria-selected={opt.value === value}
                    className={`jx-segmented__item ${
                        opt.value === value ? "jx-segmented__item--active" : ""
                    }`}
                    onClick={() => onChange(opt.value)}
                >
                    {opt.label}
                </button>
            ))}
        </div>
    );
}
