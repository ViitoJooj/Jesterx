import { useNavigate } from "react-router-dom";
import { Button, GlassCard } from "../components/ui";
import { useI18n } from "../i18n";
import "./not-found.scss";

export function NotFound() {
    const { t } = useI18n();
    const navigate = useNavigate();

    return (
        <div className="not-found">
            <span className="not-found__code" aria-hidden="true">
                {t("notFound.code")}
            </span>
            <GlassCard className="not-found__card" padding="lg">
                <h1 className="not-found__title">{t("notFound.title")}</h1>
                <p className="not-found__description">{t("notFound.description")}</p>
                <Button
                    variant="primary"
                    iconRight="arrow-right"
                    onClick={() => navigate("/")}
                >
                    {t("notFound.backHome")}
                </Button>
            </GlassCard>
        </div>
    );
}
