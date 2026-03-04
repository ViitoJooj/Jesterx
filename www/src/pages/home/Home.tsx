import React from "react";
import styles from "./Home.module.scss";
import Button from "../../components/button/Button";
import RotatingWord from "../../components/rotatingWord/RotatingWord";

export const Home: React.FC = () => {
  return (
    <main className={styles.main}>
      <div className={styles.header}>
        <h1>
          Construa seu projeto
          <br />
          <RotatingWord items={["mais rápido", "com clareza", "sem código", "do seu jeito"]} />
        </h1>

        <h2>
          Jester é a plataforma low-code para criar desde e-commerces completos até landing pages e experiências digitais em um só lugar.
          Conecte ERPs, gerencie produtos físicos e digitais e lance sua operação sem escrever código.
          Dê autonomia para sua equipe crescer e escalar.
        </h2>

        <div className={styles.cta}>
          <Button to="/pricing" variant="primary">
            Começar agora
          </Button>

          <Button to="/register" variant="secondary">
            Ver demo
          </Button>
        </div>
      </div>
    </main>
  );
};