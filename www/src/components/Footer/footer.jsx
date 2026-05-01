import { Link } from "react-router-dom";
import styles from "./footer.module.scss";
import { Brand } from "../Brand/brand";

export function Footer() {
  const year = new Date().getFullYear();

  return (
    <footer className={styles.footer}>
      <div className={styles.wrap}>
        <div className={styles.top}>
          <div className={styles.brand}>
            <Brand size="lg" />
            <p className={styles.tagline}>Do conceito ao digital, com a estrutura certa desde o início.</p>
          </div>

          <div className={styles.links}>
            <div className={styles.col}>
              <h4>Plataforma</h4>
              <Link to="/">Início</Link>
              <Link to="/pages">Páginas</Link>
              <Link to="/products">Produtos</Link>
              <Link to="/plans">Planos</Link>
            </div>
            <div className={styles.col}>
              <h4>Empresa</h4>
              <Link to="/about">Sobre nós</Link>
              <Link to="/docs">Documentação</Link>
              <Link to="/status">Status</Link>
              <a href="mailto:suporte@jesterx.com">Fale conosco</a>
            </div>
          </div>
        </div>

        <hr className={styles.hr} />

        <div className={styles.bottom}>
          <div className={styles.copy}>© {year} Jesterx. Todos os direitos reservados.</div>
          <div className={styles.mini}>
            <Link to="/privacy">Privacidade</Link>
            <Link to="/terms">Termos de uso</Link>
            <a href="mailto:suporte@jesterx.com">Contato</a>
            <a href="https://github.com/ViitoJooj" target="_blank" rel="noreferrer">
              Feito por @viitoJooj
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
