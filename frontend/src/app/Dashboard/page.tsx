"use client";

import ProtectedRoute from "@/src/components/auth/protectedRoute";
import Navbar from "@/src/components/Common/Navbar";

import { useEffect } from "react";

import { useTodoStore } from "@/src/stores/todo.store";

import TodoCard from "@/src/components/Todos/TodoCard";

function DashboardContent() {

  const todos = useTodoStore(
    (state) => state.todos
  );

  const fetchTodos = useTodoStore(
    (state) => state.fetchTodos
  );

  useEffect(() => {
    void fetchTodos();
  }, [fetchTodos]);

  return (

    <main className="min-h-screen bg-gray-100 p-8">

      <div className="mx-auto max-w-5xl">

        <h1 className="text-3xl font-bold text-gray-900">
          Dashboard
        </h1>

        <p className="mt-2 text-gray-500">
          Welcome to your Todo application.
        </p>
        
          
        <div className="mt-8 space-y-4">

          {todos.map((todo) => (

            <div
              key={todo.id}
              className="rounded-xl border bg-white p-5 shadow-sm"
            >

              <h2 className="text-lg font-semibold">
                {todo.title}
              </h2>

              <p className="mt-2 text-gray-500">
                {todo.description}
              </p>
             

            </div>

          ))}

        </div>

      </div>

    </main>

  );
}

export default function DashboardPage() {

  return (
    <ProtectedRoute>
      <Navbar />
      <DashboardContent />
    </ProtectedRoute>
  );

}