"use client";

import { useEffect } from "react";

import { useRouter } from "next/navigation";

import { useAuthStore } from "@/src/stores/auth.store";

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export default function ProtectedRoute({
  children,
}: ProtectedRouteProps) {

  const router = useRouter();

  const isAuthenticated = useAuthStore(
    (state) => state.isAuthenticated
  );

  const isCheckingAuth = useAuthStore(
    (state) => state.isCheckingAuth
  );

  useEffect(() => {

    if (!isCheckingAuth && !isAuthenticated) {
      router.replace("/Login");
    }

  }, [isAuthenticated, isCheckingAuth, router]);

  if (isCheckingAuth || !isAuthenticated) {

    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-100">

        <div className="text-center">

          <div className="mx-auto h-10 w-10 animate-spin rounded-full border-4 border-blue-600 border-t-transparent" />

          <p className="mt-4 text-sm text-gray-500">
            Checking authentication...
          </p>

        </div>

      </div>
    );

  }

  return <>{children}</>;
}