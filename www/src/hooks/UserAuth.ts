import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { apiFetch } from "./api";
import type { WebsiteId } from "./api";

type LoginRequest = { email: string; password: string };

type RegisterRequest = {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
};

type AuthResponse<T> = {
  success: boolean;
  message: string;
  data: T;
};

type LoginData = { id: string; websiteId: string; email: string };
type RefreshData =
  | { access_token: string }
  | { accessToken: string }
  | { token: string };

type MeData = {
  id: string;
  websiteId: string;
  email: string;
  role?: string;
  first_name?: string;
  last_name?: string;
};

function pickAccessToken(d: RefreshData): string {
  if ("access_token" in d) return d.access_token;
  if ("accessToken" in d) return d.accessToken;
  if ("token" in d) return d.token;
  throw new Error("refresh response missing access token");
}

export function useAuth(websiteId: WebsiteId) {
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [me, setMe] = useState<MeData | null>(null);
  const [loading, setLoading] = useState(false);

  const refreshPromiseRef = useRef<Promise<string> | null>(null);

  const isAuthenticated = useMemo(() => !!accessToken, [accessToken]);

  /**
   * LOGIN
   */
  const login = useCallback(
    async (req: LoginRequest) => {
      setLoading(true);
      try {
        await apiFetch<AuthResponse<LoginData>>(
          "/api/v1/auth/login",
          {
            method: "POST",
            websiteId,
            body: JSON.stringify(req),
          }
        );

        const token = await refresh();
        await loadMe(token);

        return token;
      } catch (error: any) {
        const message =
          error?.message || "Erro ao realizar login";
        throw new Error(message);
      } finally {
        setLoading(false);
      }
    },
    [websiteId]
  );

  /**
   * REGISTER
   */
  const register = useCallback(
    async (req: RegisterRequest) => {
      setLoading(true);
      try {
        await apiFetch<AuthResponse<unknown>>(
          "/api/v1/auth/register",
          {
            method: "POST",
            websiteId,
            body: JSON.stringify(req),
          }
        );

        // login automático após criar conta
        return await login({
          email: req.email,
          password: req.password,
        });
      } catch (error: any) {
        const message =
          error?.message || "Erro ao criar conta";
        throw new Error(message);
      } finally {
        setLoading(false);
      }
    },
    [websiteId, login]
  );

  /**
   * REFRESH
   */
  const refresh = useCallback(async () => {
    if (refreshPromiseRef.current) return refreshPromiseRef.current;

    refreshPromiseRef.current = (async () => {
      const resp = await apiFetch<AuthResponse<RefreshData>>(
        "/api/v1/auth/refresh",
        {
          method: "GET",
          websiteId,
        }
      );

      const token = pickAccessToken(resp.data);
      setAccessToken(token);
      return token;
    })();

    try {
      return await refreshPromiseRef.current;
    } finally {
      refreshPromiseRef.current = null;
    }
  }, [websiteId]);

  /**
   * LOAD ME
   */
  const loadMe = useCallback(
    async (token?: string) => {
      const t = token ?? accessToken;
      if (!t) throw new Error("no access token");

      const resp = await apiFetch<AuthResponse<MeData>>(
        "/api/v1/auth/me",
        {
          method: "GET",
          websiteId,
          accessToken: t,
        }
      );

      setMe(resp.data);
      return resp.data;
    },
    [websiteId, accessToken]
  );

  /**
   * LOGOUT
   */
  const logout = useCallback(async () => {
    await apiFetch<AuthResponse<unknown>>(
      "/api/v1/auth/logout",
      {
        method: "GET",
        websiteId,
      }
    );

    setAccessToken(null);
    setMe(null);
  }, [websiteId]);

  /**
   * AUTO REFRESH ON MOUNT
   */
  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const token = await refresh();
        if (cancelled) return;
        await loadMe(token);
      } catch {
        // silencioso
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [refresh, loadMe]);

  return {
    loading,
    isAuthenticated,
    accessToken,
    me,
    login,
    register,
    refresh,
    loadMe,
    logout,
    authHeader: accessToken
      ? { Authorization: `Bearer ${accessToken}` }
      : {},
  };
}