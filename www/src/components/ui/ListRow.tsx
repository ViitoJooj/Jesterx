import type { ReactNode } from "react";
import { IconTile } from "./IconTile";
import type { IconName } from "../icons/Icon";
import "./ListRow.scss";

export interface ListRowProps {
    icon?: IconName;
    title: string;
    subtitle?: string;
    trailing?: ReactNode;
    onClick?: () => void;
}

export function ListRow({ icon, title, subtitle, trailing, onClick }: ListRowProps) {
    const Tag = onClick ? "button" : "div";
    return (
        <Tag
            className={`jx-list-row ${onClick ? "jx-list-row--clickable" : ""}`}
            onClick={onClick}
            type={onClick ? "button" : undefined}
        >
            {icon && <IconTile icon={icon} />}
            <span className="jx-list-row__text">
                <span className="jx-list-row__title">{title}</span>
                {subtitle && <span className="jx-list-row__subtitle">{subtitle}</span>}
            </span>
            {trailing && <span className="jx-list-row__trailing">{trailing}</span>}
        </Tag>
    );
}
