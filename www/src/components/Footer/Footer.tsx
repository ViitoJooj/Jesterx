import { Link, useNavigate } from "react-router-dom";
import { useI18n } from "../../i18n";
import { Icon, type IconName } from "../icons/Icon";
import { Button, Divider } from "../ui";
import "./Footer.scss";

interface FooterLink {
    label: string;
    to: string;
}

// MOCK: colunas de links — páginas ainda não existem (ver docs/mocks.md)
const mainPages: FooterLink[][] = [
    [
        { label: "Home", to: "/" },
        { label: "Design", to: "/design" },
        { label: "Sobre", to: "/about" },
        { label: "Blog", to: "/blog" },
    ],
    [
        { label: "Preços", to: "/pricing" },
        { label: "Contato", to: "/contact" },
        { label: "Carreiras", to: "/careers" },
        { label: "Integrações", to: "/integrations" },
    ],
];

const utilityPages: FooterLink[] = [
    { label: "Entrar", to: "/login" },
    { label: "Registrar", to: "/register" },
    { label: "Licenças", to: "/licenses" },
    { label: "Changelog", to: "/changelog" },
    { label: "404", to: "/404" },
];

// MOCK: redes sociais — perfis ainda não existem (ver docs/mocks.md)
const socials: { icon: IconName; label: string; href: string }[] = [
    { icon: "facebook", label: "Facebook", href: "https://facebook.com" },
    { icon: "twitter", label: "Twitter", href: "https://twitter.com" },
    { icon: "instagram", label: "Instagram", href: "https://instagram.com" },
    { icon: "linkedin", label: "LinkedIn", href: "https://linkedin.com" },
];

export function Footer() {
    const { t } = useI18n();
    const navigate = useNavigate();

    const copyright = t("footer.copyright").replace(
        "{year}",
        String(new Date().getFullYear()),
    );

    return (
        <footer className="jx-footer">
            <div className="jx-footer__cta">
                <div className="jx-footer__cta-text">
                    <h2 className="jx-footer__cta-title">{t("footer.ctaTitle")}</h2>
                    <p className="jx-footer__cta-subtitle">{t("footer.ctaSubtitle")}</p>
                </div>
                <Button variant="primary" size="md" iconRight="arrow-right" onClick={() => navigate("/register")}>
                    {t("footer.ctaButton")}
                </Button>
            </div>

            <Divider />

            <div className="jx-footer__grid">
                <div className="jx-footer__brand">
                    <Link to="/" className="jx-footer__brand-link">
                        <img src="/favicon.svg" alt="Jesterx" className="jx-footer__logo" />
                        <span className="jx-footer__name">Jesterx</span>
                    </Link>
                    <p className="jx-footer__description">{t("footer.description")}</p>
                </div>

                <nav className="jx-footer__col-group" aria-label={t("footer.mainPages")}>
                    <span className="jx-footer__col-title">{t("footer.mainPages")}</span>
                    <div className="jx-footer__col-pair">
                        {mainPages.map((col, i) => (
                            <ul className="jx-footer__col" key={i}>
                                {col.map((link) => (
                                    <li key={link.to}>
                                        <Link className="jx-footer__link" to={link.to}>
                                            {link.label}
                                        </Link>
                                    </li>
                                ))}
                            </ul>
                        ))}
                    </div>
                </nav>

                <nav className="jx-footer__col-group" aria-label={t("footer.utilityPages")}>
                    <span className="jx-footer__col-title">{t("footer.utilityPages")}</span>
                    <ul className="jx-footer__col">
                        {utilityPages.map((link) => (
                            <li key={link.to}>
                                <Link className="jx-footer__link" to={link.to}>
                                    {link.label}
                                </Link>
                            </li>
                        ))}
                    </ul>
                </nav>
            </div>

            <Divider />

            <div className="jx-footer__bottom">
                <span className="jx-footer__copyright">{copyright}</span>
                <div className="jx-footer__socials">
                    {socials.map((social) => (
                        <a
                            key={social.label}
                            className="jx-footer__social"
                            href={social.href}
                            target="_blank"
                            rel="noopener noreferrer"
                            aria-label={social.label}
                            title={social.label}
                        >
                            <Icon name={social.icon} size={16} />
                        </a>
                    ))}
                </div>
            </div>
        </footer>
    );
}
