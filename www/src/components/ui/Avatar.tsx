import { Icon } from "../icons/Icon";
import "./Avatar.scss";

export interface AvatarProps {
    src?: string;
    name?: string;
    size?: "sm" | "md" | "lg";
    status?: "online" | "offline";
}

function initials(name: string): string {
    return name
        .split(" ")
        .filter(Boolean)
        .slice(0, 2)
        .map((p) => p[0]!.toUpperCase())
        .join("");
}

export function Avatar({ src, name, size = "md", status }: AvatarProps) {
    return (
        <span className={`jx-avatar jx-avatar--${size}`}>
            {src ? (
                <img className="jx-avatar__img" src={src} alt={name ?? "avatar"} />
            ) : name ? (
                <span className="jx-avatar__initials">{initials(name)}</span>
            ) : (
                <Icon name="user" size={size === "sm" ? 14 : size === "lg" ? 24 : 18} />
            )}
            {status && <span className={`jx-avatar__status jx-avatar__status--${status}`} />}
        </span>
    );
}
