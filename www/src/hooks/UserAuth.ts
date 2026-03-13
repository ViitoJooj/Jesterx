import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { apiFetch } from "./api";
import type { WebsiteId } from "./api";

type LoginRequest = { email: string; password: string };

type RegisterRequest = {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
  account_type: "personal" | "business";
  company_name?: string;
  trade_name?: string;
  cpf_cnpj?: string;
  phone?: string;
  zip_code?: string;
  address_street?: string;
  address_number?: string;
  address_complement?: string;
  address_city?: string;
  address_state?: string;
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
  cpf_cnpj?: string;
  avatar_url?: string;
  plan_max_sites?: number;
  plan_max_routes?: number;
  account_type?: "personal" | "business";
  company_name?: string;
  trade_name?: string;
  display_name?: string;
  birth_date?: string;
  gender?: string;
  bio?: string;
  instagram?: string;
  website_url?: string;
  whatsapp?: string;
  phone?: string;
  zip_code?: string;
  address_street?: string;
  address_number?: string;
  address_complement?: string;
  address_district?: string;
  address_city?: string;
  address_state?: string;
  address_country?: string;
};

export function useAuth(websiteId: WebsiteId) {
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [me, setMe] = useState<MeData | null>(null);
  const [loading, setLoading] = useState(false);

  const refreshPromiseRef = useRef<Promise<string> | null>(null);

  const isAuthenticated = useMemo(() => !!me, [me]);

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

  const updateProfile = useCallback(
    async (data: {
      first_name: string;
      last_name: string;
      cpf_cnpj?: string | null;
      avatar_url?: string | null;
      company_name?: string | null;
      trade_name?: string | null;
      display_name?: string | null;
      birth_date?: string | null;
      gender?: string | null;
      bio?: string | null;
      instagram?: string | null;
      website_url?: string | null;
      whatsapp?: string | null;
      phone?: string | null;
      zip_code?: string | null;
      address_street?: string | null;
      address_number?: string | null;
      address_complement?: string | null;
      address_district?: string | null;
      address_city?: string | null;
      address_state?: string | null;
      address_country?: string | null;
    }) => {
      await apiFetch<void>("/api/v1/auth/me", {
        method: "PATCH",
        websiteId,
        body: JSON.stringify(data),
      });
      await loadMe();
    },
    [websiteId, loadMe]
  );

  const cancelPlan = useCallback(async () => {
    await apiFetch<void>("/api/v1/payments/cancel", {
      method: "POST",
      websiteId,
    });
    await loadMe();
  }, [websiteId, loadMe]);

  const deleteAccount = useCallback(async () => {
    await apiFetch<void>("/api/v1/auth/me", {
      method: "DELETE",
      websiteId,
    });
    setAccessToken(null);
    setMe(null);
  }, [websiteId]);

  const logout = useCallback(async () => {
    await apiFetch<void>("/api/v1/auth/logout", {
      method: "GET",
      websiteId,
    });

    setAccessToken(null);
    setMe(null);
  }, [websiteId]);

  // refresh token and load user data on mount
  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        await refresh();
        if (cancelled) return;
        await loadMe();
      } catch {
        // silently ignore — user is just not logged in
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
    updateProfile,
    cancelPlan,
    deleteAccount,
    authHeader: accessToken
      ? { Authorization: `Bearer ${accessToken}` }
      : {},
    websiteId,
  };
}
