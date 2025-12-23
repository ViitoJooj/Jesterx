import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import styles from "../styles/pages/PaymentResult.module.scss";
import { post, get } from "../utils/api";

export function PaymentSuccess() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [loading, setLoading] = useState(true);
  const [confirmed, setConfirmed] = useState(false);
  const [plan, setPlan] = useState<string>("");

  useEffect(() => {
    const sessionId = searchParams.get("session_id");

    if (!sessionId) {
      navigate("/pricing");
      return;
    }

    async function confirmPayment() {
      try {
        const response = await post("/v1/billing/confirm", {
          session_id: sessionId,
        });

        if (response.success) {
          setConfirmed(true);
          setPlan(response.data?.plan || "");

          await get("/v1/auth/refresh");
        }
      } catch (error) {
        console.error("Erro ao confirmar pagamento:", error);
      } finally {
        setLoading(false);
      }
    }

    confirmPayment();
  }, [searchParams, navigate]);

  if (loading) {
    return (
      <main className={styles.main}>
        <div className={styles.container}>
          <div className={styles.loader}></div>
          <h2>Confirmando seu pagamento... </h2>
        </div>
      </main>
    );
  }

  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <div className={`${styles.iconCircle} ${styles.success}`}>
          <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <path d="M40 12L18 34L8 24" stroke="currentColor" strokeWidth="4" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
        </div>

        <h1>Pagamento confirmado! 🎉</h1>

        {confirmed && plan && (
          <p className={styles.description}>
            Seu plano <strong>{plan}</strong> foi ativado com sucesso!
          </p>
        )}

        <p className={styles.description}>Agora você pode criar suas páginas e começar a vender.</p>

        <div className={styles.actions}>
          <button className={styles.primaryButton} onClick={() => navigate("/pages")}>
            Criar minha primeira página
          </button>
          <button className={styles.secondaryButton} onClick={() => navigate("/dashboard")}>
            Ir para o dashboard
          </button>
        </div>
      </div>
    </main>
  );
}
