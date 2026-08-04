import type { CSSProperties } from "react";

export type IconName =
    | "moon"
    | "cloud"
    | "bell"
    | "lock"
    | "wifi"
    | "globe"
    | "mic"
    | "video"
    | "monitor"
    | "arrow-right"
    | "x"
    | "search"
    | "check"
    | "chevron-down"
    | "user"
    | "users"
    | "shield"
    | "palette"
    | "link"
    | "card"
    | "keyboard"
    | "wrench"
    | "sparkles"
    | "pen"
    | "coins"
    | "handshake"
    | "copy"
    | "sun"
    | "settings"
    | "log-out"
    | "log-in"
    | "user-plus"
    | "languages"
    | "facebook"
    | "twitter"
    | "instagram"
    | "linkedin";

const paths: Record<IconName, React.ReactNode> = {
    moon: <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79Z" />,
    cloud: (
        <>
            <path d="M17.5 19a4.5 4.5 0 0 0 .42-8.98 6 6 0 0 0-11.7 1.62A4 4 0 0 0 7 19Z" />
            <path d="M12 12v6" />
            <path d="m9.5 15.5 2.5 2.5 2.5-2.5" />
        </>
    ),
    bell: (
        <>
            <path d="M18 8a6 6 0 1 0-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
            <path d="M13.73 21a2 2 0 0 1-3.46 0" />
        </>
    ),
    lock: (
        <>
            <rect x="3" y="11" width="18" height="11" rx="2" />
            <path d="M7 11V7a5 5 0 0 1 10 0v4" />
        </>
    ),
    wifi: (
        <>
            <path d="M5 12.55a11 11 0 0 1 14.08 0" />
            <path d="M8.53 16.11a6 6 0 0 1 6.95 0" />
            <circle cx="12" cy="20" r="1" fill="currentColor" stroke="none" />
        </>
    ),
    globe: (
        <>
            <circle cx="12" cy="12" r="10" />
            <path d="M2 12h20" />
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10Z" />
        </>
    ),
    mic: (
        <>
            <rect x="9" y="2" width="6" height="12" rx="3" />
            <path d="M5 10a7 7 0 0 0 14 0" />
            <path d="M12 17v4" />
        </>
    ),
    video: (
        <>
            <rect x="2" y="6" width="13" height="12" rx="2" />
            <path d="m15 10 7-4v12l-7-4" />
        </>
    ),
    monitor: (
        <>
            <rect x="2" y="4" width="20" height="13" rx="2" />
            <path d="M8 21h8" />
            <path d="M12 17v4" />
        </>
    ),
    "arrow-right": (
        <>
            <path d="M5 12h14" />
            <path d="m13 6 6 6-6 6" />
        </>
    ),
    x: (
        <>
            <path d="M18 6 6 18" />
            <path d="m6 6 12 12" />
        </>
    ),
    search: (
        <>
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.35-4.35" />
        </>
    ),
    check: <path d="m4 12.5 5 5L20 6.5" />,
    "chevron-down": <path d="m6 9 6 6 6-6" />,
    user: (
        <>
            <circle cx="12" cy="8" r="4" />
            <path d="M4 21a8 8 0 0 1 16 0" />
        </>
    ),
    users: (
        <>
            <circle cx="9" cy="8" r="4" />
            <path d="M2 21a7 7 0 0 1 14 0" />
            <path d="M16 3.5a4 4 0 0 1 0 9" />
            <path d="M17.5 14.5A7 7 0 0 1 22 21" />
        </>
    ),
    shield: <path d="M12 22s8-3.5 8-10V5l-8-3-8 3v7c0 6.5 8 10 8 10Z" />,
    palette: (
        <>
            <path d="M12 22a10 10 0 1 1 10-10c0 2.5-1.5 3.5-3 3.5h-2.5a2.5 2.5 0 0 0-1.8 4.2c.6.7.2 2.3-2.7 2.3Z" />
            <circle cx="7.5" cy="11.5" r="1" fill="currentColor" stroke="none" />
            <circle cx="10.5" cy="7.5" r="1" fill="currentColor" stroke="none" />
            <circle cx="15" cy="7" r="1" fill="currentColor" stroke="none" />
        </>
    ),
    link: (
        <>
            <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
            <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
        </>
    ),
    card: (
        <>
            <rect x="2" y="5" width="20" height="14" rx="2" />
            <path d="M2 10h20" />
        </>
    ),
    keyboard: (
        <>
            <rect x="2" y="6" width="20" height="12" rx="2" />
            <path d="M6 10h.01M10 10h.01M14 10h.01M18 10h.01M6 14h.01M18 14h.01M9 14h6" />
        </>
    ),
    wrench: (
        <path d="M14.7 6.3a4.5 4.5 0 0 0-6 6L3 18l3 3 5.7-5.7a4.5 4.5 0 0 0 6-6L14 13l-3-3 3.7-3.7Z" />
    ),
    sparkles: (
        <>
            <path d="M12 3l1.9 5.1L19 10l-5.1 1.9L12 17l-1.9-5.1L5 10l5.1-1.9Z" />
            <path d="M19 15l.8 2.2L22 18l-2.2.8L19 21l-.8-2.2L16 18l2.2-.8Z" />
        </>
    ),
    pen: (
        <>
            <path d="M17 3a2.83 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z" />
        </>
    ),
    coins: (
        <>
            <circle cx="8" cy="8" r="6" />
            <path d="M18.09 10.37A6 6 0 1 1 10.34 18" />
            <path d="M7 6h1v4" />
            <path d="m16.71 13.88.7.71-2.82 2.82" />
        </>
    ),
    handshake: (
        <>
            <path d="m11 17 2 2a1.41 1.41 0 0 0 2-2l-2-2" />
            <path d="m14 14 2.5 2.5a1.41 1.41 0 0 0 2-2L15 11" />
            <path d="m3 11 4-4 5 2 5-2 4 4-2.5 2.5" />
            <path d="M8 12.5 5.5 15a1.41 1.41 0 0 0 2 2L10 14.5" />
        </>
    ),
    copy: (
        <>
            <rect x="9" y="9" width="13" height="13" rx="2" />
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
        </>
    ),
    sun: (
        <>
            <circle cx="12" cy="12" r="4" />
            <path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41" />
        </>
    ),
    settings: (
        <>
            <circle cx="12" cy="12" r="3" />
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1Z" />
        </>
    ),
    "log-out": (
        <>
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
            <path d="m16 17 5-5-5-5" />
            <path d="M21 12H9" />
        </>
    ),
    "log-in": (
        <>
            <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4" />
            <path d="m10 17 5-5-5-5" />
            <path d="M15 12H3" />
        </>
    ),
    "user-plus": (
        <>
            <circle cx="10" cy="8" r="4" />
            <path d="M2 21a8 8 0 0 1 16 0" />
            <path d="M19 5v6M22 8h-6" />
        </>
    ),
    languages: (
        <>
            <path d="m5 8 6 6" />
            <path d="m4 14 6-6 2-3" />
            <path d="M2 5h12" />
            <path d="M7 2h1" />
            <path d="m22 22-5-10-5 10" />
            <path d="M14 18h6" />
        </>
    ),
    facebook: <path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3Z" />,
    twitter: <path d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2Z" />,
    instagram: (
        <>
            <rect x="2" y="2" width="20" height="20" rx="5" />
            <path d="M16 11.37a4 4 0 1 1-7.75 1.26A4 4 0 0 1 16 11.37Z" />
            <path d="M17.5 6.5h.01" />
        </>
    ),
    linkedin: (
        <>
            <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-4 0v7h-4V8h4v1a6 6 0 0 1 2-1Z" />
            <rect x="2" y="9" width="4" height="12" />
            <circle cx="4" cy="4" r="2" />
        </>
    ),
};

export interface IconProps {
    name: IconName;
    size?: number;
    strokeWidth?: number;
    className?: string;
    style?: CSSProperties;
}

export function Icon({ name, size = 20, strokeWidth = 1.8, className, style }: IconProps) {
    return (
        <svg
            className={className}
            style={style}
            width={size}
            height={size}
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth={strokeWidth}
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden="true"
        >
            {paths[name]}
        </svg>
    );
}
