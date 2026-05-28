import React, { createContext, useContext, useState, useCallback } from "react";
import { login as apiLogin, logout as apiLogout, isAuthenticated } from "../lib/auth";

interface AuthContextValue {
  isAuth: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuth, setIsAuth] = useState(isAuthenticated());

  const login = useCallback(async (email: string, password: string) => {
    await apiLogin(email, password);
    setIsAuth(true);
  }, []);

  const logout = useCallback(async () => {
    await apiLogout();
    setIsAuth(false);
  }, []);

  return (
    <AuthContext.Provider value={{ isAuth, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuthContext() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuthContext must be used within AuthProvider");
  return ctx;
}
