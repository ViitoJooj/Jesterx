import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { apiFetch } from "../../hooks/api";
import { useAuthContext } from "../../hooks/AuthContext";
import styles from "./AdminDashboard.module.scss";

type Stats = {
  total_users: number;
  total_sites: number;
  total_products: number;
  total_orders: number;
  total_revenue: number;
  platform_fees: number;
  active_plans: number;
  new_users_today: number;
  orders_today: number;
};

type AdminUser = {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  role: string;
  verified_email: boolean;
  created_at: string;
  site_count: number;
};

type AdminSite = {
  id: string;
  name: string;
  type: string;
  banned: boolean;
  created_at: string;
  owner_name: string;
  owner_email: string;
  version_count: number;
};

type AdminOrder = {
  id: string;
  website_id: string;
  site_name: string;
  buyer_name: string;
  buyer_email: string;
  status: string;
  subtotal: number;
  platform_fee: number;
  total: number;
  created_at: string;
};

type DayRevenue = {
  day: string;
  order_count: number;
  gmv: number;
  fee: number;
};

type Tab = "overview" | "users" | "sites" | "orders" | "revenue";

export const AdminDashboard: React.FC = () => {
  const { websiteId } = useAuthContext();
  const navigate = useNavigate();
  const [tab, setTab] = useState<Tab>("overview");
  const [stats, setStats] = useState<Stats | null>(null);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [sites, setSites] = useState<AdminSite[]>([]);
  const [orders, setOrders] = useState<AdminOrder[]>([]);
  const [revenue, setRevenue] = useState<DayRevenue[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    Promise.all([
      apiFetch<{ data: Stats }>("/api/v1/admin/stats", { websiteId }),
      apiFetch<{ data: AdminUser[] }>("/api/v1/admin/users", { websiteId }),
      apiFetch<{ data: AdminSite[] }>("/api/v1/admin/sites", { websiteId }),
      apiFetch<{ data: AdminOrder[] }>("/api/v1/admin/orders", { websiteId }),
      apiFetch<{ data: DayRevenue[] }>("/api/v1/admin/revenue", { websiteId }),
    ])
      .then(([s, u, si, o, r]) => {
        setStats(s.data);
        setUsers(u.data ?? []);
        setSites(si.data ?? []);
        setOrders(o.data ?? []);
        setRevenue(r.data ?? []);
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [websiteId]);

  const fmt = (n: number) =>
    new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" }).format(n);
  const fmtDate = (s: string) => new Date(s).toLocaleDateString("pt-BR");
  const truncate = (s: string, n = 8) => s.slice(0, n) + "...";

  const STATUS_COLORS: Record<string, string> = {
    pending: "#f59e0b",
    paid: "#10b981",
    shipped: "#3b82f6",
    delivered: "#6366f1",
    cancelled: "#ef4444",
    refunded: "#8b5cf6",
  };

  if (loading) return <div className={styles.loading}>Carregando painel...</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  const TABS: { id: Tab; label: string }[] = [
    { id: "overview", label: "Visão Geral" },
    { id: "users", label: `Usuários (${users.length})` },
    { id: "sites", label: `Sites (${sites.length})` },
    { id: "orders", label: `Pedidos (${orders.length})` },
    { id: "revenue", label: "Receita" },
  ];

  return (
    <div className={styles.root}>
      <header className={styles.header}>
        <button className={styles.backBtn} onClick={() => navigate("/status")}>
          ← Voltar
        </button>
        <h1 className={styles.title}>⚙️ Painel Administrativo</h1>
      </header>

      <div className={styles.tabs}>
        {TABS.map((t) => (
          <button
            key={t.id}
            className={`${styles.tab} ${tab === t.id ? styles.tabActive : ""}`}
            onClick={() => setTab(t.id)}
          >
            {t.label}
          </button>
        ))}
      </div>

      <div className={styles.content}>
        {tab === "overview" && stats && (
          <div className={styles.statsGrid}>
            {[
              { label: "Usuários", value: stats.total_users, color: "#3b82f6" },
              { label: "Sites", value: stats.total_sites, color: "#6366f1" },
              { label: "Produtos", value: stats.total_products, color: "#8b5cf6" },
              { label: "Pedidos", value: stats.total_orders, color: "#f59e0b" },
              { label: "Receita Total", value: fmt(stats.total_revenue), color: "#10b981" },
              { label: "Taxa Plataforma", value: fmt(stats.platform_fees), color: "#ff5d1f" },
              { label: "Planos Ativos", value: stats.active_plans, color: "#06b6d4" },
              { label: "Novos Hoje", value: stats.new_users_today, color: "#84cc16" },
              { label: "Pedidos Hoje", value: stats.orders_today, color: "#f97316" },
            ].map((card) => (
              <div
                key={card.label}
                className={styles.statCard}
                style={{ borderTopColor: card.color }}
              >
                <div className={styles.statValue} style={{ color: card.color }}>
                  {card.value}
                </div>
                <div className={styles.statLabel}>{card.label}</div>
              </div>
            ))}
          </div>
        )}

        {tab === "users" && (
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Nome</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Verificado</th>
                  <th>Sites</th>
                  <th>Criado em</th>
                </tr>
              </thead>
              <tbody>
                {users.map((u) => (
                  <tr key={u.id}>
                    <td>
                      <code>{truncate(u.id)}</code>
                    </td>
                    <td>
                      {u.first_name} {u.last_name}
                    </td>
                    <td>{u.email}</td>
                    <td>
                      <span
                        className={`${styles.badge} ${
                          u.role === "admin" ? styles.badgeAdmin : styles.badgeUser
                        }`}
                      >
                        {u.role}
                      </span>
                    </td>
                    <td>{u.verified_email ? "✅" : "❌"}</td>
                    <td>{u.site_count}</td>
                    <td>{fmtDate(u.created_at)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {tab === "sites" && (
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Nome</th>
                  <th>Tipo</th>
                  <th>Dono</th>
                  <th>Versões</th>
                  <th>Banido</th>
                  <th>Criado em</th>
                </tr>
              </thead>
              <tbody>
                {sites.map((s) => (
                  <tr key={s.id}>
                    <td>
                      <strong>{s.name}</strong>
                    </td>
                    <td>
                      <span className={styles.badge}>{s.type}</span>
                    </td>
                    <td>
                      <div>{s.owner_name}</div>
                      <div className={styles.subtext}>{s.owner_email}</div>
                    </td>
                    <td>{s.version_count}</td>
                    <td>{s.banned ? "🚫" : "✅"}</td>
                    <td>{fmtDate(s.created_at)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {tab === "orders" && (
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Site</th>
                  <th>Comprador</th>
                  <th>Status</th>
                  <th>Total</th>
                  <th>Taxa Plat.</th>
                  <th>Data</th>
                </tr>
              </thead>
              <tbody>
                {orders.map((o) => (
                  <tr key={o.id}>
                    <td>
                      <code>{truncate(o.id)}</code>
                    </td>
                    <td>{o.site_name}</td>
                    <td>
                      <div>{o.buyer_name}</div>
                      <div className={styles.subtext}>{o.buyer_email}</div>
                    </td>
                    <td>
                      <span
                        className={styles.badge}
                        style={{
                          background: (STATUS_COLORS[o.status] ?? "#6b7a99") + "22",
                          color: STATUS_COLORS[o.status] ?? "#6b7a99",
                          borderColor: STATUS_COLORS[o.status] ?? "#6b7a99",
                        }}
                      >
                        {o.status}
                      </span>
                    </td>
                    <td className={styles.money}>{fmt(o.total)}</td>
                    <td className={styles.fee}>{fmt(o.platform_fee)}</td>
                    <td>{fmtDate(o.created_at)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {tab === "revenue" && (
          <div className={styles.tableWrap}>
            <div className={styles.revenueHeader}>
              <h3>Receita dos últimos 30 dias</h3>
              <div className={styles.revenueSummary}>
                GMV Total:{" "}
                <strong>{fmt(revenue.reduce((s, d) => s + d.gmv, 0))}</strong>
                &nbsp;·&nbsp; Taxas:{" "}
                <strong className={styles.fee}>
                  {fmt(revenue.reduce((s, d) => s + d.fee, 0))}
                </strong>
              </div>
            </div>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Data</th>
                  <th>Pedidos</th>
                  <th>GMV</th>
                  <th>Taxa Plataforma</th>
                </tr>
              </thead>
              <tbody>
                {revenue.map((d) => (
                  <tr key={d.day}>
                    <td>{d.day}</td>
                    <td>{d.order_count}</td>
                    <td className={styles.money}>{fmt(d.gmv)}</td>
                    <td className={styles.fee}>{fmt(d.fee)}</td>
                  </tr>
                ))}
                {revenue.length === 0 && (
                  <tr>
                    <td
                      colSpan={4}
                      style={{ textAlign: "center", color: "#9aa5bc", padding: "20px" }}
                    >
                      Nenhum pedido nos últimos 30 dias
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};
