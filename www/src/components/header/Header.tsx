import { useState } from "react";
import { NavLink, Link, useNavigate } from "react-router-dom";
import styles from "./Header.module.scss";
import Button from "../button/Button";
import { useAuthContext } from "../../hooks/AuthContext";

export function Header() {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();

  const { me, logout, isAuthenticated, loading } = useAuthContext();

  async function handleLogout() {
    try {
      await logout();
    } finally {
      navigate("/login", { replace: true });
    }
  }

  return (
    <>
      <header className={styles.header}>
        <div className={styles.inner}>
          <Link to="/" className={styles.brand}>
            <span className={styles.logo}>J</span>
            <span className={styles.name}>Jester</span>
          </Link>

          <nav className={styles.nav}>
            <NavLink to="/" end>Início</NavLink>
            <NavLink to="/pages">Páginas</NavLink>
            <NavLink to="/products">Produtos</NavLink>
            <NavLink to="/temas">Temas</NavLink>
            <NavLink to="/api">API</NavLink>
          </nav>

          <div className={styles.actions}>
            {/* enquanto está “bootando” (tentando refresh/me), você pode mostrar skeleton/spinner */}
            {loading && me === null && !isAuthenticated ? (
              <span>...</span>
            ) : me ? (
              <>
                <span className={styles.userLabel}>{me.email}</span>
                <Button type="button" variant="secondary" onClick={handleLogout}>
                  Sair
                </Button>
              </>
            ) : (
              <>
                <Button to="/login" variant="secondary">Entrar</Button>
                <Button to="/register" variant="primary">Criar conta</Button>
              </>
            )}
          </div>

          <button className={styles.menu_btn} onClick={() => setOpen(true)}>☰</button>
        </div>
      </header>

      <div
        className={`${styles.drawer_backdrop} ${open ? styles.open : ""}`}
        onClick={() => setOpen(false)}
      />

      <aside className={`${styles.drawer} ${open ? styles.open : ""}`}>
        <div className={styles.drawer_header}>
          <p>Menu</p>
          <button onClick={() => setOpen(false)}>×</button>
        </div>

        <div className={styles.drawer_ctas_top}>
          {me ? (
            <Button
              type="button"
              variant="secondary"
              onClick={() => {
                setOpen(false);
                handleLogout();
              }}
            >
              Logout
            </Button>
          ) : (
            <>
              <Button to="/login" variant="secondary" onClick={() => setOpen(false)}>
                Entrar
              </Button>
              <Button to="/register" variant="primary" onClick={() => setOpen(false)}>
                Criar conta
              </Button>
            </>
          )}
        </div>

        <nav className={styles.drawer_nav}>
          <NavLink to="/" onClick={() => setOpen(false)}>Início</NavLink>
          <NavLink to="/pages" onClick={() => setOpen(false)}>Páginas</NavLink>
          <NavLink to="/products" onClick={() => setOpen(false)}>Produtos</NavLink>
          <NavLink to="/temas" onClick={() => setOpen(false)}>Temas</NavLink>
          <NavLink to="/api" onClick={() => setOpen(false)}>API</NavLink>
        </nav>
      </aside>
    </>
  );
}