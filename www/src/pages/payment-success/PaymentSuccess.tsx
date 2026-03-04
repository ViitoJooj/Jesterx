import { useEffect, useMemo, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { apiFetch } from "../../hooks/api";
import { useAuthContext } from "../../hooks/AuthContext";
import styles from "./PaymentSuccess.module.scss";

type ConfirmResponse = {
  success: boolean;
  message: string;
};

export const PaymentSuccess: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { websiteId } = useAuthContext();
  const [loading, setLoading] = useState(true);
  const [confirmed, setConfirmed] = useState(false);
  const [message, setMessage] = useState("Confirmando pagamento...");

  const sessionId = useMemo(() => searchParams.get("session_id") ?? "", [searchParams]);
  const planId = useMemo(() => searchParams.get("plan_id") ?? "", [searchParams]);

  useEffect(() => {
    let cancelled = false;

    async function run() {
      if (!sessionId) {
        setLoading(false);
        setConfirmed(false);
        setMessage("Sessão de pagamento não encontrada.");
        return;
      }

      try {
        const response = await apiFetch<ConfirmResponse>(
          `/api/v1/payments/confirm?session_id=${encodeURIComponent(sessionId)}`,
          {
            method: "GET",
            websiteId,
          }
        );
        if (cancelled) return;
        setConfirmed(true);
        setMessage(response.message || "Pagamento confirmado.");
      } catch (error) {
        if (cancelled) return;
        setConfirmed(false);
        setMessage(
          error instanceof Error
            ? error.message
            : "Não foi possível confirmar seu pagamento agora."
        );
      } finally {
        if (!cancelled) setLoading(false);
      }
    }

    run();
    return () => {
      cancelled = true;
    };
  }, [sessionId, websiteId]);

  return (
    <main className={styles.main}>
      <div className={styles.container}>
        {loading ? (
          <div className={styles.loader} />
        ) : confirmed ? (
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
        ) : (
          <div className={`${styles.iconCircle} ${styles.warning}`}>
            <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
              <path
                d="M24 16V26M24 32H24.02M44 24C44 35.0457 35.0457 44 24 44C12.9543 44 4 35.0457 4 24C4 12.9543 12.9543 4 24 4C35.0457 4 44 12.9543 44 24Z"
                stroke="currentColor"
                strokeWidth="4"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </div>
        )}

        <h1>{confirmed ? "Pagamento confirmado!" : "Não foi possível confirmar agora"}</h1>

        {planId && (
          <p className={styles.description}>
            Referência do plano: <strong>{planId}</strong>
          </p>
        )}

        <p className={styles.description}>{message}</p>

        <div className={styles.actions}>
          <button
            className={styles.primaryButton}
            onClick={() => navigate("/pages")}
          >
            Ir para minhas páginas
          </button>

          <button
            className={styles.secondaryButton}
            onClick={() => navigate("/plans")}
          >
            Voltar para planos
          </button>
        </div>
      </div>
    </main>
  );
}
