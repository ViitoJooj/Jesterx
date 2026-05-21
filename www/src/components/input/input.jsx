import styles from "./Input.module.scss";

export default function Input({ value, onChange, variant = "default", error = false, className = "", ...rest }) {
  const classes = [styles.input, styles[variant], error ? styles.error : "", className].join(" ");

  return (
    <input
      value={value}
      onChange={onChange}
      className={classes}
      {...rest}
    />
  );
}
