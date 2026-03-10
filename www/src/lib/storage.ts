import { API_URL } from "../hooks/api";

/**
 * Resolves a media URL that may be a relative backend path (/files/...) or
 * an absolute URL (legacy Supabase URLs, external images, etc.).
 */
export function resolveMediaUrl(url: string | null | undefined): string | undefined {
  if (!url) return undefined;
  return url.startsWith("/") ? `${API_URL}${url}` : url;
}

export async function uploadImage(file: File, websiteId: string): Promise<string> {
  const formData = new FormData();
  formData.append("file", file);

  const res = await fetch(`${API_URL}/api/v1/upload`, {
    method: "POST",
    headers: { "X-Website-Id": websiteId },
    body: formData,
    credentials: "include",
  });

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(text || `Erro no upload: HTTP ${res.status}`);
  }

  const json = await res.json();
  const url: string = json.data.url;
  return resolveMediaUrl(url)!;
}
