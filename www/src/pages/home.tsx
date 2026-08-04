import { useNavigate } from "react-router-dom";
import { SideRays } from "../components/SideRays/SideRays";
import { Badge, Button } from "../components/ui";
import { useI18n } from "../i18n";
import "./home.scss";

export function Home() {
    const { t } = useI18n();
    const navigate = useNavigate();

    return (
        <div className="home">
            <div className="home__rays">
                <SideRays speed={1} rayColor1="#ffffff" rayColor2="#8a8a8a" intensity={1.6} spread={1.2} origin="top-left" saturation={0} blend={0.6} falloff={1.4} />
            </div>

            <main className="home__hero">
                <Badge variant="glass">{t("home.badge")}</Badge>
                <h1 className="home__title">{t("home.title")}</h1>
                <p className="home__subtitle">{t("home.subtitle")}</p>
                <div className="home__actions">
                    <Button variant="primary" size="lg" iconRight="arrow-right">
                        {t("home.ctaPrimary")}
                    </Button>
                    <Button variant="glass" size="lg" onClick={() => navigate("/design")}>
                        {t("home.ctaSecondary")}
                    </Button>
                </div>
            </main>
        </div>
    );
}
