import { useEffect, useRef, useState } from "react";
import { NavLink, Link, useNavigate } from "react-router-dom";
import styles from "./Header.module.scss";
import Button from "../button/Button";
import { useAuthContext } from "../../hooks/AuthContext";

export function Header() {
  const [open, setOpen] = useState(false);
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const userMenuRef = useRef<HTMLDivElement | null>(null);
  const navigate = useNavigate();

  const { me, logout, isAuthenticated, loading } = useAuthContext();
  const firstName = me?.first_name?.trim() || "Perfil";
  const avatarUrl = me?.avatar_url;
  const avatarFallback = firstName.slice(0, 1).toUpperCase();

  useEffect(() => {
    function onPointerDown(event: MouseEvent) {
      if (!userMenuRef.current) return;
      if (!userMenuRef.current.contains(event.target as Node)) {
        setUserMenuOpen(false);
      }
    }

    document.addEventListener("mousedown", onPointerDown);
    return () => document.removeEventListener("mousedown", onPointerDown);
  }, []);

  async function handleLogout() {
    try {
      await logout();
    } finally {
      setUserMenuOpen(false);
      setOpen(false);
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
            <NavLink to="/pages">Minhas páginas</NavLink>
            <NavLink to="/plans">Planos</NavLink>
          </nav>

          <div className={styles.actions}>
            {/* enquanto está “bootando” (tentando refresh/me), você pode mostrar skeleton/spinner */}
            {loading && me === null && !isAuthenticated ? (
              <span>...</span>
            ) : me ? (
              <div className={styles.user_menu} ref={userMenuRef}>
                <button
                  type="button"
                  className={styles.user_trigger}
                  onClick={() => setUserMenuOpen((prev) => !prev)}
                >
                  <span className={styles.avatar_btn}>
                    {avatarUrl ? (
                      <img src={avatarUrl} alt={firstName} className={styles.avatar_img} />
                    ) : (
                      <span className={styles.avatar_fallback}>{avatarFallback}</span>
                    )}
                  </span>
                  <span className={styles.avatar_name}>{firstName}</span>
                </button>

                <div
                  className={`${styles.user_menu_dropdown} ${userMenuOpen ? styles.open : ""}`}
                >
                  <Link
                    to="/profile"
                    className={styles.user_menu_item}
                    onClick={() => setUserMenuOpen(false)}
                  >
                    Perfil
                  </Link>
                  <Link
                    to="/settings"
                    className={styles.user_menu_item}
                    onClick={() => setUserMenuOpen(false)}
                  >
                    Configurações
                  </Link>
                  <div className={styles.menu_divider} />
                  <button
                    type="button"
                    className={`${styles.user_menu_item_btn} ${styles.danger}`}
                    onClick={handleLogout}
                  >
                    Logout
                  </button>
                </div>
              </div>
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
            <>
              <p className={styles.drawer_user_label}>{firstName}</p>
              <Button to="/profile" variant="secondary" onClick={() => setOpen(false)}>
                Perfil
              </Button>
              <Button to="/settings" variant="secondary" onClick={() => setOpen(false)}>
                Configurações
              </Button>
              <Button
                type="button"
                variant="secondary"
                onClick={handleLogout}
              >
                Logout
              </Button>
            </>
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
          <NavLink to="/pages" onClick={() => setOpen(false)}>Minhas páginas</NavLink>
          <NavLink to="/plans" onClick={() => setOpen(false)}>Planos</NavLink>
        </nav>
      </aside>
    </>
  );
}
