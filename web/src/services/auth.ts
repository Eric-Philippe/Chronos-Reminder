import { httpClient } from "./http";
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  SessionData,
  ApiResponse,
} from "./types";

// User data storage key
const USER_DATA_STORAGE_KEY = "user_data";

/**
 * Authentication Service
 * Handles login, registration, logout, and session management
 */
class AuthService {
  /**
   * Register a new user
   */
  async register(data: RegisterRequest): Promise<RegisterResponse> {
    const response = await httpClient.post<ApiResponse<RegisterResponse>>(
      "/api/auth/register",
      data
    );
    return (response.data || response) as RegisterResponse;
  }

  /**
   * Login user
   */
  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await httpClient.post<ApiResponse<LoginResponse>>(
      "/api/auth/login",
      data
    );
    const loginData = (response.data || response) as LoginResponse;

    // Store token and user data
    httpClient.setToken(loginData.token, new Date(loginData.expires_at));
    this.setUserData({
      user_id: loginData.id,
      email: loginData.email,
      username: loginData.username,
      expires_at: loginData.expires_at,
    });

    return loginData;
  }

  /**
   * Logout user
   */
  async logout(): Promise<void> {
    try {
      await httpClient.post("/api/auth/logout", {});
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      this.clearSession();
    }
  }

  /**
   * Check if user is authenticated
   */
  isAuthenticated(): boolean {
    return httpClient.isAuthenticated();
  }

  /**
   * Get current user data
   */
  getUserData(): SessionData | null {
    const data = localStorage.getItem(USER_DATA_STORAGE_KEY);
    return data ? JSON.parse(data) : null;
  }

  /**
   * Store user data
   */
  private setUserData(data: SessionData): void {
    localStorage.setItem(USER_DATA_STORAGE_KEY, JSON.stringify(data));
  }

  /**
   * Clear all session data
   */
  private clearSession(): void {
    localStorage.removeItem(USER_DATA_STORAGE_KEY);
    httpClient.clearAuth();
  }

  /**
   * Set authentication token and user data (used by OAuth flows)
   */
  setAuthentication(
    token: string,
    expiresAtStr: string,
    userData: SessionData
  ): void {
    httpClient.setToken(token, new Date(expiresAtStr));
    this.setUserData(userData);
  }
}

// Export singleton instance
export const authService = new AuthService();
