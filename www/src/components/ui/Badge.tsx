import type { ReactNode } from "react";
import "./Badge.scss";

export interface BadgeProps {
    children: ReactNode;
    variant?: "glass" | "solid";
}

export function Badge({ children, variant = "glass" }: BadgeProps) {
    return <span className={`jx-badge jx-badge--${variant}`}>{children}</span>;
}
