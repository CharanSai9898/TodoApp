"use client";

import { useEffect } from "react";

import { useAuthStore } from "@/src/stores/auth.store";

interface AuthProviderProps {
  children: React.ReactNode;
}

export default function AuthProvider({
  children,
}: AuthProviderProps) {

  const checkAuth = useAuthStore(
    (state) => state.checkAuth
  );

  useEffect(() => {
    void checkAuth();
  }, [checkAuth]);

  return <>{children}</>;
}