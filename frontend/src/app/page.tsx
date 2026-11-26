import Link from "next/link";

export default function HomePage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-6 bg-gray-100 p-6 text-center">
      <h1 className="text-4xl font-bold text-gray-900">Todo App</h1>
      <p className="text-gray-600">Organize your daily tasks in one place.</p>
      <nav className="flex gap-4">
        <Link
          href="/Login"
          className="rounded-lg bg-blue-600 px-5 py-3 font-semibold text-white hover:bg-blue-700"
        >
          Log in
        </Link>
        <Link
          href="/Register"
          className="rounded-lg border border-gray-300 bg-white px-5 py-3 font-semibold text-gray-900 hover:bg-gray-50"
        >
          Create account
        </Link>
      </nav>
    </main>
  );
}