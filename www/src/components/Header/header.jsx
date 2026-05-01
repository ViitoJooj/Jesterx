import { useEffect, useRef, useState } from "react";
import { NavLink } from "react-router-dom";
import styles from "./header.module.scss";
import Button from "../Button/button";
import { Brand } from "../Brand/brand";

const NAV_LINKS = [
  { to: "/", end: true, label: "Início" },
  { to: "/pages", label: "Minhas páginas" },
  { to: "/plans", label: "Planos" },
];

export function Header() {
  const [theme, setTheme] = useState("light");
  const [open, setOpen] = useState(false);
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const userMenuRef = useRef(null);

  const me = {
    first_name: "João",
    avatar_url: null,
  };

  const firstName = me?.first_name?.trim() || "Perfil";
  const avatarUrl = me?.avatar_url;
  const avatarFallback = firstName.slice(0, 1).toUpperCase();

  useEffect(() => {
    const saved = window.localStorage.getItem("jx-theme");
    const initialTheme =
      saved === "light" || saved === "dark"
        ? saved
        : window.matchMedia("(prefers-color-scheme: dark)").matches
          ? "dark"
          : "light";

    setTheme(initialTheme);
    document.documentElement.setAttribute("data-theme", initialTheme);
  }, []);

  useEffect(() => {
    function onPointerDown(event) {
      if (!userMenuRef.current) return;
      if (!userMenuRef.current.contains(event.target)) {
        setUserMenuOpen(false);
      }
    }

    document.addEventListener("mousedown", onPointerDown);
    return () => document.removeEventListener("mousedown", onPointerDown);
  }, []);

  function toggleTheme() {
    setTheme((prev) => {
      const next = prev === "dark" ? "light" : "dark";
      document.documentElement.setAttribute("data-theme", next);
      window.localStorage.setItem("jx-theme", next);
      return next;
    });
  }

  function handleLogout() {
    setUserMenuOpen(false);
    setOpen(false);
    alert("Logout (visual apenas)");
  }

  return (
    <>
      <header className={styles.header}>
        <div className={styles.inner}>
          <Brand to="/" />

          <nav className={styles.nav}>
            {NAV_LINKS.map(({ to, end, label }) => (
              <NavLink key={to} to={to} end={end}>
                {label}
              </NavLink>
            ))}
          </nav>

          <div className={styles.actions}>
            <button
              type="button"
              className={styles.theme_toggle}
              onClick={toggleTheme}
              aria-label="Alternar tema"
              title={theme === "dark" ? "Ativar tema claro" : "Ativar tema escuro"}
            >
              {theme === "dark" ? "☀" : "☾"}
            </button>

            {me ? (
              <div className={styles.user_menu} ref={userMenuRef}>
                <button
                  type="button"
                  className={styles.user_trigger}
                  onClick={() => setUserMenuOpen((prev) => !prev)}
                >
                  <span className={styles.avatar_btn}>
                    {avatarUrl ? (
                      <img
                        src={avatarUrl}
                        alt={firstName}
                        className={styles.avatar_img}
                      />
                    ) : (
                      <span className={styles.avatar_fallback}>
                        {avatarFallback}
                      </span>
                    )}
                  </span>
                  <span className={styles.avatar_name}>{firstName}</span>
                </button>

                <div
                  className={`${styles.user_menu_dropdown} ${
                    userMenuOpen ? styles.open : ""
                  }`}
                >
                  <NavLink
                    to="/profile"
                    className={styles.user_menu_item}
                    onClick={() => setUserMenuOpen(false)}
                  >
                    Perfil
                  </NavLink>

                  <NavLink
                    to="/settings"
                    className={styles.user_menu_item}
                    onClick={() => setUserMenuOpen(false)}
                  >
                    Configurações
                  </NavLink>

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
                <Button to="/login" variant="secondary">
                  Entrar
                </Button>
                <Button to="/register" variant="primary">
                  Criar conta
                </Button>
              </>
            )}
          </div>

          <button className={styles.menu_btn} onClick={() => setOpen(true)}>
            ☰
          </button>
        </div>
      </header>

      <div
        className={`${styles.drawer_backdrop} ${open ? styles.open : ""}`}
        onClick={() => setOpen(false)}
      />

      <aside className={`${styles.drawer} ${open ? styles.open : ""}`}>
        <div className={styles.drawer_header}>
          <p>Menu</p>
          <button className={styles.drawer_close} onClick={() => setOpen(false)}>×</button>
        </div>

        <div className={styles.drawer_ctas_top}>
          <button
            type="button"
            className={styles.theme_toggle}
            onClick={toggleTheme}
            aria-label="Alternar tema"
          >
            {theme === "dark" ? "Tema claro ☀" : "Tema escuro ☾"}
          </button>

          {me ? (
            <>
              <p className={styles.drawer_user_label}>{firstName}</p>

              <Button
                to="/profile"
                variant="secondary"
                onClick={() => setOpen(false)}
              >
                Perfil
              </Button>

              <Button
                to="/settings"
                variant="secondary"
                onClick={() => setOpen(false)}
              >
                Configurações
              </Button>

              <Button type="button" variant="secondary" onClick={handleLogout}>
                Logout
              </Button>
            </>
          ) : (
            <>
              <Button
                to="/login"
                variant="secondary"
                onClick={() => setOpen(false)}
              >
                Entrar
              </Button>

              <Button
                to="/register"
                variant="primary"
                onClick={() => setOpen(false)}
              >
                Criar conta
              </Button>
            </>
          )}
        </div>

        <nav className={styles.drawer_nav}>
          {NAV_LINKS.map(({ to, end, label }) => (
            <NavLink key={to} to={to} end={end} onClick={() => setOpen(false)}>
              {label}
            </NavLink>
          ))}
        </nav>
      </aside>
    </>
  );
}
