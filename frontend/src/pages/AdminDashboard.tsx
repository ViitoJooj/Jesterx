import { useEffect, useMemo, useState } from "react";
import { get, put, del } from "../utils/api";
import { url } from "../config/Vars";
import styles from "../styles/pages/Admin.module.scss";

type Plan = {
  id: string;
  name: string;
  price_cents: number;
  description: string;
  features: string[];
  site_limit: number;
};

type AdminUser = {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  profile_img: string;
  plan: string;
  role: string;
  banned: boolean;
  created_at: string;
};

type MetricPoint = {
  label: string;
  value: number;
};

type AdminStats = {
  total_users: number;
  active_users: number;
  banned_users: number;
  paid_total_cents: number;
  paid_last_30_days_cents: number;
  new_users_series: MetricPoint[];
  payments_series: MetricPoint[];
  new_users_last_30_days: number;
  created_last_24h: number;
  payments_last_24h_cents: number;
  plans_by_usage: MetricPoint[];
  average_ticket_cents: number;
  paying_users: number;
};

function formatCurrency(cents: number) {
  if (!cents) return "Grátis";
  return (cents / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" });
}

function formatDate(iso: string) {
  if (!iso) return "-";
  const date = new Date(iso);
  return date.toLocaleDateString("pt-BR");
}

export function AdminDashboard() {
  const [plans, setPlans] = useState<Plan[]>([]);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");
  const [savingPlanId, setSavingPlanId] = useState<string | null>(null);
  const [updatingUserId, setUpdatingUserId] = useState<string | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  async function loadData() {
    try {
      setLoading(true);
      setMessage("");
      const [planRes, userRes, statsRes] = await Promise.all([get<Plan[]>("/v1/admin/plans"), get<AdminUser[]>("/v1/admin/users"), get<AdminStats>("/v1/admin/stats/overview")]);

      if (planRes.success && Array.isArray(planRes.data)) {
        setPlans(planRes.data as Plan[]);
      }
      if (userRes.success && Array.isArray(userRes.data)) {
        setUsers(userRes.data as AdminUser[]);
      }
      if (statsRes.success && statsRes.data) {
        setStats(statsRes.data as AdminStats);
      }
    } catch (error: any) {
      setMessage(error?.message || "Falha ao carregar dados de admin.");
    } finally {
      setLoading(false);
    }
  }

  function updatePlanField(id: string, field: keyof Plan, value: any) {
    setPlans((prev) => prev.map((plan) => (plan.id === id ? { ...plan, [field]: value } : plan)));
  }

  async function savePlan(plan: Plan) {
    try {
      setSavingPlanId(plan.id);
      const response = await put(`/v1/admin/plans/${plan.id}`, {
        name: plan.name,
        price_cents: Number(plan.price_cents) || 0,
        description: plan.description,
        features: plan.features,
        site_limit: Number(plan.site_limit) || 0,
      });

      const data = response.data as Plan[] | undefined;
      if (response.success && Array.isArray(data) && data.length > 0) {
        setPlans((prev) => prev.map((p) => (p.id === plan.id ? data[0] : p)));
        setMessage("Plano atualizado.");
      } else if (!response.success) {
        setMessage(response.message || "Não foi possível atualizar o plano.");
      }
    } catch (error: any) {
      setMessage(error?.message || "Erro ao salvar o plano.");
    } finally {
      setSavingPlanId(null);
    }
  }

  async function updateUser(user: AdminUser, changes: Partial<AdminUser>) {
    try {
      setUpdatingUserId(user.id);
      const response = await put(`/v1/admin/users/${user.id}`, {
        first_name: changes.first_name ?? user.first_name,
        last_name: changes.last_name ?? user.last_name,
        profile_img: changes.profile_img ?? user.profile_img,
        plan: changes.plan ?? user.plan,
        role: changes.role ?? user.role,
      });

      if (response.success && response.data) {
        setUsers((prev) => prev.map((u) => (u.id === user.id ? (response.data as AdminUser) : u)));
      } else if (!response.success) {
        setMessage(response.message || "Não foi possível atualizar o usuário.");
      }
    } catch (error: any) {
      setMessage(error?.message || "Erro ao atualizar usuário.");
    } finally {
      setUpdatingUserId(null);
    }
  }

  async function toggleBan(user: AdminUser, banned: boolean) {
    try {
      setUpdatingUserId(user.id);
      const response = await put(`/v1/admin/users/${user.id}/ban`, { banned });
      if (response.success && response.data) {
        setUsers((prev) => prev.map((u) => (u.id === user.id ? (response.data as AdminUser) : u)));
      }
    } catch (error: any) {
      setMessage(error?.message || "Erro ao atualizar status do usuário.");
    } finally {
      setUpdatingUserId(null);
    }
  }

  async function deleteUser(user: AdminUser) {
    if (!window.confirm(`Remover ${user.email}?`)) return;
    try {
      setUpdatingUserId(user.id);
      const response = await del(`/v1/admin/users/${user.id}`);
      if (response.success) {
        setUsers((prev) => prev.filter((u) => u.id !== user.id));
      } else {
        setMessage(response.message || "Não foi possível remover usuário.");
      }
    } catch (error: any) {
      setMessage(error?.message || "Erro ao remover usuário.");
    } finally {
      setUpdatingUserId(null);
    }
  }

  async function exportUsers() {
    try {
      setMessage("Gerando exportação...");
      const res = await fetch(`${url}/v1/admin/users/export`, {
        method: "GET",
        credentials: "include",
      });
      if (!res.ok) throw new Error("Falha ao exportar usuários");
      const blob = await res.blob();
      const link = document.createElement("a");
      link.href = window.URL.createObjectURL(blob);
      link.download = "usuarios.xlsx";
      link.click();
      window.URL.revokeObjectURL(link.href);
      setMessage("Exportação concluída.");
    } catch (error: any) {
      setMessage(error?.message || "Erro ao exportar usuários.");
    }
  }

  const planOptions = useMemo(() => plans.map((p) => ({ id: p.id, name: p.name })), [plans]);

  return (
    <main className={styles.container}>
      <header className={styles.header}>
        <div>
          <p className={styles.kicker}>Área administrativa</p>
          <h1>Console do administrador</h1>
          <p className={styles.subtitle}>Gerencie planos, usuários e acompanhe métricas críticas do produto.</p>
        </div>
        <div className={styles.actions}>
          <button className={styles.ghost} onClick={exportUsers}>
            Exportar usuários (XLSX)
          </button>
          <button className={styles.primary} onClick={loadData} disabled={loading}>
            Atualizar dados
          </button>
        </div>
      </header>

      {message && <div className={styles.banner}>{message}</div>}

      <section className={styles.grid}>
        <div className={styles.card}>
          <p className={styles.label}>Usuários totais</p>
          <strong className={styles.value}>{stats?.total_users ?? "-"}</strong>
          <span className={styles.meta}>+{stats?.new_users_last_30_days ?? 0} nos últimos 30 dias</span>
        </div>
        <div className={styles.card}>
          <p className={styles.label}>Usuários ativos</p>
          <strong className={styles.value}>{stats?.active_users ?? "-"}</strong>
          <span className={styles.meta}>{stats?.banned_users ?? 0} banidos</span>
        </div>
        <div className={styles.card}>
          <p className={styles.label}>Receita total</p>
          <strong className={styles.value}>{formatCurrency(stats?.paid_total_cents || 0)}</strong>
          <span className={styles.meta}>{formatCurrency(stats?.payments_last_24h_cents || 0)} nas últimas 24h</span>
        </div>
        <div className={styles.card}>
          <p className={styles.label}>Ticket médio</p>
          <strong className={styles.value}>{formatCurrency(stats?.average_ticket_cents || 0)}</strong>
          <span className={styles.meta}>{stats?.paying_users ?? 0} clientes pagantes</span>
        </div>
      </section>

      <section className={styles.section}>
        <div className={styles.sectionHeader}>
          <div>
            <p className={styles.kicker}>Planos</p>
            <h2>Editar valores e limites</h2>
          </div>
        </div>
        <div className={styles.planGrid}>
          {plans.map((plan) => (
            <div key={plan.id} className={styles.planCard}>
              <div className={styles.planHeader}>
                <input className={styles.input} value={plan.name} onChange={(e) => updatePlanField(plan.id, "name", e.target.value)} />
                <span className={styles.planId}>{plan.id}</span>
              </div>
              <label className={styles.field}>
                <span>Preço (centavos)</span>
                <input type="number" className={styles.input} value={plan.price_cents} onChange={(e) => updatePlanField(plan.id, "price_cents", Number(e.target.value))} />
              </label>
              <label className={styles.field}>
                <span>Descrição</span>
                <textarea className={styles.input} value={plan.description} onChange={(e) => updatePlanField(plan.id, "description", e.target.value)} />
              </label>
              <label className={styles.field}>
                <span>Funcionalidades (uma por linha)</span>
                <textarea
                  className={styles.input}
                  value={plan.features.join("\n")}
                  onChange={(e) => updatePlanField(plan.id, "features", e.target.value.split("\n").map((f) => f.trim()).filter(Boolean))}
                />
              </label>
              <label className={styles.field}>
                <span>Limite de sites</span>
                <input type="number" className={styles.input} value={plan.site_limit} onChange={(e) => updatePlanField(plan.id, "site_limit", Number(e.target.value))} />
              </label>
              <button className={styles.primary} onClick={() => savePlan(plan)} disabled={savingPlanId === plan.id}>
                {savingPlanId === plan.id ? "Salvando..." : "Salvar plano"}
              </button>
            </div>
          ))}
        </div>
      </section>

      <section className={styles.section}>
        <div className={styles.sectionHeader}>
          <div>
            <p className={styles.kicker}>Usuários</p>
            <h2>Gestão de contas</h2>
          </div>
        </div>

        <div className={styles.tableWrapper}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Email</th>
                <th>Plano</th>
                <th>Função</th>
                <th>Criado em</th>
                <th>Status</th>
                <th>Ações</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id}>
                  <td>
                    <div className={styles.userCell}>
                      {user.profile_img ? <img src={user.profile_img} alt={user.email} /> : <span className={styles.avatar}>{user.first_name?.[0] || "U"}</span>}
                      <div>
                        <strong>{user.email}</strong>
                        <small>
                          {user.first_name} {user.last_name}
                        </small>
                      </div>
                    </div>
                  </td>
                  <td>
                    <select className={styles.input} value={user.plan} onChange={(e) => updateUser(user, { plan: e.target.value })} disabled={updatingUserId === user.id}>
                      {planOptions.map((plan) => (
                        <option key={plan.id} value={plan.id}>
                          {plan.name}
                        </option>
                      ))}
                    </select>
                  </td>
                  <td>
                    <select className={styles.input} value={user.role} onChange={(e) => updateUser(user, { role: e.target.value })} disabled={updatingUserId === user.id}>
                      <option value="platform_user">Usuário</option>
                      <option value="platform_admin">Admin</option>
                      <option value="customer">Cliente</option>
                      <option value="admin">Admin (tenant)</option>
                      <option value="owner">Owner (tenant)</option>
                    </select>
                  </td>
                  <td>{formatDate(user.created_at)}</td>
                  <td>
                    <span className={`${styles.badge} ${user.banned ? styles.danger : styles.success}`}>{user.banned ? "Banido" : "Ativo"}</span>
                  </td>
                  <td className={styles.actionsCell}>
                    <button className={styles.ghost} onClick={() => toggleBan(user, !user.banned)} disabled={updatingUserId === user.id}>
                      {user.banned ? "Reativar" : "Banir"}
                    </button>
                    <button className={styles.danger} onClick={() => deleteUser(user)} disabled={updatingUserId === user.id}>
                      Deletar
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {users.length === 0 && !loading && <p className={styles.empty}>Nenhum usuário encontrado.</p>}
        </div>
      </section>
    </main>
  );
}
