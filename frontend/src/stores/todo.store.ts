import { create } from "zustand";
import type { Todo } from "@/src/types/todo.types";

interface TodoState {
  todos: Todo[];
  isLoading: boolean;
  error: string | null;
  fetchTodos: () => Promise<void>;
}

export const useTodoStore = create<TodoState>((set) => ({
  todos: [],
  isLoading: false,
  error: null,
  fetchTodos: async () => {
    set({ todos: [] });
  },
}));
