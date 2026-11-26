"use client";

import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { ShieldCheck, ArrowLeft, RefreshCw } from "lucide-react";

import { loginSchema, LoginInput } from "@/src/schemas/auth.schema";
import { useAuthStore } from "@/src/stores/auth.store";
import PasswordViewer from "@/src/components/auth/passwordViewer";
import GoogleSignInButton from "@/src/components/auth/GoogleSignInButton";

const LOGIN_BACKGROUND_GIF = "/1.gif";

export default function LoginForm() {
  const router = useRouter();

  const login = useAuthStore((state) => state.login);
  const verifyOtp = useAuthStore((state) => state.verifyOtp);
  const resendOtp = useAuthStore((state) => state.resendOtp);
  const isLoading = useAuthStore((state) => state.isLoading);
  const error = useAuthStore((state) => state.error);

  // 2FA state
  const [step, setStep] = useState<"credentials" | "otp">("credentials");
  const [tempToken, setTempToken] = useState<string>("");
  const [userEmail, setUserEmail] = useState<string>("");
  const [otp, setOtp] = useState<string>("");
  const [resendCooldown, setResendCooldown] = useState<number>(0);
  const [feedbackMessage, setFeedbackMessage] = useState<string | null>(null);

  // Countdown timer for resend
  useEffect(() => {
    if (resendCooldown <= 0) return;
    const timer = setInterval(() => {
      setResendCooldown((prev) => prev - 1);
    }, 1000);
    return () => clearInterval(timer);
  }, [resendCooldown]);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginInput>({
    resolver: zodResolver(loginSchema),
  });

  const onCredentialsSubmit = async (data: LoginInput) => {
    setFeedbackMessage(null);
    try {
      const res = await login(data.email, data.password);
      if (res?.requires2FA && res?.tempToken) {
        setTempToken(res.tempToken);
        setUserEmail(res.email || data.email);
        setStep("otp");
        setResendCooldown(30);
      } else {
        router.push("/Dashboard");
      }
    } catch {
      // Error is handled in Zustand
    }
  };

  const onOtpSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!otp || otp.trim().length !== 6) return;
    setFeedbackMessage(null);

    try {
      await verifyOtp(tempToken, otp.trim());
      router.push("/Dashboard");
    } catch {
      // Error is handled in Zustand
    }
  };

  const handleResend = async () => {
    if (resendCooldown > 0 || !tempToken) return;
    try {
      await resendOtp(tempToken);
      setFeedbackMessage("A new verification code has been sent!");
      setResendCooldown(45);
    } catch {
      // Error handled in store
    }
  };

  return (
    <div
      className="relative flex min-h-screen items-center justify-center bg-cover bg-center bg-no-repeat px-4"
      style={{ backgroundImage: `url("${LOGIN_BACKGROUND_GIF}")` }}
    >
      <div
        aria-hidden="true"
        className="absolute inset-0 bg-slate-950/45"
      />

      {/* <div className="relative w-full max-w-md rounded-2xl bg-white p-8 shadow-lg"> */}
     {/* <div className="w-full max-w-md rounded-2xl border border-white/20 bg-white/0 p-8 shadow-2xl" > */}
     <div className="relative z-10 w-full max-w-md rounded-2xl border border-white/20 bg-slate-950/60 p-8 shadow-2xl">
        {step === "credentials" ? (
          <>
            <h1 className="mb-2 text-center text-3xl font-bold text-yellow-400">
              Welcome Back
            </h1>

            <p className="mb-8 text-center text-yellow-200">
              Login to your account
            </p>

            {error && (
              <p
                role="alert"
                className="mb-5 rounded-lg bg-red-50 p-3 text-sm text-red-600"
              >
                {error}
              </p>
            )}

            <form
              onSubmit={handleSubmit(onCredentialsSubmit)}
              className="space-y-5"
            >
              <div>
                <label className="mb-2 block text-sm font-medium text-green-100">
                  Email
                </label>

                <input
                  {...register("email")}
                  type="email"
                  placeholder="Enter your email"
                  //className="w-full rounded-lg border border-gray-300 px-4 py-3 outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-100 placeholder:text-green-100"
                  className="w-full rounded-lg border border-gray-300 bg-slate-900/70 px-4 py-3 text-white outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-100 placeholder:text-green-100"
                />

                {errors.email && (
                  <p className="mt-1 text-sm text-red-500">
                    {errors.email.message}
                  </p>
                )}
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-green-100">
                  Password
                </label>

                <PasswordViewer
                  {...register("password")}
                  placeholder="Enter your password"
                  className="placeholder:text-green-100 text-white"
                />

                {errors.password && (
                  <p className="mt-1 text-sm text-red-500">
                    {errors.password.message}
                  </p>
                )}
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="w-full rounded-lg bg-blue-600 py-3 font-semibold text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {isLoading ? "Checking credentials..." : "Login"}
              </button>
            </form>

            <div className="relative my-6 flex items-center justify-center">
              <div className="w-full border-t border-gray-200" />
              <span className="absolute bg-white px-3 text-xs uppercase text-gray-500">
                or
              </span>
            </div>

            <GoogleSignInButton text="Sign in with Google" />

            <p className="mt-6 text-center text-sm text-gray-600">
              Don&apos;t have an account?{" "}
              <Link
                href="/Register"
                className="font-semibold text-blue-600 hover:underline"
              >
                Create one
              </Link>
            </p>
          </>
        ) : (
          /* STEP 2: 2FA OTP ENTRY */
          <div>
            <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-blue-50 text-blue-600">
              <ShieldCheck className="h-8 w-8" />
            </div>

            <h1 className="mb-2 text-center text-2xl font-bold text-gray-900">
              Two-Factor Authentication
            </h1>

            <p className="mb-1 text-center text-sm text-gray-500">
              Enter the 6-digit verification code sent to:
            </p>
            <p className="mb-6 text-center text-sm font-semibold text-gray-800">
              {userEmail}
            </p>

            {error && (
              <p
                role="alert"
                className="mb-4 rounded-lg bg-red-50 p-3 text-center text-sm text-red-600"
              >
                {error}
              </p>
            )}

            {feedbackMessage && (
              <p className="mb-4 rounded-lg bg-green-50 p-3 text-center text-sm text-green-700">
                {feedbackMessage}
              </p>
            )}

            <form onSubmit={onOtpSubmit} className="space-y-5">
              <div>
                <input
                  type="text"
                  maxLength={6}
                  value={otp}
                  onChange={(e) =>
                    setOtp(e.target.value.replace(/\D/g, ""))
                  }
                  placeholder="······"
                  autoFocus
                  className="w-full rounded-lg border border-gray-300 py-3 text-center font-mono text-3xl font-bold tracking-[0.5em] text-gray-800 outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-100"
                />
                <p className="mt-2 text-center text-xs text-gray-400">
                  Code expires in 5 minutes
                </p>
              </div>

              <button
                type="submit"
                disabled={isLoading || otp.length !== 6}
                className="w-full rounded-lg bg-blue-600 py-3 font-semibold text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {isLoading ? "Verifying..." : "Verify Code"}
              </button>
            </form>

            <div className="mt-6 flex flex-col items-center gap-3">
              <button
                type="button"
                onClick={handleResend}
                disabled={resendCooldown > 0 || isLoading}
                className="inline-flex items-center gap-2 text-sm font-medium text-blue-600 hover:text-blue-700 disabled:cursor-not-allowed disabled:text-gray-400"
              >
                <RefreshCw className="h-4 w-4" />
                {resendCooldown > 0
                  ? `Resend code in ${resendCooldown}s`
                  : "Resend verification code"}
              </button>

              <button
                type="button"
                onClick={() => {
                  setStep("credentials");
                  setOtp("");
                  setFeedbackMessage(null);
                }}
                className="inline-flex items-center gap-2 text-sm text-gray-500 hover:text-gray-800"
              >
                <ArrowLeft className="h-4 w-4" />
                Back to login
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
