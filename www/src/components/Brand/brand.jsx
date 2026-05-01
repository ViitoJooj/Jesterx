import { Link } from "react-router-dom";
import styles from "./brand.module.scss";

export function Brand({ to, size = "sm" }) {
  const Tag = to ? Link : "div";
  const tagProps = to ? { to } : {};
  return (
    <Tag {...tagProps} className={`${styles.brand} ${styles[size]}`}>
      <span className={styles.logo}>J</span>
      <span className={styles.name}>Jester</span>
    </Tag>
  );
}
