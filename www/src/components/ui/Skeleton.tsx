import "./Skeleton.scss";

export interface SkeletonProps {
    width?: string | number;
    height?: string | number;
    circle?: boolean;
}

export function Skeleton({ width = "100%", height = 16, circle }: SkeletonProps) {
    const w = typeof width === "number" ? `${width}px` : width;
    const h = typeof height === "number" ? `${height}px` : height;
    return (
        <span
            className={`jx-skeleton ${circle ? "jx-skeleton--circle" : ""}`}
            style={{ width: w, height: h }}
        />
    );
}
