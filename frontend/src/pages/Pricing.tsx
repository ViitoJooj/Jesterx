import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import styles from "../styles/pages/Pricing.module.scss";
import { get, post } from "../utils/api";

type PlanType = string;

interface Plan {
  id: PlanType;
  name: string;
  price_cents: number;
  features: string[];
  popular?: boolean;
  description?: string;
}

const fallbackPlans: Plan[] = [
  {
    id: "business",
    name: "Business",
    price_cents: 4900,
    features: ["1 site", "Páginas ilimitadas", "Suporte por email", "Templates básicos", "SSL incluído"],
  },
  {
    id: "pro",
    name: "Pro",
    price_cents: 9900,
    popular: true,
    features: ["Até 10 sites", "Páginas ilimitadas", "Suporte prioritário", "Todos os templates", "SSL incluído", "Analytics avançado", "Integrações premium"],
  },
  {
    id: "enterprise",
    name: "Enterprise",
    price_cents: 19900,
    features: ["Até 50 sites", "Páginas ilimitadas", "Suporte 24/7", "Templates customizados", "SSL incluído", "Analytics avançado", "Todas as integrações", "API dedicada", "Suporte a White Label"],
  },
];

export function Pricing() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [selectedPlan, setSelectedPlan] = useState<PlanType | null>(null);
  const [plans, setPlans] = useState<Plan[]>(fallbackPlans);

  useEffect(() => {
    (async () => {
      try {
        const response = await get<Plan[]>("/v1/plans");
        if (response.success && Array.isArray(response.data) && response.data.length > 0) {
          const normalized = (response.data as Plan[]).map((plan) => ({
            ...plan,
            price_cents: Number(plan.price_cents) || 0,
            features: Array.isArray(plan.features) ? plan.features : [],
          }));
          setPlans(normalized);
        }
      } catch {
        setPlans(fallbackPlans);
      }
    })();
  }, []);

  function formatPrice(cents: number) {
    if (!cents) return "R$ 0,00";
    return (cents / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" });
  }

  async function handleSelectPlan(planId: PlanType) {
    setSelectedPlan(planId);
    setLoading(true);

    try {
      const response = await post("/v1/billing/checkout", {
        plan: planId,
      });

      if (response.success && response.data?.checkout_url) {
        // Redirecionar para o Stripe Checkout
        window.location.href = response.data.checkout_url;
      } else {
        alert("Erro ao criar sessão de checkout");
        setLoading(false);
      }
    } catch (error: any) {
      console.error("Erro ao criar checkout:", error);
      alert(error?.message || "Erro ao processar pagamento.  Tente novamente.");
      setLoading(false);
      setSelectedPlan(null);
    }
  }

  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <div className={styles.header}>
          <h1>Escolha seu plano</h1>
          <p>Comece a criar suas páginas e e-commerces hoje mesmo</p>
        </div>

        <div className={styles.plansGrid}>
              {plans.map((plan) => (
                <div key={plan.id} className={`${styles.planCard} ${plan.popular ? styles.popular : ""}`}>
                  {plan.popular && <span className={styles.badge}>Mais Popular</span>}

                  <h2 className={styles.planName}>{plan.name}</h2>
                  {plan.description && <p className={styles.planDescription}>{plan.description}</p>}

                  <div className={styles.priceSection}>
                    <span className={styles.price}>{formatPrice(plan.price_cents)}</span>
                    <span className={styles.period}>/mês</span>
                  </div>

              <ul className={styles.features}>
                {plan.features.map((feature, index) => (
                  <li key={index}>
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
                      <path d="M16.6667 5L7.50004 14.1667L3.33337 10" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                    </svg>
                    {feature}
                  </li>
                ))}
              </ul>

              <button className={`${styles.selectButton} ${plan.popular ? styles.primary : styles.secondary}`} onClick={() => handleSelectPlan(plan.id)} disabled={loading}>
                {loading && selectedPlan === plan.id ? "Processando..." : "Começar agora"}
              </button>
            </div>
          ))}
        </div>

        <div className={styles.footer}>
          <p>Todos os planos incluem 14 dias de teste grátis. Sem cartão de crédito necessário.</p>
          <p>
            <a href="/contact">Precisa de um plano customizado?</a>
          </p>
        </div>
      </div>
    </main>
  );
}
