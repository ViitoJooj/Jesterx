export const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export type WebsiteId = string;

export function makeHeaders(websiteId: WebsiteId, accessToken?: string) {
  const h: Record<string, string> = {
    "Content-Type": "application/json",
    "X-Website-Id": websiteId,
  };
  if (accessToken) h["Authorization"] = `Bearer ${accessToken}`;
  return h;
}

export async function apiFetch<T>(
  input: string,
  init: RequestInit & { websiteId: WebsiteId; accessToken?: string }
): Promise<T> {
  const { websiteId, accessToken, ...rest } = init;

  const res = await fetch(`${API_URL}${input}`, {
    ...rest,
    headers: {
      ...(rest.headers ?? {}),
      ...makeHeaders(websiteId, accessToken),
    },
    credentials: "include",
  });

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(text || `HTTP ${res.status}`);
  }

  return (await res.json()) as T;
}