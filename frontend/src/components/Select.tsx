import { useState, useRef, useEffect } from "react";
import styles from "../styles/components/Select.module.scss";

type Option = { label: string; value: string };

type CustomSelectProps = {
    options: Option[];
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
};

export function CustomSelect({ options, value, onChange, placeholder }: CustomSelectProps) {
    const [open, setOpen] = useState(false);
    const ref = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (ref.current && !ref.current.contains(e.target as Node)) {
                setOpen(false);
            }
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    const selectedLabel = options.find((o) => o.value === value)?.label;

    return (
        <div className={styles.selectWrapper} ref={ref}>
            <button
                type="button"
                className={styles.selectButton}
                onClick={() => setOpen((prev) => !prev)}
            >
                {selectedLabel || placeholder || "Selecione..."}
                <span className={styles.arrow}>▼</span>
            </button>

            {open && (
                <ul className={styles.selectDropdown}>
                    {options.map((opt) => (
                        <li
                            key={opt.value}
                            className={`${styles.selectOption} ${opt.value === value ? styles.active : ""}`}
                            onClick={() => {
                                onChange(opt.value);
                                setOpen(false);
                            }}
                        >
                            {opt.label}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}
