import { Link } from "react-router-dom";
import styles from "./button.module.scss";

export default function Button({
  children,
  variant = "secondary",
  to,
  className = "",
  ...rest
}) {
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