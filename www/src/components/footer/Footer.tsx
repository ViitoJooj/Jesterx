import { Link } from "react-router-dom";
import styles from "./Footer.module.scss";

export const Footer: React.FC = () => {
  const year = new Date().getFullYear();

  return (
    <footer className={styles.footer}>
      <div className={styles.wrap}>
        <div className={styles.top}>
          <div className={styles.brand}>
            <div className={styles.logoRow}>
              <div className={styles.logo}>J</div>
              <div className={styles.name}>Jester</div>
            </div>
            <p className={styles.tagline}>SaaS para criar e gerenciar e-commerces com simplicidade.</p>
          </div>

          <div className={styles.links}>
            <div className={styles.col}>
              <h4>Produto</h4>
              <Link to="/">Início</Link>
              <Link to="/pages">Minhas páginas</Link>
              <Link to="/products">Meus produtos</Link>
              <Link to="/plans">Planos</Link>
            </div>
            <div className={styles.col}>
              <h4>Empresa</h4>
              <Link to="/about">Sobre</Link>
              <Link to="/docs">Documentação</Link>
              <Link to="/status">Status</Link>
              <a href="mailto:suporte@jesterx.com">Suporte</a>
            </div>
          </div>
        </div>

        <hr className={styles.hr} />

        <div className={styles.bottom}>
          <div className={styles.copy}>© {year} Admin. Todos os direitos reservados.</div>
          <div className={styles.mini}>
            <Link to="/privacy">Privacidade</Link>
            <Link to="/terms">Termos</Link>
            <a href="mailto:hello@example.com">Contato</a>
            <a href="https://github.com/ViitoJooj" target="_blank" rel="noreferrer">
              Por @viitoJooj
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}