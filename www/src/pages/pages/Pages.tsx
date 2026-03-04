import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import Button from "../../components/button/Button";
import Input from "../../components/input/input";
import { useAuthContext } from "../../hooks/AuthContext";
import { apiFetch } from "../../hooks/api";
import styles from "./Pages.module.scss";

type WebsiteType = "landing_page" | "ecommerce" | "blog";

type CreatedWebsite = {
  id: string;
  type: string;
  name: string;
  short_description?: string;
  description?: string;
  created_at?: string;
};

type WebsiteResponse = {
  success: boolean;
  message: string;
  data: CreatedWebsite;
};

const typeOptions: { value: WebsiteType; label: string }[] = [
  { value: "landing_page", label: "Landing Page" },
  { value: "ecommerce", label: "E-commerce" },
  { value: "blog", label: "Blog" },
];

export const Pages: React.FC = () => {
  const navigate = useNavigate();
  const { isAuthenticated, websiteId, loading } = useAuthContext();

  const [type, setType] = useState<WebsiteType>("landing_page");
  const [name, setName] = useState("");
  const [shortDescription, setShortDescription] = useState("");
  const [description, setDescription] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [created, setCreated] = useState<CreatedWebsite | null>(null);

  const canSubmit = useMemo(() => {
    return isAuthenticated && name.trim().length >= 3 && !submitting;
  }, [isAuthenticated, name, submitting]);

  async function handleCreate(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!canSubmit) return;

    setError(null);
    setSubmitting(true);

    try {
      const payload = {
        type,
        name: name.trim(),
        short_description: shortDescription.trim(),
        description: description.trim(),
      };

      const resp = await apiFetch<WebsiteResponse>("/api/v1/websites", {
        method: "POST",
        websiteId,
        body: JSON.stringify(payload),
      });
      setCreated(resp.data);
      setName("");
      setShortDescription("");
      setDescription("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Falha ao criar página");
    } finally {
      setSubmitting(false);
    }
  }

  if (!loading && !isAuthenticated) {
    return (
      <main className={styles.main}>
        <section className={styles.container}>
          <h1>Área de páginas</h1>
          <p>Faça login para criar e gerenciar suas páginas via backend.</p>
          <div className={styles.actions}>
            <Button to="/login" variant="primary">
              Entrar
            </Button>
            <Button to="/register" variant="secondary">
              Criar conta
            </Button>
          </div>
        </section>
      </main>
    );
  }

  return (
    <main className={styles.main}>
      <section className={styles.container}>
        <div className={styles.header}>
          <h1>Minhas páginas</h1>
          <p>Crie uma nova página conectada diretamente ao endpoint `/api/v1/websites`.</p>
        </div>

        <form onSubmit={handleCreate} className={styles.form} noValidate>
          <div className={styles.fieldGroup}>
            <label htmlFor="websiteType">Tipo</label>
            <select
              id="websiteType"
              value={type}
              onChange={(e) => setType(e.target.value as WebsiteType)}
            >
              {typeOptions.map((item) => (
                <option key={item.value} value={item.value}>
                  {item.label}
                </option>
              ))}
            </select>
          </div>

          <div className={styles.fieldGroup}>
            <label htmlFor="websiteName">Nome do projeto</label>
            <Input
              id="websiteName"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Ex: Loja Fitness Aurora"
              minLength={3}
              required
            />
          </div>

          <div className={styles.fieldGroup}>
            <label htmlFor="shortDescription">Descrição curta</label>
            <Input
              id="shortDescription"
              value={shortDescription}
              onChange={(e) => setShortDescription(e.target.value)}
              placeholder="Uma frase objetiva para o projeto"
            />
          </div>

          <div className={styles.fieldGroup}>
            <label htmlFor="description">Descrição completa</label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Contexto, objetivo, público e diferenciais."
              rows={5}
            />
          </div>

          {error && <p className={styles.error}>{error}</p>}

          <div className={styles.actions}>
            <Button type="submit" variant="primary" disabled={!canSubmit}>
              {submitting ? "Criando..." : "Criar página"}
            </Button>
            <Button type="button" variant="secondary" onClick={() => navigate("/plans")}>
              Ver planos
            </Button>
          </div>
        </form>

        {created && (
          <article className={styles.result}>
            <h2>Projeto criado com sucesso</h2>
            <p>
              <strong>ID:</strong> {created.id}
            </p>
            <p>
              <strong>Nome:</strong> {created.name}
            </p>
            <p>
              <strong>Tipo:</strong> {created.type}
            </p>
          </article>
        )}
      </section>
    </main>
  );
};
