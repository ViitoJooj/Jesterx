import React from "react";
import styles from "./Input.module.scss";

type InputProps = Omit<
  React.InputHTMLAttributes<HTMLInputElement>,
  "onChange" | "value"
> & {
  value?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  variant?: "default" | "primary";
  error?: boolean;
};

const Input: React.FC<InputProps> = ({
  value, onChange, variant = "default", error = false, className = "", ...rest}) => {
  const classes = [styles.input, styles[variant],  error ? styles.error : "", className].join(" ");

  return (
    <input
      value={value}
      onChange={onChange}
      className={classes}
      {...rest}
    />
  );
};

export default Input;