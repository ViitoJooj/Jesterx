import { createContext, useContext } from "react";
import { useAuth } from "./UserAuth";

const AuthContext = createContext(null);

export function AuthProvider({ websiteId, children }) {
  const auth = useAuth(websiteId);
  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
}

export function useAuthContext() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuthContext must be used within <AuthProvider />");
  return ctx;
}
