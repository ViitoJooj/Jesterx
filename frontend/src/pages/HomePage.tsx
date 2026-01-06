import { Link } from "react-router-dom";
import styles from "../styles/pages/HomePage.module.scss";
import buttonStyles from "../styles/components/Button.module.scss";
import RotatingWord from "../components/RotatingWord";

export function HomePage() {
  return (
    <>
      <main className={styles.main}>
        <div className={styles.header}>
          <h1>
            Construa seu projeto
            <br />
            <RotatingWord items={["mais rápido", "com clareza", "sem código", "do seu jeito"]} />
          </h1>

          <h2>
            Jester é a plataforma low-code para criar desde e-commerces completos até landing pages e experiências digitais em um só lugar. Conecte ERPs, gerencie produtos físicos e digitais e lance sua operação sem escrever código. Dê autonomia para sua equipe crescer e escalar.
          </h2>

          <div className={styles.cta}>
            <Link to="/pricing" className={`${buttonStyles.default_button} ${buttonStyles["default_button--primary"]}`}>
              Começar agora
            </Link>
            <Link to="/register" className={`${buttonStyles.default_button} ${buttonStyles["default_button--secondary"]}`}>
              Ver demo
            </Link>
          </div>
        </div>
      </main>

      {/* Second part */}
      <section></section>
    </>
  );
}
