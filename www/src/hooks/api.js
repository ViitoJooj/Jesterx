export const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";
export const DEFAULT_WEBSITE_ID = "00000000-0000-0000-0000-000000000001";

export function makeHeaders(websiteId, accessToken, includeJsonContentType = true) {
  const resolvedWebsiteId = websiteId?.trim() || DEFAULT_WEBSITE_ID;
  const h = { "X-Website-Id": resolvedWebsiteId };
  if (includeJsonContentType) h["Content-Type"] = "application/json";
  if (accessToken) h["Authorization"] = `Bearer ${accessToken}`;
  return h;
}

function hasJsonContentType(headers) {
  return headers.get("content-type")?.includes("application/json") ?? false;
}

export async function apiFetch(input, init) {
  const { websiteId, accessToken, includeJsonContentType = true, ...rest } = init;

  const res = await fetch(`${API_URL}${input}`, {
    ...rest,
    headers: {
      ...(rest.headers ?? {}),
      ...makeHeaders(websiteId, accessToken, includeJsonContentType),
    },
    credentials: "include",
  });

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(text || `HTTP ${res.status}`);
  }

  if (res.status === 204) return undefined;
  if (hasJsonContentType(res.headers)) return res.json();
  return res.text();
}
