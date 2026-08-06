const API_BASE_URL = import.meta.env.VITE_API_URL ?? "";

export interface ApiError {
  error: string;
  message: string;
}

export class ApiClientError extends Error {
  status: number;
  code: string;

  constructor(status: number, code: string, message: string) {
    super(message);
    this.name = "ApiClientError";
    this.status = status;
    this.code = code;
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let body: ApiError | null = null;
    try {
      body = await response.json();
    } catch {
      // non-JSON error body, ignore
    }
    throw new ApiClientError(
      response.status,
      body?.error ?? "UNKNOWN",
      body?.message ?? response.statusText,
    );
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json() as Promise<T>;
}

function buildHeaders(init?: RequestInit): HeadersInit {
  const headers: Record<string, string> = {};

  if (!(init?.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }

  const websiteUuid =
    (typeof document !== "undefined" &&
      document.documentElement.dataset.websiteUuid) ??
    undefined;
  if (websiteUuid) {
    headers["X-Website-UUID"] = websiteUuid;
  }

  return headers;
}

export const api = {
  async get<T>(path: string, init?: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      ...init,
      method: "GET",
      headers: { ...buildHeaders(init), ...(init?.headers as Record<string, string>) },
      credentials: "include",
    });
    return handleResponse<T>(response);
  },

  async post<T>(path: string, body?: unknown, init?: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      ...init,
      method: "POST",
      headers: { ...buildHeaders(init), ...(init?.headers as Record<string, string>) },
      credentials: "include",
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
    return handleResponse<T>(response);
  },

  async put<T>(path: string, body?: unknown, init?: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      ...init,
      method: "PUT",
      headers: { ...buildHeaders(init), ...(init?.headers as Record<string, string>) },
      credentials: "include",
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
    return handleResponse<T>(response);
  },

  async delete<T>(path: string, init?: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      ...init,
      method: "DELETE",
      headers: { ...buildHeaders(init), ...(init?.headers as Record<string, string>) },
      credentials: "include",
    });
    return handleResponse<T>(response);
  },
};

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
  save_login: boolean;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterResponse {
  uuid: string;
  name: string;
  email: string;
  created_at: string;
}

export interface LoginResponse {
  token: string;
  token_type: string;
  expires_in: number;
  user_uuid: string;
  website_uuid: string;
  name: string;
  email: string;
}

export const authService = {
  register(data: RegisterRequest): Promise<RegisterResponse> {
    return api.post<RegisterResponse>("/auth/register", data);
  },

  login(data: LoginRequest): Promise<LoginResponse> {
    return api.post<LoginResponse>("/auth/login", data);
  },
};
