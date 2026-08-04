import { useEffect, useState } from "react";

export type Theme = "dark" | "light";

const STORAGE_KEY = "theme";

function detectTheme(): Theme {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved === "dark" || saved === "light") return saved;
    return window.matchMedia("(prefers-color-scheme: light)").matches
        ? "light"
        : "dark";
}

export function useTheme() {
    const [theme, setTheme] = useState<Theme>(detectTheme);

    useEffect(() => {
        document.documentElement.dataset.theme = theme;
        localStorage.setItem(STORAGE_KEY, theme);
    }, [theme]);

    const toggleTheme = () =>
        setTheme((prev) => (prev === "dark" ? "light" : "dark"));

    return { theme, toggleTheme };
}
