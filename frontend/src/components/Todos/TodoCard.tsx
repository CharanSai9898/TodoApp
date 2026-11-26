"use client"
import { Todo } from "@/src/types/todo.types"
import { useState } from "react"

interface TodoCardProps {
    todo: Todo
}

export default function TodoCard({
    todo,
}: TodoCardProps) {
    const [isCompleted, setIsCompleted] = useState(todo.completed)
    return (
        <li className="flex items-center justify-between gap-2 rounded-lg border border-gray-300 px-4 py-3">

            <div className="flex items-center gap-3">
                <input
                    type="checkbox"
                    checked={isCompleted}
                    onChange={(e) => setIsCompleted(e.target.checked)}
                    className="h-5 w-5 cursor-pointer rounded border-gray-300 text-blue-600 focus:ring-2 focus:ring-blue-500"
                />
                <p
                    className={`text-sm leading-relaxed ${isCompleted ? "text-gray-500 line-through" : "text-slate-200"
                        }`}
                >
                    {todo.title}
                </p>
            </div>

            <div className="flex items-center gap-3">
                <button className="rounded-lg border border-gray-300 px-4 py-2 text-xs font-medium text-white transition hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500">
                    Edit
                </button>
                <button className="rounded-lg border border-gray-300 px-4 py-2 text-xs font-medium text-red-500 transition hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-red-500">
                    Delete
                </button>
            </div>

        </li>
    )
}
