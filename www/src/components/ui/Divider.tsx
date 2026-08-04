import "./Divider.scss";

export interface DividerProps {
    label?: string;
}

export function Divider({ label }: DividerProps) {
    if (!label) return <hr className="jx-divider" />;
    return (
        <div className="jx-divider-labeled">
            <span className="jx-divider" />
            <span className="jx-divider-labeled__text">{label}</span>
            <span className="jx-divider" />
        </div>
    );
}
