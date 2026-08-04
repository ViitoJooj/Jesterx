import { useEffect, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../../auth";
import { useI18n } from "../../i18n";
import { useTheme } from "../../theme";
import { Icon } from "../icons/Icon";
import { Avatar, Button, Divider, IconButton } from "../ui";
import "./Header.scss";

export function Header() {
    const { t, locale, setLocale } = useI18n();
    const { theme, toggleTheme } = useTheme();
    const { user, logout } = useAuth();
    const navigate = useNavigate();

    const [openMenu, setOpenMenu] = useState<"notifications" | "user" | null>(null);
    const notificationsRef = useRef<HTMLDivElement>(null);
    const userMenuRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (!openMenu) return;

        const onPointerDown = (e: PointerEvent) => {
            const target = e.target as Node;
            if (
                !notificationsRef.current?.contains(target) &&
                !userMenuRef.current?.contains(target)
            ) {
                setOpenMenu(null);
            }
        };
        const onKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Escape") setOpenMenu(null);
        };

        document.addEventListener("pointerdown", onPointerDown);
        document.addEventListener("keydown", onKeyDown);
        return () => {
            document.removeEventListener("pointerdown", onPointerDown);
            document.removeEventListener("keydown", onKeyDown);
        };
    }, [openMenu]);

    const toggleLanguage = () => setLocale(locale === "pt-br" ? "en" : "pt-br");

    const toggleMenu = (menu: "notifications" | "user") =>
        setOpenMenu((prev) => (prev === menu ? null : menu));

    const go = (path: string) => {
        setOpenMenu(null);
        navigate(path);
    };

    return (
        <header className="jx-header">
            <Link to="/" className="jx-header__brand">
                <img src="/favicon.svg" alt="Jesterx" className="jx-header__logo" />
                <span className="jx-header__name">Jesterx</span>
            </Link>

            <div className="jx-header__actions">
                {user && (
                    <div className="jx-header__notifications" ref={notificationsRef}>
                        <IconButton
                            icon="bell"
                            label={t("header.notifications")}
                            onClick={() => toggleMenu("notifications")}
                            aria-haspopup="menu"
                            aria-expanded={openMenu === "notifications"}
                        />

                        {openMenu === "notifications" && (
                            <div className="jx-header__menu jx-header__menu--notifications" role="menu">
                                <span className="jx-header__menu-title">{t("header.notifications")}</span>
                                <Divider />
                                <div className="jx-header__menu-empty">
                                    <Icon name="bell" size={24} />
                                    <span>{t("header.noNotifications")}</span>
                                </div>
                            </div>
                        )}
                    </div>
                )}
                <IconButton
                    icon={theme === "dark" ? "sun" : "moon"}
                    label={t("header.toggleTheme")}
                    onClick={toggleTheme}
                />
                <IconButton
                    icon="languages"
                    label={t("header.toggleLanguage")}
                    onClick={toggleLanguage}
                />

                {user ? (
                    <div className="jx-header__user" ref={userMenuRef}>
                        <button
                            className="jx-header__avatar-btn"
                            onClick={() => toggleMenu("user")}
                            aria-haspopup="menu"
                            aria-expanded={openMenu === "user"}
                            aria-label={user.name}
                        >
                            <Avatar src={user.avatar} name={user.name} size="sm" />
                        </button>

                        {openMenu === "user" && (
                            <div className="jx-header__menu" role="menu">
                                <div className="jx-header__menu-user">
                                    <Avatar src={user.avatar} name={user.name} size="md" />
                                    <div className="jx-header__menu-user-info">
                                        <span className="jx-header__menu-user-name">{user.name}</span>
                                        <span className="jx-header__menu-user-email">{user.email}</span>
                                    </div>
                                </div>

                                <Divider />

                                <button className="jx-header__menu-item" role="menuitem" onClick={() => go("/profile")}>
                                    <Icon name="user" size={16} />
                                    <span>{t("header.profile")}</span>
                                </button>
                                <button className="jx-header__menu-item" role="menuitem" onClick={() => go("/settings")}>
                                    <Icon name="settings" size={16} />
                                    <span>{t("header.settings")}</span>
                                </button>
                                <button className="jx-header__menu-item" role="menuitem" onClick={() => go("/login?add=true")}>
                                    <Icon name="user-plus" size={16} />
                                    <span>{t("header.addAccount")}</span>
                                </button>

                                <Divider />

                                <button
                                    className="jx-header__menu-item jx-header__menu-item--danger"
                                    role="menuitem"
                                    onClick={() => {
                                        setOpenMenu(null);
                                        logout();
                                    }}
                                >
                                    <Icon name="log-out" size={16} />
                                    <span>{t("header.logout")}</span>
                                </button>
                            </div>
                        )}
                    </div>
                ) : (
                    <div className="jx-header__guest">
                        <Button variant="ghost" size="sm" icon="log-in" onClick={() => go("/login")}>
                            {t("header.login")}
                        </Button>
                        <Button variant="primary" size="sm" onClick={() => go("/register")}>
                            {t("header.register")}
                        </Button>
                    </div>
                )}
            </div>
        </header>
    );
}
