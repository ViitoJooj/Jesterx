import type { ReactNode } from "react";
import "./GlassCard.scss";

export interface GlassCardProps {
    children: ReactNode;
    className?: string;
    padding?: "sm" | "md" | "lg";
}

export function GlassCard({ children, className = "", padding = "md" }: GlassCardProps) {
    return (
        <div className={`jx-glass-card jx-glass-card--${padding} ${className}`}>
            {children}
        </div>
    );
}
