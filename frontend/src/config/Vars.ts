export const url = "http://localhost:8080";

export function getOAuthUrl(provider: "google" | "github" | "twitter"): string {
  return `${url}/v1/auth/${provider}`;
}
