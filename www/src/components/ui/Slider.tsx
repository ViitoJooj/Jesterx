import "./Slider.scss";

export interface SliderProps {
    value: number;
    onChange: (value: number) => void;
    min?: number;
    max?: number;
    step?: number;
    label?: string;
    marks?: string[];
}

export function Slider({ value, onChange, min = 0, max = 100, step = 1, label, marks }: SliderProps) {
    const percent = ((value - min) / (max - min)) * 100;
    return (
        <div className="jx-slider">
            {label && <span className="jx-slider__label">{label}</span>}
            <input
                type="range"
                className="jx-slider__input"
                min={min}
                max={max}
                step={step}
                value={value}
                onChange={(e) => onChange(Number(e.target.value))}
                style={{ "--fill": `${percent}%` } as React.CSSProperties}
            />
            {marks && (
                <div className="jx-slider__marks">
                    {marks.map((m) => (
                        <span key={m}>{m}</span>
                    ))}
                </div>
            )}
        </div>
    );
}
