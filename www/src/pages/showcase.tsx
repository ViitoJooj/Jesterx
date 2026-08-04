import { useEffect, useState } from "react";
import {
    Avatar,
    Badge,
    Button,
    Checkbox,
    Divider,
    GlassCard,
    IconButton,
    IconTile,
    Input,
    ListRow,
    Modal,
    ProgressDots,
    Radio,
    SegmentedControl,
    Select,
    Skeleton,
    Slider,
    Textarea,
    Toast,
    Toggle,
    Tooltip,
} from "../components/ui";
import { Icon } from "../components/icons/Icon";
import "./showcase.scss";

export function Showcase() {
    const [toggles, setToggles] = useState({ dark: true, backup: false, notif: true });
    const [theme, setTheme] = useState("dark");
    const [radius, setRadius] = useState(12);
    const [modalOpen, setModalOpen] = useState(false);
    const [toastOpen, setToastOpen] = useState(false);
    const [check, setCheck] = useState(true);
    const [radio, setRadio] = useState("a");
    const [font, setFont] = useState("sf");

    useEffect(() => {
        const root = document.documentElement;
        if (theme === "auto") {
            const mq = window.matchMedia("(prefers-color-scheme: light)");
            const apply = () => {
                root.dataset.theme = mq.matches ? "light" : "dark";
            };
            apply();
            mq.addEventListener("change", apply);
            return () => mq.removeEventListener("change", apply);
        }
        root.dataset.theme = theme;
    }, [theme]);

    return (
        <div className="showcase">
            <GlassCard className="showcase__card" padding="lg">
                <h1>Componentes</h1>
                <p className="showcase__dim">Metal UI / Liquid Glass — preto e branco</p>

                <section>
                    <h2>Segmented</h2>
                    <SegmentedControl
                        value={theme}
                        onChange={setTheme}
                        options={[
                            { value: "auto", label: "Auto" },
                            { value: "light", label: "Light" },
                            { value: "dark", label: "Dark" },
                        ]}
                    />
                </section>

                <section>
                    <h2>Slider</h2>
                    <Slider
                        label="Corner radius"
                        min={4}
                        max={16}
                        value={radius}
                        onChange={setRadius}
                        marks={["4", "8", "12", "16"]}
                    />
                </section>

                <section>
                    <h2>Inputs</h2>
                    <Input label="Email" icon="user" placeholder="voce@exemplo.com" />
                    <Input icon="search" placeholder="Buscar..." />
                </section>

                <section>
                    <h2>List rows</h2>
                    <div className="showcase__stack">
                        <ListRow
                            icon="moon"
                            title="Dark mode"
                            subtitle="Mais facil para os olhos e bateria."
                            trailing={
                                <Toggle
                                    checked={toggles.dark}
                                    onChange={(v) => setToggles({ ...toggles, dark: v })}
                                    label="Dark mode"
                                />
                            }
                        />
                        <ListRow
                            icon="cloud"
                            title="Auto Backup"
                            subtitle="Salva seus dados na nuvem."
                            trailing={
                                <Toggle
                                    checked={toggles.backup}
                                    onChange={(v) => setToggles({ ...toggles, backup: v })}
                                    label="Auto Backup"
                                />
                            }
                        />
                        <ListRow
                            icon="bell"
                            title="Notificacoes"
                            subtitle="Alertas de atualizacoes."
                            trailing={<Badge variant="solid">3</Badge>}
                            onClick={() => setToggles({ ...toggles, notif: !toggles.notif })}
                        />
                        <ListRow
                            icon="lock"
                            title="Privacy Lock"
                            subtitle="Protege acoes com biometria."
                            trailing={<Icon name="arrow-right" size={18} />}
                            onClick={() => {}}
                        />
                    </div>
                </section>

                <section>
                    <h2>Botoes</h2>
                    <div className="showcase__row">
                        <Button variant="primary" size="lg">
                            Continue
                        </Button>
                        <Button variant="glass" icon="sparkles">
                            Glass
                        </Button>
                        <Button variant="ghost" iconRight="arrow-right">
                            Skip
                        </Button>
                        <IconButton icon="x" label="Fechar" />
                    </div>
                </section>

                <section>
                    <h2>Icon tiles</h2>
                    <div className="showcase__row">
                        <IconTile icon="mic" size="lg" />
                        <IconTile icon="video" size="lg" />
                        <IconTile icon="monitor" size="lg" />
                        <IconTile icon="wifi" size="lg" />
                        <IconTile icon="globe" size="lg" />
                    </div>
                </section>

                <Divider label="Novos" />

                <section>
                    <h2>Avatar</h2>
                    <div className="showcase__row">
                        <Avatar name="Viito Jooj" size="lg" status="online" />
                        <Avatar name="Ana Silva" />
                        <Avatar size="sm" />
                    </div>
                </section>

                <section>
                    <h2>Select e Textarea</h2>
                    <Select
                        label="Font"
                        value={font}
                        onChange={(e) => setFont(e.target.value)}
                        options={[
                            { value: "sf", label: "SF Pro" },
                            { value: "inter", label: "Inter" },
                            { value: "roboto", label: "Roboto" },
                        ]}
                    />
                    <Textarea label="Bio" placeholder="Conte sobre voce..." />
                </section>

                <section>
                    <h2>Checkbox e Radio</h2>
                    <div className="showcase__row">
                        <Checkbox checked={check} onChange={setCheck} label="Lembrar de mim" />
                        <Radio checked={radio === "a"} onChange={() => setRadio("a")} label="Opcao A" name="demo" />
                        <Radio checked={radio === "b"} onChange={() => setRadio("b")} label="Opcao B" name="demo" />
                    </div>
                </section>

                <section>
                    <h2>Progress, Tooltip, Skeleton</h2>
                    <ProgressDots steps={5} current={1} />
                    <div className="showcase__row">
                        <Tooltip content="Dica contextual">
                            <Button variant="glass" size="sm">
                                Passe o mouse
                            </Button>
                        </Tooltip>
                        <Skeleton circle width={40} height={40} />
                        <Skeleton width={120} height={14} />
                    </div>
                </section>

                <section>
                    <h2>Modal e Toast</h2>
                    <div className="showcase__row">
                        <Button variant="glass" onClick={() => setModalOpen(true)}>
                            Abrir modal
                        </Button>
                        <Button variant="glass" onClick={() => setToastOpen(true)}>
                            Mostrar toast
                        </Button>
                    </div>
                </section>
            </GlassCard>

            <Modal
                open={modalOpen}
                onClose={() => setModalOpen(false)}
                title="Just a minute..."
                subtitle="De acesso aos seus dispositivos"
            >
                <div className="showcase__stack">
                    <ListRow
                        icon="mic"
                        title="Microphone"
                        subtitle="Para todos te ouvirem"
                        trailing={<Toggle checked={toggles.notif} onChange={(v) => setToggles({ ...toggles, notif: v })} />}
                    />
                    <Button variant="primary" size="lg" onClick={() => setModalOpen(false)}>
                        Continue
                    </Button>
                </div>
            </Modal>

            <Toast open={toastOpen} message="Tema atualizado" onClose={() => setToastOpen(false)} />
        </div>
    );
}
