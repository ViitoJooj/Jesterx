export const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";
export const DEFAULT_WEBSITE_ID = "00000000-0000-0000-0000-000000000001";

export type WebsiteId = string;

type ApiFetchInit = RequestInit & {
  websiteId: WebsiteId;
  accessToken?: string;
  includeJsonContentType?: boolean;
};

export function makeHeaders(
  websiteId: WebsiteId,
  accessToken?: string,
  includeJsonContentType = true
) {
  const resolvedWebsiteId = websiteId?.trim() || DEFAULT_WEBSITE_ID;

  const h: Record<string, string> = {
    "X-Website-Id": resolvedWebsiteId,
  };
  if (includeJsonContentType) h["Content-Type"] = "application/json";
  if (accessToken) h["Authorization"] = `Bearer ${accessToken}`;
  return h;
}

function hasJsonContentType(headers: Headers): boolean {
  const contentType = headers.get("content-type");
  return contentType?.includes("application/json") ?? false;
}

export async function apiFetch<T>(input: string, init: ApiFetchInit): Promise<T> {
  const {
    websiteId,
    accessToken,
    includeJsonContentType = true,
    ...rest
  } = init;

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

  if (res.status === 204) return undefined as T;

  if (hasJsonContentType(res.headers)) {
    return (await res.json()) as T;
  }

  const text = await res.text();
  return text as T;
}
