import "./ProgressDots.scss";

export interface ProgressDotsProps {
    steps: number;
    current: number;
}

export function ProgressDots({ steps, current }: ProgressDotsProps) {
    return (
        <div className="jx-progress-dots" role="progressbar" aria-valuenow={current} aria-valuemax={steps}>
            {Array.from({ length: steps }, (_, i) => (
                <span
                    key={i}
                    className={`jx-progress-dots__dot ${
                        i < current ? "jx-progress-dots__dot--done" : ""
                    } ${i === current ? "jx-progress-dots__dot--current" : ""}`}
                />
            ))}
        </div>
    );
}
