export interface User {
  id: string;
  email: string;
}

export interface AuthResponse {
  user: User;
}

export interface AuthInput {
  email: string;
  password: string;
}

export interface LoginResponse {
  requires2FA?: boolean;
  tempToken?: string;
  email?: string;
  message?: string;
}
