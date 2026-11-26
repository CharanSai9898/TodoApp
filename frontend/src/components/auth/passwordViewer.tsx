"use client";

import React, { useState, useCallback, forwardRef, useRef } from "react";
import { Eye, EyeOff } from "lucide-react";

export type PasswordViewerProps = Omit<
  React.InputHTMLAttributes<HTMLInputElement>,
  "type"
>;

export const PasswordViewer = forwardRef<HTMLInputElement, PasswordViewerProps>(
  ({ className = "", disabled = false, ...props }, forwardedRef) => {
    const [showPassword, setShowPassword] = useState<boolean>(false);
    const innerRef = useRef<HTMLInputElement>(null);

    // Merge forwarded ref with internal ref
    const setRef = useCallback(
      (node: HTMLInputElement | null) => {
        (innerRef as React.MutableRefObject<HTMLInputElement | null>).current =
          node;
        if (typeof forwardedRef === "function") {
          forwardedRef(node);
        } else if (forwardedRef) {
          (
            forwardedRef as React.MutableRefObject<HTMLInputElement | null>
          ).current = node;
        }
      },
      [forwardedRef]
    );

    return (
      <div className="relative w-full">
        <input
          {...props}
          ref={setRef}
          type={showPassword ? "text" : "password"}
          disabled={disabled}
          className={`w-full rounded-lg border border-gray-300 px-4 py-3 pr-11 outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-100 disabled:cursor-not-allowed disabled:bg-gray-100 ${className}`}
        />
        <button
          type="button"
          onClick={() => setShowPassword((prev) => !prev)}
          disabled={disabled}
          tabIndex={-1}
          aria-label={showPassword ? "Hide password" : "Show password"}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-800 disabled:cursor-not-allowed"
        >
          {showPassword ? (
            <EyeOff className="h-5 w-5" />
          ) : (
            <Eye className="h-5 w-5" />
          )}
        </button>
      </div>
    );
  }
);

PasswordViewer.displayName = "PasswordViewer";

export default PasswordViewer;
