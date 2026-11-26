"use client";

import { useAuthStore } from "@/src/stores/auth.store";

import { useRouter } from "next/navigation";

export default function Navbar() {

  const router = useRouter();

  const user = useAuthStore(
    (state) => state.user
  );

  const logout = useAuthStore(
    (state) => state.logout
  );

  const handleLogout = async () => {

    await logout();

    router.replace("/Login");

  };

  return (

    <nav className="flex items-center justify-between border-b bg-white px-8 py-4">

      <h1 className="text-xl font-bold text-blue-600">
        Todo App
      </h1>

      <div className="flex items-center gap-5">

        <span className="text-sm text-gray-600">
          {user?.email}
        </span>

        <button
          onClick={handleLogout}
          className="rounded-lg bg-red-500 px-4 py-2 text-sm font-medium text-white transition hover:bg-red-600"
        >
          Logout
        </button>

      </div>

    </nav>

  );

}