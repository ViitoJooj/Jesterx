import { createContext, useContext, useState, type ReactNode } from "react";

export interface User {
    name: string;
    email: string;
    avatar?: string;
}

interface AuthContextValue {
    user: User | null;
    logout: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

const STORAGE_KEY = "auth-user";

function loadUser(): User | null {
    try {
        const raw = localStorage.getItem(STORAGE_KEY);
        return raw ? (JSON.parse(raw) as User) : null;
    } catch {
        return null;
    }
}

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(loadUser);

    const logout = () => {
        localStorage.removeItem(STORAGE_KEY);
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ user, logout }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth(): AuthContextValue {
    const ctx = useContext(AuthContext);
    if (!ctx) throw new Error("useAuth deve ser usado dentro de AuthProvider");
    return ctx;
}
