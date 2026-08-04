import { Icon, type IconName } from "../icons/Icon";
import "./IconTile.scss";

export interface IconTileProps {
    icon: IconName;
    size?: "sm" | "md" | "lg";
}

export function IconTile({ icon, size = "md" }: IconTileProps) {
    const iconSize = size === "sm" ? 16 : size === "lg" ? 24 : 20;
    return (
        <div className={`jx-icon-tile jx-icon-tile--${size}`}>
            <Icon name={icon} size={iconSize} />
        </div>
    );
}
