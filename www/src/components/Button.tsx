import React from "react";

type ButtonProps = {
  label: string;
  onClick?: () => void;
  type?: "button" | "submit";
  disabled?: boolean;
  className?: string;
};

const Button: React.FC<ButtonProps> = ({ label, onClick, type = "button", disabled = false, className = "" }) => {
  return (
    <button type={type} onClick={onClick} disabled={disabled} className={`${className}`}>
      {label}
    </button>
  );
};

export default Button;