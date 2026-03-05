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
  data?: T;
};

type LoginData = { id: string; website_id: string; email: string };
type RegisterData = { id?: string; website_id?: string; email?: string };

type MeData = {
  id: string;
  website_id?: string;
  email: string;
  role?: string;
  user_plan?: string;
  first_name?: string;
  last_name?: string;
  avatar_url?: string;
};

export function useAuth(websiteId: WebsiteId) {
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [me, setMe] = useState<MeData | null>(null);
  const [loading, setLoading] = useState(false);

  const refreshPromiseRef = useRef<Promise<string> | null>(null);

  const isAuthenticated = useMemo(() => !!me, [me]);

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

        await refresh();
        await loadMe();

        return "cookie-session";
      } catch (error: any) {
        const message = error?.message || "Erro ao realizar login";
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
        const resp = await apiFetch<AuthResponse<RegisterData>>(
          "/api/v1/auth/register",
          {
            method: "POST",
            websiteId,
            body: JSON.stringify(req),
          }
        );
        return resp;
      } catch (error: any) {
        const message = error?.message || "Erro ao criar conta";
        throw new Error(message);
      } finally {
        setLoading(false);
      }
    },
    [websiteId]
  );

  /**
   * REFRESH
   */
  const refresh = useCallback(async () => {
    if (refreshPromiseRef.current) return refreshPromiseRef.current;

    refreshPromiseRef.current = (async () => {
      await apiFetch<AuthResponse<unknown>>("/api/v1/auth/refresh", {
        method: "GET",
        websiteId,
      });
      const token = `cookie-session-${Date.now()}`;
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
    async () => {
      const resp = await apiFetch<MeData>("/api/v1/auth/me", {
        method: "GET",
        websiteId,
      });

      setMe(resp);
      return resp;
    },
    [websiteId]
  );

  /**
   * LOGOUT
   */
  const logout = useCallback(async () => {
    await apiFetch<void>("/api/v1/auth/logout", {
      method: "GET",
      websiteId,
    });

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
        await refresh();
        if (cancelled) return;
        await loadMe();
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
    websiteId,
  };
}
