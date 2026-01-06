import { Link } from "react-router-dom";
import styles from "../styles/pages/NotFound.module.scss";
import buttonStyles from "../styles/components/Button.module.scss";

export function NotFound() {
  return (
    <main className={styles.wrap} aria-labelledby="nf-title">
      <section className={styles.card}>
        <div className={styles.code}>404</div>
        <h1 id="nf-title" className={styles.title}>
          Página não encontrada
        </h1>
        <p className={styles.desc}>A página que você procura não existe ou foi movida.</p>
        <div className={styles.actions}>
          <Link to="/" className={`${buttonStyles.default_button} ${buttonStyles["default_button--primary"]}`}>
            Ir para início
          </Link>
          <Link to="/pages" className={`${buttonStyles.default_button} ${buttonStyles["default_button--secondary"]}`}>
            Minhas páginas
          </Link>
        </div>
      </section>
    </main>
  );
}
