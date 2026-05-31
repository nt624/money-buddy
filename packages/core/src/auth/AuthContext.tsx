"use client";

import { createContext, useContext, useEffect, useState } from "react";
import type { AuthPort, AuthUser } from "./types";

type AuthContextType = {
  user: AuthUser | null;
  loading: boolean;
  signIn: (email: string, password: string) => Promise<void>;
  signUp: (email: string, password: string) => Promise<void>;
  signInWithGoogle: () => Promise<void>;
  signOut: () => Promise<void>;
  getIdToken: () => Promise<string | null>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

/**
 * Shared auth provider. The platform supplies an AuthPort (Firebase Web on web,
 * Firebase on React Native for mobile); this component holds the reactive auth
 * state and exposes the same surface to both apps.
 */
export function AuthProvider({
  authPort,
  children,
}: {
  authPort: AuthPort;
  children: React.ReactNode;
}) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const unsubscribe = authPort.onAuthStateChanged((nextUser) => {
      setUser(nextUser);
      setLoading(false);
    });
    return unsubscribe;
  }, [authPort]);

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        signIn: (email, password) => authPort.signIn(email, password),
        signUp: (email, password) => authPort.signUp(email, password),
        signInWithGoogle: () => authPort.signInWithGoogle(),
        signOut: () => authPort.signOut(),
        getIdToken: () => authPort.getToken(),
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
}
