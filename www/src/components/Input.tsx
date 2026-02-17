import React from "react";

type InputProps = Omit<React.InputHTMLAttributes<HTMLInputElement>, "onChange" | "value"> & {
  value?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
};

const Input: React.FC<InputProps> = ({ value, onChange, className = "", ...rest }) => {
  return <input value={value} onChange={onChange} className={className} {...rest} />;
};

export default Input;