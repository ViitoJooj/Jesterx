import { forwardRef, type TextareaHTMLAttributes } from "react";
import "./Textarea.scss";

export interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
    label?: string;
    error?: string;
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
    ({ label, error, className = "", id, ...rest }, ref) => {
        const inputId = id ?? label?.toLowerCase().replace(/\s+/g, "-");
        return (
            <div className={`jx-field ${error ? "jx-field--error" : ""} ${className}`}>
                {label && (
                    <label className="jx-field__label" htmlFor={inputId}>
                        {label}
                    </label>
                )}
                <textarea ref={ref} id={inputId} className="jx-textarea" {...rest} />
                {error && <span className="jx-field__error">{error}</span>}
            </div>
        );
    }
);

Textarea.displayName = "Textarea";
