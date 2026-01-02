import { useState, useMemo, useRef, useEffect } from "react";
import { createPortal } from "react-dom";
import { NavLink, Link, useNavigate } from "react-router-dom";
import styles from "../styles/components/Header.module.scss";
import buttonStyles from "../styles/components/Button.module.scss";
import { useUser } from "../config/UserContext";
import { url } from "../config/Vars";

export function Header() {
  const { user, setUser } = useUser();
  const [open, setOpen] = useState(false);
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const userMenuRef = useRef<HTMLDivElement | null>(null);
  const navigate = useNavigate();

  const initials = useMemo(() => {
    if (!user) return "A";
    const f = (user.first_name?.[0] ?? "").toUpperCase();
    const l = (user.last_name?.[0] ?? "").toUpperCase();
    return (f + l || f || user.id?.[0] || "U").toUpperCase();
  }, [user]);

  function getCookie(name: string) {
    return document.cookie
      .split("; ")
      .find((row) => row.startsWith(name + "="))
      ?.split("=")[1];
  }

  async function handleLogout() {
    try {
      const csrf = getCookie("csrf");
      await fetch(`${url}/v1/auth/logout`, {
        method: "GET",
        credentials: "include",
        headers: csrf ? { "X-CSRF-Token": csrf } : undefined,
      });
    } catch {}

    localStorage.removeItem("userId");
    localStorage.removeItem("userProfileImg");
    localStorage.removeItem("userFirstName");
    localStorage.removeItem("userLastName");
    localStorage.removeItem("userEmail");
    localStorage.removeItem("userRole");
    localStorage.removeItem("userPlan");

    setUser(null);
    navigate("/login", { replace: true });
  }

  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") {
        setUserMenuOpen(false);
        setOpen(false);
      }
    }

    function onClickOutside(e: MouseEvent) {
      if (userMenuRef.current && !userMenuRef.current.contains(e.target as Node)) {
        setUserMenuOpen(false);
      }
    }

    document.addEventListener("keydown", onKey);
    document.addEventListener("mousedown", onClickOutside);

    return () => {
      document.removeEventListener("keydown", onKey);
      document.removeEventListener("mousedown", onClickOutside);
    };
  }, []);

  useEffect(() => {
    if (!open) return;
    const prev = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = prev;
    };
  }, [open]);

  const drawerPortal =
    typeof document !== "undefined"
      ? createPortal(
          <>
            <div className={`${styles.drawer_backdrop} ${open ? styles.open : ""}`} onClick={() => setOpen(false)} aria-hidden />
            <aside id="drawer" className={`${styles.drawer} ${open ? styles.open : ""}`} aria-hidden={!open} role="dialog" aria-modal="true">
              <div className={styles.drawer_header}>
                <span>Menu</span>
                <button className={styles.drawer_close} onClick={() => setOpen(false)} aria-label="Fechar menu">
                  ×
                </button>
              </div>

              {user === null ? (
                <div className={styles.drawer_ctas_top}>
                  <Link to="/login" onClick={() => setOpen(false)} className={`${buttonStyles.default_button} ${buttonStyles["default_button--secondary"]}`}>
                    Login
                  </Link>
                  <Link to="/register" onClick={() => setOpen(false)} className={`${buttonStyles.default_button} ${buttonStyles["default_button--primary"]}`}>
                    Register
                  </Link>
                </div>
              ) : (
                <div className={styles.drawer_ctas_top}>
                  <button
                    onClick={() => {
                      setOpen(false);
                      handleLogout();
                    }}
                    className={`${buttonStyles.default_button} ${buttonStyles["default_button--secondary"]}`}
                  >
                    Logout
                  </button>
                </div>
              )}

              <nav className={styles.drawer_nav} aria-label="Drawer navigation">
                <NavLink to="/" end onClick={() => setOpen(false)}>
                  Home
                </NavLink>
                <NavLink to="/pages" onClick={() => setOpen(false)}>
                  My Store
                </NavLink>
                <NavLink to="/products" onClick={() => setOpen(false)}>
                  My Products
                </NavLink>
                <NavLink to="/api" onClick={() => setOpen(false)}>
                  API
                </NavLink>
                {user?.role === "platform_admin" && (
                  <NavLink to="/admin" onClick={() => setOpen(false)}>
                    Admin
                  </NavLink>
                )}
              </nav>
            </aside>
          </>,
          document.body
        )
      : null;

  if (user === undefined) {
    return (
      <header className={styles.header}>
        <div className={styles.inner}>
          <Link to="/" className={styles.brand} aria-label="Go to homepage">
            <span className={styles.logo}>J</span>
            <span className={styles.name}>Jester</span>
          </Link>
        </div>
      </header>
    );
  }

  return (
    <>
      <header className={styles.header}>
        <div className={styles.inner}>
          <Link to="/" className={styles.brand} aria-label="Go to homepage">
            <span className={styles.logo}>J</span>
            <span className={styles.name}>Jester</span>
          </Link>

          <nav className={styles.nav} aria-label="Primary">
            <NavLink to="/" end className={({ isActive }) => (isActive ? "active" : undefined)}>
              Home
            </NavLink>
            <NavLink to="/pages" className={({ isActive }) => (isActive ? "active" : undefined)}>
              My Store
            </NavLink>
              <NavLink to="/products" className={({ isActive }) => (isActive ? "active" : undefined)}>
                My Products
              </NavLink>
              <NavLink to="/api" className={({ isActive }) => (isActive ? "active" : undefined)}>
                API
              </NavLink>
              {user?.role === "platform_admin" && (
                <NavLink to="/admin" className={({ isActive }) => (isActive ? "active" : undefined)}>
                  Admin
                </NavLink>
              )}
            </nav>

          <div className={styles.user}>
            {user === null ? (
              <>
                <Link to="/login" className={`${buttonStyles.default_button} ${buttonStyles["default_button--secondary"]}`}>
                  Login
                </Link>
                <Link to="/register" className={`${buttonStyles.default_button} ${buttonStyles["default_button--primary"]}`}>
                  Register
                </Link>
              </>
            ) : (
              <div className={styles.user_menu} ref={userMenuRef}>
                <button className={styles.avatar_btn} aria-haspopup="menu" aria-expanded={userMenuOpen} onClick={() => setUserMenuOpen((v) => !v)}>
                  {user.profile_img ? (
                    <img className={styles.avatar_img} src={user.profile_img} alt={`${user.first_name} ${user.last_name}`} draggable={false} />
                  ) : (
                    <span className={styles.avatar_fallback} aria-hidden="true">
                      {initials}
                    </span>
                  )}
                </button>

                <span className={styles.avatar_name}>{user.first_name}</span>

                <div className={`${styles.user_menu_dropdown} ${userMenuOpen ? styles.open : ""}`} role="menu" aria-label="User menu">
                  <div className={styles.menu_section_label}>Conta</div>
                  <Link to={`/${user.id}/user`} role="menuitem" className={styles.user_menu_item}>
                    Perfil
                  </Link>
                  <Link to="/settings" role="menuitem" className={styles.user_menu_item}>
                    Configurações
                  </Link>

                  <div className={styles.menu_divider} />

                  <div className={styles.menu_section_label}>Loja</div>
                  <Link to="/store" role="menuitem" className={styles.user_menu_item}>
                    Minha Loja
                  </Link>
                  <Link to="/products" role="menuitem" className={styles.user_menu_item}>
                    Meus Produtos
                  </Link>

                  <div className={styles.menu_divider} />

                  <button role="menuitem" onClick={handleLogout} className={`${styles.user_menu_item_btn} ${styles.danger}`}>
                    Logout
                  </button>
                </div>
              </div>
            )}
          </div>

          <button className={styles.menu_btn} aria-expanded={open} aria-controls="drawer" onClick={() => setOpen((v) => !v)}>
            <svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true">
              <path d="M3 6h18M3 12h18M3 18h18" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
            </svg>
            Menu
          </button>
        </div>
      </header>

      {drawerPortal}
    </>
  );
}
