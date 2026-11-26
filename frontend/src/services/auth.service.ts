import api from "@/src/lib/axios";
import { AuthInput, LoginResponse } from "@/src/types/auth.types";

// Register
export const registerUser = async (data: AuthInput) => {
  const response = await api.post("/auth/register", data);
  return response.data;
};

// Login
export const loginUser = async (data: AuthInput): Promise<LoginResponse> => {
  const response = await api.post<LoginResponse>("/auth/login", data);
  return response.data;
};

// Verify OTP for 2FA
export const verifyOtp = async (tempToken: string, otp: string) => {
  const response = await api.post("/auth/verify-otp", {
    tempToken,
    otp,
  });
  return response.data;
};

// Resend OTP for 2FA
export const resendOtp = async (tempToken: string) => {
  const response = await api.post("/auth/resend-otp", {
    tempToken,
  });
  return response.data;
};

// Logout
export const logoutUser = async () => {
  const response = await api.post("/auth/logout");
  return response.data;
};

// Get Current User
export const getCurrentUser = async () => {
  const response = await api.get("/auth/me");
  return response.data;
};
