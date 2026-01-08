import { url } from "../config/Vars";

export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
}

export interface ApiError {
  message: string;
  status?: number;
}

export function getCurrentTenant(): string | null {
  return localStorage.getItem("current_tenant");
}

export function setCurrentTenant(tenantPageId: string) {
  localStorage.setItem("current_tenant", tenantPageId);
}

export function clearCurrentTenant() {
  localStorage.removeItem("current_tenant");
}

export async function apiRequest<T = any>(endpoint: string, options: RequestInit = {}, useTenant: boolean = true): Promise<ApiResponse<T>> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (useTenant) {
    const tenant = getCurrentTenant();
    if (tenant) {
      headers["X-Tenant-Page-Id"] = tenant;
    }
  }

  const response = await fetch(`${url}${endpoint}`, {
    ...options,
    headers: {
      ...headers,
      ...options.headers,
    },
    credentials: "include",
  });

  const text = await response.text();
  let data: any = {};

  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = {};
    }
  }

  if (!response.ok) {
    const error = new Error(data.message || data.error || "Request failed") as ApiError;
    error.status = response.status;
    throw error;
  }

  return text ? (data as ApiResponse<T>) : ({ success: true } as ApiResponse<T>);
}

export function get<T = any>(endpoint: string) {
  return apiRequest<T>(endpoint, { method: "GET" }, true);
}

export function getPublic<T = any>(endpoint: string) {
  return apiRequest<T>(endpoint, { method: "GET" }, false);
}

export function post<T = any>(endpoint: string, body: any) {
  return apiRequest<T>(endpoint, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export function put<T = any>(endpoint: string, body: any) {
  return apiRequest<T>(endpoint, {
    method: "PUT",
    body: JSON.stringify(body),
  });
}

export function del<T = any>(endpoint: string) {
  return apiRequest<T>(endpoint, { method: "DELETE" });
}
