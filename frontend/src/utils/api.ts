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

export async function apiRequest<T = any>(endpoint: string, options: RequestInit = {}): Promise<ApiResponse<T>> {
  const defaultHeaders: Record<string, string> = {
    "Content-Type": "application/json",
  };

  const tenant = getCurrentTenant();
  if (tenant) {
    defaultHeaders["X-Tenant-Page-Id"] = tenant;
  }

  const config: RequestInit = {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
    credentials: "include",
  };

  try {
    const response = await fetch(`${url}${endpoint}`, config);
    const data = await response.json();

    if (!response.ok) {
      throw {
        message: data.message || data.error || "Request failed",
        status: response.status,
      } as ApiError;
    }

    return data;
  } catch (error) {
    if ((error as ApiError).message) {
      throw error;
    }
    throw {
      message: "Network error.  Please try again later.",
    } as ApiError;
  }
}

export async function post<T = any>(endpoint: string, body: any): Promise<ApiResponse<T>> {
  return apiRequest<T>(endpoint, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export async function get<T = any>(endpoint: string): Promise<ApiResponse<T>> {
  return apiRequest<T>(endpoint, {
    method: "GET",
  });
}

export async function put<T = any>(endpoint: string, body: any): Promise<ApiResponse<T>> {
  return apiRequest<T>(endpoint, {
    method: "PUT",
    body: JSON.stringify(body),
  });
}

export async function del<T = any>(endpoint: string): Promise<ApiResponse<T>> {
  return apiRequest<T>(endpoint, {
    method: "DELETE",
  });
}
