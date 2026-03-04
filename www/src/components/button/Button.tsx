import { Link } from "react-router-dom";
import React from "react";
import styles from "./Button.module.scss";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "ghost";
  to?: string;
};

export default function Button({
  children,
  variant = "secondary",
  to,
  className = "",
  ...rest
}: ButtonProps) {

  const classes = [
    styles.default_button,
    styles[`default_button--${variant}`],
    className,
  ].join(" ");

  if (to) {
    return (
      <Link to={to} className={classes}>
        {children}
      </Link>
    );
  }

  return (
    <button className={classes} {...rest}>
      {children}
    </button>
  );
}