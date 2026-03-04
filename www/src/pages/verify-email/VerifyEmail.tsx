import { Link, useSearchParams } from "react-router-dom";
import Button from "../../components/button/Button";
import styles from "./VerifyEmail.module.scss";

export const VerifyEmail: React.FC = () => {
  const [searchParams] = useSearchParams();
  const email = searchParams.get("email");

  return (
    <main className={styles.main}>
      <section className={styles.card}>
        <div className={styles.icon}>
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
            <path
              d="M4 6H20V18H4V6ZM4 7L12 13L20 7"
              stroke="currentColor"
              strokeWidth="1.8"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        </div>

        <h1>Verifique seu email</h1>
        <p>
          Enviamos um link de confirmação para{" "}
          <strong>{email || "o email cadastrado"}</strong>. Abra sua caixa de
          entrada e confirme para ativar sua conta.
        </p>
        <p className={styles.subtle}>
          Não encontrou? Verifique também spam e promoções.
        </p>

        <div className={styles.actions}>
          <Button to="/login" variant="primary">
            Ir para login
          </Button>
          <Link to="/register" className={styles.link}>
            Cadastrar outro email
          </Link>
        </div>
      </section>
    </main>
  );
};
