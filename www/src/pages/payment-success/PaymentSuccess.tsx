import { useNavigate } from "react-router-dom";
import styles from "./PaymentSuccess.module.scss";

export const PaymentSuccess: React.FC = () => {
  const navigate = useNavigate();
  const confirmed = true;
  const plan = "Pro";

  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <div className={`${styles.iconCircle} ${styles.success}`}>
          <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <path
              d="M40 12L18 34L8 24"
              stroke="currentColor"
              strokeWidth="4"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        </div>

        <h1>Pagamento confirmado! 🎉</h1>

        {confirmed && plan && (
          <p className={styles.description}>
            Seu plano <strong>{plan}</strong> foi ativado com sucesso!
          </p>
        )}

        <p className={styles.description}>
          Agora você pode criar suas páginas e começar a vender.
        </p>

        <div className={styles.actions}>
          <button
            className={styles.primaryButton}
            onClick={() => navigate("/pages")}
          >
            Criar minha primeira página
          </button>

          <button
            className={styles.secondaryButton}
            onClick={() => navigate("/dashboard")}
          >
            Ir para o dashboard
          </button>
        </div>
      </div>
    </main>
  );
}
