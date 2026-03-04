import styles from "./NotFound.module.scss";
import Button from "../../components/button/Button";

export const NotFound: React.FC = () => {
  return (
    <main className={styles.wrap} aria-labelledby="nf-title">
      <section className={styles.card}>
        <div className={styles.code}>404</div>
        <h1 id="nf-title" className={styles.title}>
          Página não encontrada
        </h1>
        <p className={styles.desc}>A página que você procura não existe ou foi movida.</p>
        <div className={styles.actions}>
          <Button to="/" variant="primary">
            Ir para início
          </Button>
            <Button to="/websites" variant="secondary">
                Minhas páginas
            </Button>
        </div>
      </section>
    </main>
  );
}