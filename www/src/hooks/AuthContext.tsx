import React, { createContext, useContext } from "react";
import type { WebsiteId } from "./api";
import { useAuth } from "./UserAuth";

type AuthContextValue = ReturnType<typeof useAuth>;

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({
  websiteId,
  children,
}: React.PropsWithChildren<{ websiteId: WebsiteId }>) {
  const auth = useAuth(websiteId);
  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
}

export function useAuthContext() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuthContext must be used within <AuthProvider />");
  return ctx;
}