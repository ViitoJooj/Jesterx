import { useEffect } from "react";
import { createPortal } from "react-dom";
import { Icon, type IconName } from "../icons/Icon";
import "./Toast.scss";

export interface ToastProps {
    open: boolean;
    message: string;
    icon?: IconName;
    duration?: number;
    onClose: () => void;
}

export function Toast({ open, message, icon = "check", duration = 3000, onClose }: ToastProps) {
    useEffect(() => {
        if (!open) return;
        const t = setTimeout(onClose, duration);
        return () => clearTimeout(t);
    }, [open, duration, onClose]);

    if (!open) return null;

    return createPortal(
        <div className="jx-toast" role="status">
            <span className="jx-toast__icon">
                <Icon name={icon} size={16} />
            </span>
            {message}
        </div>,
        document.body
    );
}
