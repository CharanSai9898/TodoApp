import { create } from "zustand";
import axios from "axios";

import {
  loginUser,
  registerUser,
  logoutUser,
  getCurrentUser,
  verifyOtp as verifyOtpApi,
  resendOtp as resendOtpApi,
} from "@/src/services/auth.service";

import type { User, LoginResponse } from "@/src/types/auth.types";

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isCheckingAuth: boolean;
  error: string | null;

  login: (email: string, password: string) => Promise<LoginResponse>;
  register: (email: string, password: string) => Promise<void>;
  verifyOtp: (tempToken: string, otp: string) => Promise<void>;
  resendOtp: (tempToken: string) => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,
  isCheckingAuth: true,
  error: null,

  // LOGIN
  login: async (email, password) => {
    try {
      set({
        isLoading: true,
        error: null,
      });

      const response = await loginUser({
        email,
        password,
      });

      set({
        isLoading: false,
      });

      return response;
    } catch (error) {
      let message = "Login failed";
      if (axios.isAxiosError(error)) {
        message = (error.response?.data as { error?: string })?.error || message;
      }

      set({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: message,
      });

      throw error;
    }
  },

  // VERIFY OTP (2FA)
  verifyOtp: async (tempToken: string, otp: string) => {
    try {
      set({
        isLoading: true,
        error: null,
      });

      await verifyOtpApi(tempToken, otp);

      const response = await getCurrentUser();

      set({
        user: response.user,
        isAuthenticated: true,
        isLoading: false,
      });
    } catch (error) {
      let message = "Verification failed";
      if (axios.isAxiosError(error)) {
        message = (error.response?.data as { error?: string })?.error || message;
      }

      set({
        isLoading: false,
        error: message,
      });

      throw error;
    }
  },

  // RESEND OTP
  resendOtp: async (tempToken: string) => {
    try {
      set({ error: null });
      await resendOtpApi(tempToken);
    } catch (error) {
      let message = "Failed to resend code";
      if (axios.isAxiosError(error)) {
        message = (error.response?.data as { error?: string })?.error || message;
      }
      set({ error: message });
      throw error;
    }
  },

  // REGISTER
  register: async (email, password) => {
    try {
      set({
        isLoading: true,
        error: null,
      });

      await registerUser({
        email,
        password,
      });

      set({
        isLoading: false,
      });
    } catch (error) {
      let message = "Registration failed";
      if (axios.isAxiosError(error)) {
        message = (error.response?.data as { error?: string })?.error || message;
      }

      set({
        isLoading: false,
        error: message,
      });

      throw error;
    }
  },

  // LOGOUT
  logout: async () => {
    try {
      await logoutUser();
    } finally {
      set({
        user: null,
        isAuthenticated: false,
        isLoading: false,
        error: null,
      });
    }
  },

  // CHECK AUTHENTICATION
  checkAuth: async () => {
    try {
      set({
        isCheckingAuth: true,
      });

      const response = await getCurrentUser();

      set({
        user: response.user,
        isAuthenticated: true,
        isCheckingAuth: false,
      });
    } catch {
      set({
        user: null,
        isAuthenticated: false,
        isCheckingAuth: false,
      });
    }
  },
}));
