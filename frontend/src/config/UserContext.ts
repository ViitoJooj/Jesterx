import React, { createContext, useContext, useEffect, useRef, useState } from "react";
import { url } from "./Vars";

export type User = {
  id: string;
  profile_img: string;
  first_name: string;
  last_name: string;
  email: string;
  role: string;
  plan: string;
} | null;

type MeResponse = {
  success: boolean;
  message: string;
  data: {
    id: string;
    profile_img: string;
    first_name: string;
    last_name: string;
    email: string;
    role: string;
    plan: string;
  };
};

type UserContextValue = {
  user: User | undefined;
  setUser: React.Dispatch<React.SetStateAction<User | undefined>>;
};

const UserContext = createContext<UserContextValue | undefined>(undefined);

function toStr(v: unknown) {
  return v === null || v === undefined ? "" : String(v);
}

function readCachedUser(): User | undefined {
  const id = localStorage.getItem("userId");
  if (!id) return undefined;

  return {
    id,
    profile_img: localStorage.getItem("userProfileImg") || "",
    first_name: localStorage.getItem("userFirstName") || "",
    last_name: localStorage.getItem("userLastName") || "",
    email: localStorage.getItem("userEmail") || "",
    role: localStorage.getItem("userRole") || "",
    plan: localStorage.getItem("userPlan") || "",
  };
}

function writeCachedUser(u: Exclude<User, null>) {
  localStorage.setItem("userId", u.id);
  localStorage.setItem("userProfileImg", u.profile_img);
  localStorage.setItem("userFirstName", u.first_name);
  localStorage.setItem("userLastName", u.last_name);
  localStorage.setItem("userEmail", u.email);
  localStorage.setItem("userRole", u.role);
  localStorage.setItem("userPlan", u.plan);
}

function clearCachedUser() {
  localStorage.removeItem("userId");
  localStorage.removeItem("userProfileImg");
  localStorage.removeItem("userFirstName");
  localStorage.removeItem("userLastName");
  localStorage.removeItem("userEmail");
  localStorage.removeItem("userRole");
  localStorage.removeItem("userPlan");
}

function clearAllCookies() {
  if (typeof document === "undefined") return;
  const raw = document.cookie;
  if (!raw) return;
  const cookies = raw.split(";");
  for (const cookie of cookies) {
    const parts = cookie.split("=");
    const name = parts[0]?.trim();
    if (!name) continue;
    document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`;
  }
}

function handleUnauthorized() {
  let hadUser = false;
  try {
    if (typeof window !== "undefined" && window.localStorage) {
      hadUser = !!localStorage.getItem("userId");
    }
  } catch {
    hadUser = false;
  }
  clearCachedUser();
  clearAllCookies();
  if (hadUser && typeof window !== "undefined") {
    window.location.reload();
  }
}

if (typeof globalThis.fetch === "function") {
  const originalFetch = globalThis.fetch;
  globalThis.fetch = async (input: RequestInfo | URL, init?: RequestInit) => {
    const response = await originalFetch(input, init);
    if (response.status === 401) {
      handleUnauthorized();
    }
    return response;
  };
}

type UserProviderProps = {
  children: React.ReactNode;
};

export function UserProvider({ children }: UserProviderProps) {
  const [user, setUser] = useState<User | undefined>(undefined);
  const ran = useRef(false);

  useEffect(() => {
    if (ran.current) return;
    ran.current = true;

    const cached = readCachedUser();
    if (cached) {
      setUser(cached);
      return;
    }

    (async () => {
      try {
        const res = await fetch(`${url}/v1/auth/me`, {
          method: "GET",
          credentials: "include",
        });

        if (!res.ok) {
          setUser(null);
          return;
        }

        const payload = (await res.json()) as MeResponse;
        const data = payload?.data;

        if (!data?.id) {
          clearCachedUser();
          setUser(null);
          return;
        }

        const fresh = {
          id: toStr(data.id),
          profile_img: toStr(data.profile_img),
          first_name: toStr(data.first_name),
          last_name: toStr(data.last_name),
          email: toStr(data.email),
          role: toStr(data.role),
          plan: toStr(data.plan),
        };

        writeCachedUser(fresh);
        setUser(fresh);
      } catch {
        setUser(null);
      }
    })();
  }, []);

  return React.createElement(
    UserContext.Provider,
    { value: { user, setUser } },
    children
  );
}

export function useUser() {
  const ctx = useContext(UserContext);
  if (!ctx) throw new Error("useUser must be used within UserProvider");
  return ctx;
}
