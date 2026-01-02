import { useState } from "react";
import styles from "../styles/pages/Register.module.scss";
import buttonStyles from "../styles/components/Button.module.scss";
import inputStyles from "../styles/components/Input.module.scss";
import Button from "../components/Button";
import Input from "../components/Input";
import { getOAuthUrl } from "../config/Vars";
import { url } from "../config/Vars";

export function Register() {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [keepMeLoggedIn, setKeepMeLoggedIn] = useState(false);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState<string>("");
  const [status, setStatus] = useState<"success" | "error" | null>(null);

  const handleCheckboxChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setKeepMeLoggedIn(event.target.checked);
  };

  function validate() {
    if (!firstName || !lastName || !email || !password || !confirmPassword) {
      setMessage("Please fill in all fields.");
      setStatus("error");
      return false;
    }
    if (password.length < 8) {
      setMessage("Password must be at least 8 characters.");
      setStatus("error");
      return false;
    }
    if (password !== confirmPassword) {
      setMessage("Passwords do not match.");
      setStatus("error");
      return false;
    }
    return true;
  }

  async function handleRegister(e?: React.FormEvent) {
    e?.preventDefault();
    setMessage("");
    setStatus(null);

    if (!validate()) return;

    try {
      setLoading(true);
      const res = await fetch(`${url}/v1/auth/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          first_name: firstName,
          last_name: lastName,
          email: email,
          password: password,
          keep_me_logged_in: keepMeLoggedIn,
        }),
      });

      const data = await res.json();

      if (res.ok) {
        setStatus("success");
        setMessage(data.message || "Account created successfully!");
        setTimeout(() => {
          window.location.href = "/";
        }, 400);
      } else {
        setStatus("error");
        setMessage(data.message || "Could not create your account.");
      }
    } catch (err) {
      setStatus("error");
      setMessage("Network error. Please try again later.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <h1>Welcome to our new store!</h1>

        <div className={styles.oauth_buttons}>
          <button type="button" onClick={() => window.location.href = getOAuthUrl("google")} className={styles.oauth_button}>
            <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M17.64 9.20454C17.64 8.56636 17.5827 7.95272 17.4764 7.36363H9V10.845H13.8436C13.635 11.97 13.0009 12.9231 12.0477 13.5613V15.8195H14.9564C16.6582 14.2527 17.64 11.9454 17.64 9.20454Z" fill="#4285F4"/>
              <path d="M9 18C11.43 18 13.4673 17.1941 14.9564 15.8195L12.0477 13.5613C11.2418 14.1013 10.2109 14.4204 9 14.4204C6.65591 14.4204 4.67182 12.8372 3.96409 10.71H0.957275V13.0418C2.43818 15.9831 5.48182 18 9 18Z" fill="#34A853"/>
              <path d="M3.96409 10.71C3.78409 10.17 3.68182 9.59318 3.68182 9C3.68182 8.40682 3.78409 7.83 3.96409 7.29V4.95818H0.957275C0.347727 6.17318 0 7.54772 0 9C0 10.4523 0.347727 11.8268 0.957275 13.0418L3.96409 10.71Z" fill="#FBBC05"/>
              <path d="M9 3.57955C10.3214 3.57955 11.5077 4.03364 12.4405 4.92545L15.0218 2.34409C13.4632 0.891818 11.4259 0 9 0C5.48182 0 2.43818 2.01682 0.957275 4.95818L3.96409 7.29C4.67182 5.16273 6.65591 3.57955 9 3.57955Z" fill="#EA4335"/>
            </svg>
            Continue with Google
          </button>
          <button type="button" onClick={() => window.location.href = getOAuthUrl("github")} className={styles.oauth_button}>
            <svg width="18" height="18" viewBox="0 0 18 18" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
              <path fillRule="evenodd" clipRule="evenodd" d="M9 0C4.0275 0 0 4.13211 0 9.22838C0 13.3065 2.5785 16.7648 6.15375 17.9841C6.60375 18.0709 6.76875 17.7853 6.76875 17.5403C6.76875 17.3212 6.76125 16.7405 6.7575 15.9712C4.254 16.5277 3.726 14.7332 3.726 14.7332C3.3165 13.6681 2.72475 13.3832 2.72475 13.3832C1.90875 12.8111 2.78775 12.8229 2.78775 12.8229C3.6915 12.887 4.16625 13.7737 4.16625 13.7737C4.96875 15.1847 6.273 14.777 6.7875 14.5414C6.8685 13.9443 7.10025 13.5381 7.3575 13.3065C5.35875 13.0748 3.258 12.2911 3.258 8.75524C3.258 7.74565 3.60825 6.92133 4.18425 6.27274C4.083 6.03748 3.77925 5.0985 4.263 3.82654C4.263 3.82654 5.01675 3.57808 6.738 4.77104C7.458 4.56516 8.223 4.46353 8.988 4.45837C9.753 4.46353 10.518 4.56516 11.238 4.77104C12.948 3.57808 13.7017 3.82654 13.7017 3.82654C14.1855 5.0985 13.8818 6.03748 13.7917 6.27274C14.3655 6.92133 14.7142 7.74565 14.7142 8.75524C14.7142 12.3007 12.6105 13.0713 10.608 13.2990C10.923 13.5797 11.2155 14.1304 11.2155 15.0034C11.2155 16.2435 11.2043 17.2455 11.2043 17.5403C11.2043 17.7868 11.3625 18.0753 11.823 17.9826C15.4237 16.7609 18 13.3045 18 9.22838C18 4.13211 13.9703 0 9 0Z"/>
            </svg>
            Continue with GitHub
          </button>
        </div>

        <div className={styles.divider}>
          <span>or</span>
        </div>

        <form onSubmit={handleRegister} noValidate>
          <div className={styles.input_container}>
            <div>
              <label htmlFor="firstName">First Name:</label>
              <Input id="firstName" value={firstName} onChange={(e) => setFirstName(e.target.value)} placeholder="John" type="text" autoComplete="given-name" className={inputStyles.default_input} required />
            </div>

            <div>
              <label htmlFor="lastName">Last Name:</label>
              <Input id="lastName" value={lastName} onChange={(e) => setLastName(e.target.value)} placeholder="Doe" type="text" autoComplete="family-name" className={inputStyles.default_input} required />
            </div>
          </div>

          <div className={`${styles.row} ${styles["row--full"]}`}>
            <div>
              <label htmlFor="email">Email:</label>
              <Input id="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="example@example.com" type="email" autoComplete="email" className={inputStyles.default_input} required />
            </div>
          </div>

          <div className={styles.input_container}>
            <div>
              <label htmlFor="password">Password:</label>
              <Input id="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="********" type="password" autoComplete="new-password" className={inputStyles.default_input} required />
            </div>

            <div>
              <label htmlFor="confirmPassword">Confirm Password:</label>
              <Input id="confirmPassword" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} placeholder="********" type="password" autoComplete="new-password" className={inputStyles.default_input} required />
            </div>
          </div>

          <div className={styles.checkbox_row}>
            <input id="keepMeLoggedIn" type="checkbox" checked={keepMeLoggedIn} onChange={handleCheckboxChange} />
            <label htmlFor="keepMeLoggedIn">Keep me logged in</label>
          </div>

          <Button type="submit" label={loading ? "Creating account…" : "Register"} disabled={loading} className={`${buttonStyles.default_button} ${buttonStyles["default_button--primary"]} ${styles.cta_button}`} />
        </form>

        {message && (
          <div className={`${styles.feedback} ${status === "success" ? styles["feedback--success"] : styles["feedback--error"]}`} role={status === "error" ? "alert" : "status"} aria-live={status === "error" ? "assertive" : "polite"}>
            {message}
          </div>
        )}
      </div>
    </main>
  );
}
