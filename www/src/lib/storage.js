import { API_URL } from "../hooks/api";

export function resolveMediaUrl(url) {
  if (!url) return undefined;
  return url.startsWith("/") ? `${API_URL}${url}` : url;
}

export async function uploadImage(file, websiteId) {
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
  return resolveMediaUrl(json.data.url);
}
