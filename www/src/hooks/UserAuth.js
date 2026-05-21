import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { apiFetch } from "./api";

export function useAuth(websiteId) {
  const [accessToken, setAccessToken] = useState(null);
  const [me, setMe] = useState(null);
  const [loading, setLoading] = useState(false);

  const refreshPromiseRef = useRef(null);

  const isAuthenticated = useMemo(() => !!me, [me]);

  const login = useCallback(
    async (req) => {
      setLoading(true);
      try {
        await apiFetch("/api/v1/auth/login", {
          method: "POST",
          websiteId,
          body: JSON.stringify(req),
        });

        await refresh();
        await loadMe();

        return "cookie-session";
      } catch (error) {
        const message = error?.message || "Erro ao realizar login";
        throw new Error(message);
      } finally {
        setLoading(false);
      }
    },
    [websiteId]
  );

  const register = useCallback(
    async (req) => {
      setLoading(true);
      try {
        const resp = await apiFetch("/api/v1/auth/register", {
          method: "POST",
          websiteId,
          body: JSON.stringify(req),
        });
        return resp;
      } catch (error) {
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
      await apiFetch("/api/v1/auth/refresh", {
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
      const resp = await apiFetch("/api/v1/auth/me", {
        method: "GET",
        websiteId,
      });

      setMe(resp);
      return resp;
    },
    [websiteId]
  );

  const updateProfile = useCallback(
    async (data) => {
      await apiFetch("/api/v1/auth/me", {
        method: "PATCH",
        websiteId,
        body: JSON.stringify(data),
      });
      await loadMe();
    },
    [websiteId, loadMe]
  );

  const listAddresses = useCallback(async () => {
    return apiFetch("/api/v1/auth/addresses", { method: "GET", websiteId });
  }, [websiteId]);

  const createAddress = useCallback(async (data) => {
    return apiFetch("/api/v1/auth/addresses", {
      method: "POST",
      websiteId,
      body: JSON.stringify(data),
    });
  }, [websiteId]);

  const updateAddress = useCallback(async (id, data) => {
    return apiFetch(`/api/v1/auth/addresses/${id}`, {
      method: "PATCH",
      websiteId,
      body: JSON.stringify(data),
    });
  }, [websiteId]);

  const deleteAddress = useCallback(async (id) => {
    return apiFetch(`/api/v1/auth/addresses/${id}`, { method: "DELETE", websiteId });
  }, [websiteId]);

  const setDefaultAddress = useCallback(async (id) => {
    return apiFetch(`/api/v1/auth/addresses/${id}/default`, { method: "POST", websiteId });
  }, [websiteId]);

  const cancelPlan = useCallback(async () => {
    await apiFetch("/api/v1/payments/cancel", {
      method: "POST",
      websiteId,
    });
    await loadMe();
  }, [websiteId, loadMe]);

  const deleteAccount = useCallback(async () => {
    await apiFetch("/api/v1/auth/me", {
      method: "DELETE",
      websiteId,
    });
    setAccessToken(null);
    setMe(null);
  }, [websiteId]);

  const logout = useCallback(async () => {
    await apiFetch("/api/v1/auth/logout", {
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
    listAddresses,
    createAddress,
    updateAddress,
    deleteAddress,
    setDefaultAddress,
    cancelPlan,
    deleteAccount,
    authHeader: accessToken
      ? { Authorization: `Bearer ${accessToken}` }
      : {},
    websiteId,
  };
}
