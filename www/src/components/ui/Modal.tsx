import { useEffect, type ReactNode } from "react";
import { createPortal } from "react-dom";
import { IconButton } from "./Button";
import "./Modal.scss";

export interface ModalProps {
    open: boolean;
    onClose: () => void;
    title?: string;
    subtitle?: string;
    children: ReactNode;
}

export function Modal({ open, onClose, title, subtitle, children }: ModalProps) {
    useEffect(() => {
        if (!open) return;
        const onKey = (e: KeyboardEvent) => {
            if (e.key === "Escape") onClose();
        };
        window.addEventListener("keydown", onKey);
        return () => window.removeEventListener("keydown", onKey);
    }, [open, onClose]);

    if (!open) return null;

    return createPortal(
        <div className="jx-modal-overlay" onClick={onClose}>
            <div
                className="jx-modal"
                role="dialog"
                aria-modal="true"
                aria-label={title}
                onClick={(e) => e.stopPropagation()}
            >
                <div className="jx-modal__close">
                    <IconButton icon="x" label="Fechar" size="sm" onClick={onClose} />
                </div>
                {(title || subtitle) && (
                    <header className="jx-modal__header">
                        {title && <h2 className="jx-modal__title">{title}</h2>}
                        {subtitle && <p className="jx-modal__subtitle">{subtitle}</p>}
                    </header>
                )}
                <div className="jx-modal__body">{children}</div>
            </div>
        </div>,
        document.body
    );
}
