import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuthContext } from "../../hooks/AuthContext";

import styles from "./Register.module.scss";
import Input from "../../components/input/input";
import Button from "../../components/button/Button";

export const Register: React.FC = () => {
  const navigate = useNavigate();
  const { register, loading } = useAuthContext();

  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState<string | null>(null);

  const passwordsMatch = password === confirmPassword;

  const canSubmit =
    firstName.trim().length > 0 &&
    lastName.trim().length > 0 &&
    email.trim().length > 0 &&
    password.length > 0 &&
    confirmPassword.length > 0 &&
    passwordsMatch &&
    !loading;

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!canSubmit) return;

    setError(null);

    try {
      await register({
        first_name: firstName.trim(),
        last_name: lastName.trim(),
        email: email.trim(),
        password,
      });

      navigate("/", { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Falha no registro");
    }
  }

  return (
    <main className={styles.main}>
      <div className={styles.register_container}>
        <h1 className={styles.title}>Criar conta</h1>

          <div className={styles.oauth_row}>
            <button type="button" className={styles.oauth_box} disabled>
              {/* Google */}
              <svg width="20" height="20" viewBox="0 0 24 24">
                <path
                  fill="#EA4335"
                  d="M12 10.2v3.9h5.5c-.2 1.2-1.4 3.5-5.5 3.5-3.3 0-6-2.7-6-6s2.7-6 6-6c1.9 0 3.2.8 3.9 1.5l2.7-2.6C17.1 2.9 14.8 2 12 2 6.9 2 2.9 6 2.9 11s4 9 9.1 9c5.3 0 8.8-3.7 8.8-8.9 0-.6-.1-1.1-.2-1.6H12z"
                />
              </svg>
            </button>

            <button type="button" className={styles.oauth_box} disabled>
              {/* GitHub */}
              <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 .5C5.7.5.9 5.3.9 11.6c0 4.9 3.2 9.1 7.6 10.6.6.1.8-.3.8-.6v-2.2c-3.1.7-3.8-1.3-3.8-1.3-.5-1.2-1.1-1.5-1.1-1.5-.9-.6.1-.6.1-.6 1 .1 1.6 1 1.6 1 .9 1.6 2.3 1.1 2.9.9.1-.7.4-1.1.7-1.4-2.5-.3-5.1-1.2-5.1-5.5 0-1.2.4-2.2 1-3-.1-.3-.4-1.5.1-3.2 0 0 .9-.3 3 .9a10.4 10.4 0 0 1 5.4 0c2.1-1.2 3-.9 3-.9.5 1.7.2 2.9.1 3.2.6.8 1 1.8 1 3 0 4.3-2.6 5.2-5.1 5.5.4.3.8 1 .8 2v3c0 .3.2.7.8.6 4.4-1.5 7.6-5.7 7.6-10.6C23.1 5.3 18.3.5 12 .5z"/>
              </svg>
            </button>

            <button type="button" className={styles.oauth_box} disabled>
              {/* X */}
              <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                <path d="M18.9 2H22l-6.5 7.4L23.6 22h-6.6l-5.1-6.6L5.8 22H2.7l7-8L.4 2h6.7l4.6 6L18.9 2z"/>
              </svg>
            </button>
          </div>

        <div className={styles.divider}>
          <span>ou</span>
        </div>

        <form noValidate onSubmit={handleSubmit}>
          <div className={styles.form_group}>

            <div className={styles.row}>
              <div>
                <label htmlFor="first_name">Primeiro nome</label>
                <Input
                  id="first_name"
                  type="text"
                  autoComplete="given-name"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  required
                />
              </div>

              <div>
                <label htmlFor="last_name">Sobrenome</label>
                <Input
                  id="last_name"
                  type="text"
                  autoComplete="family-name"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  required
                />
              </div>
            </div>

            <div>
              <label htmlFor="email">Email</label>
              <Input
                id="email"
                type="email"
                autoComplete="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>

            <div className={styles.row}>
              <div>
                <label htmlFor="password">Senha</label>
                <Input
                  id="password"
                  type="password"
                  autoComplete="new-password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>

              <div>
                <label htmlFor="confirmPassword">Confirmar senha</label>
                <Input
                  id="confirmPassword"
                  type="password"
                  autoComplete="new-password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  required
                />
              </div>
            </div>

          </div>

          {!passwordsMatch && confirmPassword.length > 0 && (
            <p className={styles.error}>As senhas não coincidem</p>
          )}

          {error && <p className={styles.error}>{error}</p>}

          <Button
            type="submit"
            variant="primary"
            className={styles.cta_button}
            disabled={!canSubmit}
          >
            {loading ? "Criando conta..." : "Registrar"}
          </Button>
        </form>

        <div className={styles.links}>
          <p>
            Já tem conta? <a href="/login">Entrar</a>
          </p>
        </div>
      </div>
    </main>
  );
};