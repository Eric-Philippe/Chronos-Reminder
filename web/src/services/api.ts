import axios, {
  type AxiosInstance,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from "axios";

// API Response types
interface ApiResponse<T = unknown> {
  data?: T;
  error?: string;
  message?: string;
}

interface LoginResponse {
  id: string;
  email: string;
  username: string;
  token: string;
  expires_at: string;
  message: string;
}

interface RegisterResponse {
  id: string;
  email: string;
  username: string;
  message: string;
}

interface SessionData {
  user_id: string;
  email: string;
  username: string;
  expires_at: string;
}

// Auth Request types
interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  timezone: string;
}

interface LoginRequest {
  email: string;
  password: string;
  remember_me: boolean;
}

// Token storage keys
const TOKEN_STORAGE_KEY = "auth_token";
const EXPIRES_AT_STORAGE_KEY = "token_expires_at";
const USER_DATA_STORAGE_KEY = "user_data";

/**
 * API Client with Token Management
 * Handles authentication, token refresh, and session management
 */
class ApiClient {
  private axiosInstance: AxiosInstance;
  private tokenRefreshTimer: ReturnType<typeof setTimeout> | null = null;
  private isRefreshing = false;
  private refreshSubscribers: Array<() => void> = [];

  constructor(
    baseURL: string = import.meta.env.VITE_API_URL || "http://localhost:8080"
  ) {
    this.axiosInstance = axios.create({
      baseURL,
      withCredentials: true, // Enable cookies
      timeout: 10000,
      headers: {
        "Content-Type": "application/json",
      },
    });

    // Request interceptor to add token
    this.axiosInstance.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        const token = this.getToken();
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error: unknown) => Promise.reject(error)
    );

    // Response interceptor to handle token expiration and refresh
    this.axiosInstance.interceptors.response.use(
      (response: AxiosResponse) => response,
      (error: unknown) => {
        if (axios.isAxiosError(error)) {
          const originalRequest = error.config as InternalAxiosRequestConfig & {
            _retry?: boolean;
          };

          // Don't auto-redirect for auth endpoints (login/register)
          const isAuthEndpoint =
            originalRequest.url?.includes("/api/auth/login") ||
            originalRequest.url?.includes("/api/auth/register");

          if (
            error.response?.status === 401 &&
            !originalRequest._retry &&
            !isAuthEndpoint
          ) {
            if (!this.isRefreshing) {
              this.isRefreshing = true;
              originalRequest._retry = true;

              // Token expired or invalid - redirect to login
              this.clearAuth();
              window.location.href = "/login";
            } else {
              // Queue request until token is refreshed
              return new Promise((resolve) => {
                this.refreshSubscribers.push(() => {
                  resolve(this.axiosInstance(originalRequest));
                });
              });
            }
          }
        }

        return Promise.reject(error);
      }
    );

    // Restore session on initialization
    this.restoreSession();
  }

  /**
   * Register a new user
   */
  async register(data: RegisterRequest): Promise<RegisterResponse> {
    try {
      const response = await this.axiosInstance.post<
        ApiResponse<RegisterResponse>
      >("/api/auth/register", data);
      return (response.data.data || response.data) as RegisterResponse;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  /**
   * Login user
   */
  async login(data: LoginRequest): Promise<LoginResponse> {
    try {
      const response = await this.axiosInstance.post<
        ApiResponse<LoginResponse>
      >("/api/auth/login", data);
      const loginData = (response.data.data || response.data) as LoginResponse;

      // Store token and user data
      this.setToken(loginData.token, new Date(loginData.expires_at));
      this.setUserData({
        user_id: loginData.id,
        email: loginData.email,
        username: loginData.username,
        expires_at: loginData.expires_at,
      });

      // Set up token refresh
      this.setupTokenRefresh();

      return loginData;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  /**
   * Logout user
   */
  async logout(): Promise<void> {
    try {
      await this.axiosInstance.post("/api/auth/logout");
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      this.clearAuth();
    }
  }

  /**
   * Check if user is authenticated
   */
  isAuthenticated(): boolean {
    const token = this.getToken();
    const expiresAt = this.getTokenExpiresAt();

    if (!token || !expiresAt) {
      return false;
    }

    return new Date() < expiresAt;
  }

  /**
   * Get current user data
   */
  getUserData(): SessionData | null {
    const data = localStorage.getItem(USER_DATA_STORAGE_KEY);
    return data ? JSON.parse(data) : null;
  }

  /**
   * Get stored token
   */
  getToken(): string | null {
    return localStorage.getItem(TOKEN_STORAGE_KEY);
  }

  /**
   * Get token expiration time
   */
  private getTokenExpiresAt(): Date | null {
    const expiresAt = localStorage.getItem(EXPIRES_AT_STORAGE_KEY);
    return expiresAt ? new Date(expiresAt) : null;
  }

  /**
   * Store token with expiration
   */
  private setToken(token: string, expiresAt: Date): void {
    localStorage.setItem(TOKEN_STORAGE_KEY, token);
    localStorage.setItem(EXPIRES_AT_STORAGE_KEY, expiresAt.toISOString());
  }

  /**
   * Store user data
   */
  private setUserData(data: SessionData): void {
    localStorage.setItem(USER_DATA_STORAGE_KEY, JSON.stringify(data));
  }

  /**
   * Clear all auth data
   */
  private clearAuth(): void {
    localStorage.removeItem(TOKEN_STORAGE_KEY);
    localStorage.removeItem(EXPIRES_AT_STORAGE_KEY);
    localStorage.removeItem(USER_DATA_STORAGE_KEY);

    if (this.tokenRefreshTimer) {
      clearTimeout(this.tokenRefreshTimer);
      this.tokenRefreshTimer = null;
    }

    this.isRefreshing = false;
    this.refreshSubscribers = [];
  }

  /**
   * Setup automatic token refresh
   * Refreshes token 5 minutes before expiration
   */
  private setupTokenRefresh(): void {
    if (this.tokenRefreshTimer) {
      clearTimeout(this.tokenRefreshTimer);
    }

    const expiresAt = this.getTokenExpiresAt();
    if (!expiresAt) return;

    const now = new Date();
    const refreshTime = expiresAt.getTime() - now.getTime() - 5 * 60 * 1000; // 5 minutes before expiry

    if (refreshTime > 0) {
      this.tokenRefreshTimer = setTimeout(() => {
        this.refreshToken();
      }, refreshTime);
    }
  }

  /**
   * Refresh the token (placeholder - implement based on backend)
   */
  private async refreshToken(): Promise<void> {
    try {
      // This would call a refresh endpoint if your backend provides one
      // For now, we'll just notify subscribers and clear auth
      this.refreshSubscribers.forEach((callback) => callback());
      this.refreshSubscribers = [];
      this.isRefreshing = false;

      // If token refresh fails, clear auth and redirect to login
      this.clearAuth();
      window.location.href = "/login";
    } catch (error) {
      console.error("Token refresh failed:", error);
      this.clearAuth();
      window.location.href = "/login";
    }
  }

  /**
   * Restore session from localStorage if available
   */
  private restoreSession(): void {
    if (this.isAuthenticated()) {
      this.setupTokenRefresh();
    } else {
      this.clearAuth();
    }
  }

  /**
   * Handle API errors
   */
  private handleError(error: unknown): Error {
    if (axios.isAxiosError(error)) {
      const message =
        (error.response?.data as Record<string, string>)?.error ||
        (error.response?.data as Record<string, string>)?.message ||
        error.message ||
        "An error occurred";
      return new Error(message);
    }
    return error instanceof Error ? error : new Error("An error occurred");
  }

  /**
   * Get axios instance for custom requests
   */
  getAxiosInstance(): AxiosInstance {
    return this.axiosInstance;
  }
}

// Export singleton instance
export const apiClient = new ApiClient();

// Export types
export type {
  ApiResponse,
  LoginResponse,
  RegisterResponse,
  SessionData,
  RegisterRequest,
  LoginRequest,
};

export default apiClient;
