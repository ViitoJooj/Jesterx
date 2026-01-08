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
  route_limit: number;
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
  return new Date(iso).toLocaleDateString("pt-BR");
}

function StatCard({ label, value, meta }: { label: string; value: any; meta?: any }) {
  return (
    <div className={styles.card}>
      <p className={styles.label}>{label}</p>
      <strong className={styles.value}>{value}</strong>
      {meta && <span className={styles.meta}>{meta}</span>}
    </div>
  );
}

export function AdminDashboard() {
  const [plans, setPlans] = useState<Plan[]>([]);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");
  const [savingPlanId, setSavingPlanId] = useState<string | null>(null);
  const [updatingUserId, setUpdatingUserId] = useState<string | null>(null);
  const [openPlanId, setOpenPlanId] = useState<string | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  async function loadData() {
    try {
      setLoading(true);
      const [p, u, s] = await Promise.all([
        get<Plan[]>("/v1/admin/plans"),
        get<AdminUser[]>("/v1/admin/users"),
        get<AdminStats>("/v1/admin/stats/overview"),
      ]);
      if (p.success) setPlans(p.data || []);
      if (u.success) setUsers(u.data || []);
      if (s.success) setStats(s.data || null);
    } catch (e: any) {
      setMessage(e.message);
    } finally {
      setLoading(false);
    }
  }

  function updatePlan(id: string, field: keyof Plan, value: any) {
    setPlans((prev) => prev.map((p) => (p.id === id ? { ...p, [field]: value } : p)));
  }

  async function savePlan(plan: Plan) {
    try {
      setSavingPlanId(plan.id);
      const res = await put(`/v1/admin/plans/${plan.id}`, plan);
      if (res.success && res.data) {
        setPlans((prev) => prev.map((p) => (p.id === plan.id ? res.data[0] : p)));
        setOpenPlanId(null);
      }
    } finally {
      setSavingPlanId(null);
    }
  }

  async function updateUser(user: AdminUser, changes: Partial<AdminUser>) {
    try {
      setUpdatingUserId(user.id);
      const res = await put(`/v1/admin/users/${user.id}`, { ...user, ...changes });
      if (res.success) {
        setUsers((prev) => prev.map((u) => (u.id === user.id ? res.data : u)));
      }
    } finally {
      setUpdatingUserId(null);
    }
  }

  async function toggleBan(user: AdminUser) {
    try {
      setUpdatingUserId(user.id);
      const res = await put(`/v1/admin/users/${user.id}/ban`, { banned: !user.banned });
      if (res.success) {
        setUsers((prev) => prev.map((u) => (u.id === user.id ? res.data : u)));
      }
    } finally {
      setUpdatingUserId(null);
    }
  }

  async function deleteUser(user: AdminUser) {
    if (!window.confirm(`Remover ${user.email}?`)) return;

    try {
      setUpdatingUserId(user.id);
      const res = await del(`/v1/admin/users/${user.id}`);
      if (res.success) setUsers((prev) => prev.filter((u) => u.id !== user.id));
    } finally {
      setUpdatingUserId(null);
    }
  }

  async function exportUsers() {
    const res = await fetch(`${url}/v1/admin/users/export`, { credentials: "include" });
    const blob = await res.blob();
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = "usuarios.xlsx";
    link.click();
    URL.revokeObjectURL(link.href);
  }

  const planOptions = useMemo(() => plans.map((p) => ({ id: p.id, name: p.name })), [plans]);

  return (
    <main className={styles.container}>
      <header className={styles.header}>
        <div>
          <p className={styles.kicker}>Área administrativa</p>
          <h1>Console do administrador</h1>
          <p className={styles.subtitle}>Gerencie planos, usuários e métricas</p>
        </div>
        <div className={styles.actions}>
          <button className={styles.ghost} onClick={exportUsers}>Exportar usuários</button>
          <button className={styles.primary} onClick={loadData} disabled={loading}>Atualizar</button>
        </div>
      </header>

      {message && <div className={styles.banner}>{message}</div>}

      <section className={styles.grid}>
        <StatCard label="Usuários totais" value={stats?.total_users ?? "-"} meta={`+${stats?.new_users_last_30_days ?? 0} em 30 dias`} />
        <StatCard label="Usuários ativos" value={stats?.active_users ?? "-"} meta={`${stats?.banned_users ?? 0} banidos`} />
        <StatCard label="Receita total" value={formatCurrency(stats?.paid_total_cents || 0)} meta={`${formatCurrency(stats?.payments_last_24h_cents || 0)} em 24h`} />
        <StatCard label="Ticket médio" value={formatCurrency(stats?.average_ticket_cents || 0)} meta={`${stats?.paying_users ?? 0} pagantes`} />
      </section>

      <section className={styles.section}>
        <h2>Planos</h2>
        <div className={styles.planGrid}>
          {plans.map((plan) => (
            <div key={plan.id} className={styles.planCard}>
              <div className={styles.planHeader}>
                <strong>{plan.name}</strong>
                <button className={styles.ghost} onClick={() => setOpenPlanId(plan.id === openPlanId ? null : plan.id)}>
                  {openPlanId === plan.id ? "Fechar" : "Editar"}
                </button>
              </div>

              {openPlanId === plan.id && (
                <>
                  <input className={styles.input} value={plan.name} onChange={(e) => updatePlan(plan.id, "name", e.target.value)} />
                  <input type="number" className={styles.input} value={plan.price_cents} onChange={(e) => updatePlan(plan.id, "price_cents", Number(e.target.value))} />
                  <textarea className={styles.input} value={plan.description} onChange={(e) => updatePlan(plan.id, "description", e.target.value)} />
                  <textarea
                    className={styles.input}
                    value={plan.features.join("\n")}
                    onChange={(e) => updatePlan(plan.id, "features", e.target.value.split("\n").filter(Boolean))}
                  />
                  <input type="number" className={styles.input} value={plan.site_limit} onChange={(e) => updatePlan(plan.id, "site_limit", Number(e.target.value))} />
                  <input type="number" className={styles.input} value={plan.route_limit} onChange={(e) => updatePlan(plan.id, "route_limit", Number(e.target.value))} />
                  <button className={styles.primary} onClick={() => savePlan(plan)} disabled={savingPlanId === plan.id}>
                    {savingPlanId === plan.id ? "Salvando..." : "Salvar"}
                  </button>
                </>
              )}
            </div>
          ))}
        </div>
      </section>

      <section className={styles.section}>
        <h2>Usuários</h2>
        <div className={styles.tableWrapper}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Email</th>
                <th>Plano</th>
                <th>Função</th>
                <th>Criado</th>
                <th>Status</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.id}>
                  <td>{u.email}</td>
                  <td>
                    <select className={styles.input} value={u.plan} onChange={(e) => updateUser(u, { plan: e.target.value })}>
                      {planOptions.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
                    </select>
                  </td>
                  <td>
                    <select className={styles.input} value={u.role} onChange={(e) => updateUser(u, { role: e.target.value })}>
                      <option value="platform_user">Usuário</option>
                      <option value="platform_admin">Admin</option>
                      <option value="customer">Cliente</option>
                      <option value="admin">Admin</option>
                      <option value="owner">Owner</option>
                    </select>
                  </td>
                  <td>{formatDate(u.created_at)}</td>
                  <td>
                    <span className={`${styles.badge} ${u.banned ? styles.danger : styles.success}`}>
                      {u.banned ? "Banido" : "Ativo"}
                    </span>
                  </td>
                  <td className={styles.actionsCell}>
                    <button className={styles.ghost} onClick={() => toggleBan(u)} disabled={updatingUserId === u.id}>
                      {u.banned ? "Reativar" : "Banir"}
                    </button>
                    <button className={styles.danger} onClick={() => deleteUser(u)} disabled={updatingUserId === u.id}>
                      Deletar
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </main>
  );
}
