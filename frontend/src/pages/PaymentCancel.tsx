import { useNavigate } from "react-router-dom";
import styles from "../styles/pages/PaymentResult.module.scss";

export function PaymentCancel() {
  const navigate = useNavigate();

  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <div className={`${styles.iconCircle} ${styles.warning}`}>
          <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <path d="M24 16V26M24 32H24.02M44 24C44 35.0457 35.0457 44 24 44C12.9543 44 4 35.0457 4 24C4 12.9543 12.9543 4 24 4C35.0457 4 44 12.9543 44 24Z" stroke="currentColor" strokeWidth="4" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
        </div>

        <h1>Pagamento cancelado</h1>

        <p className={styles.description}>Você cancelou o processo de pagamento. Nenhuma cobrança foi realizada.</p>

        <p className={styles.description}>Caso tenha alguma dúvida sobre os planos, entre em contato com nosso suporte.</p>

        <div className={styles.actions}>
          <button className={styles.primaryButton} onClick={() => navigate("/pricing")}>
            Ver planos novamente
          </button>
          <button className={styles.secondaryButton} onClick={() => navigate("/")}>
            Voltar para o início
          </button>
        </div>
      </div>
    </main>
  );
}
