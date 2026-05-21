import { Link } from "react-router-dom";
import styles from "./Home.module.scss";
import Button from "../../components/button/Button";
import RotatingWord from "../../components/rotatingWord/RotatingWord";

export const Home = () => {
  return (
 <main className={styles.main}>
      <div className={styles.header}>
        <h1>
          Seu negócio digital,
          <br />
          <RotatingWord items={["do zero.", "com precisão.", "sem código.", "do seu jeito."]} />
        </h1>

        <h2>
          E-commerces, landing pages e lojas digitais, criados do jeito certo desde o
          início. Catálogo, pagamentos e conteúdo em um único lugar. Sem escrever uma
          linha de código.
        </h2>

        <div className={styles.cta}>
          <Button to="/plans" variant="primary">
            Começar agora
          </Button>

          <Button to="/register" variant="secondary">
            Criar conta grátis
          </Button>
        </div>
      </div>
    </main>
  );
};
