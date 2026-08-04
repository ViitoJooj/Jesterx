import type { ReactNode } from "react";
import "./Tooltip.scss";

export interface TooltipProps {
    content: string;
    children: ReactNode;
    position?: "top" | "bottom";
}

export function Tooltip({ content, children, position = "top" }: TooltipProps) {
    return (
        <span className={`jx-tooltip jx-tooltip--${position}`}>
            {children}
            <span className="jx-tooltip__bubble" role="tooltip">
                {content}
            </span>
        </span>
    );
}
