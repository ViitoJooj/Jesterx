import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuthContext } from "../../hooks/AuthContext";
import { apiFetch } from "../../hooks/api";
import styles from "./Plans.module.scss"

export const Plans = () => {
  const navigate = useNavigate();
  const { isAuthenticated, websiteId } = useAuthContext();
  const [plans, setPlans] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [checkoutLoadingPlanId, setCheckoutLoadingPlanId] = useState(null);
  const [selectedPlan, setSelectedPlan] = useState(null);

  function formatPrice(price) {
    return price.toLocaleString("pt-BR", {
      style: "currency",
      currency: "BRL",
    });
  }

  const popularPlanId = useMemo(() => {
    if (!plans.length) return null;
    const sortedByPrice = [...plans].sort((a, b) => a.price - b.price);
    if (sortedByPrice.length < 2) return sortedByPrice[0].id;
    return sortedByPrice[1].id;
  }, [plans]);

  useEffect(() => {
    let cancelled = false;

    async function loadPlans() {
      try {
        setLoading(true);
        setError(null);
        const response = await apiFetch("/api/v1/plans", {
          method: "GET",
          websiteId,
          includeJsonContentType: false,
        });

        if (cancelled) return;

        const mapped = response.data.map((item) => ({
          id: item.id,
          name: item.name,
          price: item.price,
          billing_cycle: item.billing_cycle || "monthly",
          description: item.description,
          features: extractFeatures(item.description_md, item.description),
        }));

        setPlans(mapped);
      } catch (err) {
        if (cancelled) return;
        setError(
          err instanceof Error
            ? err.message
            : "Nao foi possivel carregar os planos no momento."
        );
      } finally {
        if (!cancelled) setLoading(false);
      }
    }

    loadPlans();
    return () => {
      cancelled = true;
    };
  }, [websiteId]);

  async function handleSelectPlan(planId) {
    setSelectedPlan(planId);

    if (isAuthenticated) {
      try {
        setCheckoutLoadingPlanId(planId);
        setError(null);

        const response = await apiFetch("/api/v1/payments/checkout", {
          method: "POST",
          websiteId,
          body: JSON.stringify({
            plan_id: planId,
            quantity: 1,
          }),
        });

        if (!response?.data?.checkout_url) {
          throw new Error("Checkout indisponivel no momento.");
        }

        window.location.assign(response.data.checkout_url);
      } catch (err) {
        setError(
          err instanceof Error
            ? err.message
            : "Nao foi possivel iniciar o pagamento."
        );
      } finally {
        setCheckoutLoadingPlanId(null);
      }
      return;
    }

    navigate(`/login?next=/plans&plan=${planId}`);
  }

  return (
    <main className={styles.main}>
      <div className={styles.hero}>
        <p className={styles.eyebrow}>Planos</p>
        <h1>Escolha seu plano</h1>
        <p>Comece a criar suas páginas e e-commerces hoje mesmo</p>
      </div>

      <div className={styles.container}>
        {loading && <p className={styles.infoText}>Carregando planos...</p>}
        {error && <p className={styles.errorText}>{error}</p>}

        {!loading && !error && (
          <div className={styles.plansGrid}>
          {plans.map((plan) => (
            <div
              key={plan.id}
              className={`${styles.planCard} ${
                plan.id === popularPlanId ? styles.popular : ""
              }`}
            >
              {plan.id === popularPlanId && (
                <span className={styles.badge}>Mais Popular</span>
              )}

              <h2 className={styles.planName}>{plan.name}</h2>
              {plan.description && (
                <p className={styles.planDescription}>{plan.description}</p>
              )}

              <div className={styles.priceSection}>
                <span className={styles.price}>
                  {formatPrice(plan.price)}
                </span>
                <span className={styles.period}>
                  /{plan.billing_cycle === "monthly" ? "mês" : plan.billing_cycle}
                </span>
              </div>

              <ul className={styles.features}>
                {plan.features.map((feature, index) => (
                  <li key={index}>
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                    >
                      <path
                        d="M16.6667 5L7.50004 14.1667L3.33337 10"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      />
                    </svg>
                    {feature}
                  </li>
                ))}
              </ul>

              <button
                className={`${styles.selectButton} ${
                  plan.id === popularPlanId ? styles.primary : styles.secondary
                }`}
                onClick={() => handleSelectPlan(plan.id)}
                disabled={checkoutLoadingPlanId !== null}
              >
                {checkoutLoadingPlanId === plan.id
                  ? "Redirecionando para pagamento..."
                  : selectedPlan === plan.id
                    ? "Plano selecionado"
                    : isAuthenticated
                      ? "Assinar agora"
                      : "Fazer login para assinar"}
              </button>
            </div>
          ))}
          </div>
        )}

        <div className={styles.footer}>
          <p>Pagamento processado com segurança pela Stripe.</p>
        </div>
      </div>
    </main>
  );
}

function extractFeatures(descriptionMd, description) {
  const source = descriptionMd?.trim() || "";
  if (source.length > 0) {
    const byLines = source
      .split(/\r?\n/)
      .map((line) => line.trim().replace(/^[-*•]\s*/, ""))
      .filter((line) => line.length > 0);

    if (byLines.length > 1) return byLines;

    const bySeparator = source
      .split(/[;|]/)
      .map((item) => item.trim())
      .filter((item) => item.length > 0);

    if (bySeparator.length > 1) return bySeparator;
  }

  if (description?.trim()) return [description.trim()];
  return ["Recursos sob demanda para o seu projeto."];
}
