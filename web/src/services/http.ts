import axios, {
  type AxiosInstance,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from "axios";

// Token storage keys
const TOKEN_STORAGE_KEY = "auth_token";
const EXPIRES_AT_STORAGE_KEY = "token_expires_at";

/**
 * HTTP Client with Token Management
 * Handles authentication headers, token refresh, and error handling
 */
export class HttpClient {
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
   * Make a GET request
   */
  async get<T>(url: string, config?: Record<string, unknown>): Promise<T> {
    try {
      const response = await this.axiosInstance.get<T>(url, config);
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  /**
   * Make a POST request
   */
  async post<T>(
    url: string,
    data?: Record<string, unknown>,
    config?: Record<string, unknown>
  ): Promise<T> {
    try {
      const response = await this.axiosInstance.post<T>(url, data, config);
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  /**
   * Make a PUT request
   */
  async put<T>(
    url: string,
    data?: Record<string, unknown>,
    config?: Record<string, unknown>
  ): Promise<T> {
    try {
      const response = await this.axiosInstance.put<T>(url, data, config);
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  /**
   * Make a DELETE request
   */
  async delete<T>(url: string, config?: Record<string, unknown>): Promise<T> {
    try {
      const response = await this.axiosInstance.delete<T>(url, config);
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
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
  getTokenExpiresAt(): Date | null {
    const expiresAt = localStorage.getItem(EXPIRES_AT_STORAGE_KEY);
    return expiresAt ? new Date(expiresAt) : null;
  }

  /**
   * Store token with expiration
   */
  setToken(token: string, expiresAt: Date): void {
    localStorage.setItem(TOKEN_STORAGE_KEY, token);
    localStorage.setItem(EXPIRES_AT_STORAGE_KEY, expiresAt.toISOString());
    this.setupTokenRefresh();
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
   * Clear all auth data
   */
  clearAuth(): void {
    localStorage.removeItem(TOKEN_STORAGE_KEY);
    localStorage.removeItem(EXPIRES_AT_STORAGE_KEY);

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
    if (!expiresAt) {
      return;
    }

    const now = new Date();
    const timeUntilExpiry = expiresAt.getTime() - now.getTime();
    const refreshTime = timeUntilExpiry - 5 * 60 * 1000; // 5 minutes before expiry

    // Maximum safe timeout in JavaScript is ~24.8 days (2^31 - 1 milliseconds)
    const MAX_SAFE_TIMEOUT = 2147483647;

    if (refreshTime > 0 && refreshTime <= MAX_SAFE_TIMEOUT) {
      this.tokenRefreshTimer = setTimeout(() => {
        this.refreshToken();
      }, refreshTime);
    }
  }

  /**
   * Refresh the token
   */
  private async refreshToken(): Promise<void> {
    try {
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
    const authenticated = this.isAuthenticated();
    if (authenticated) {
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
export const httpClient = new HttpClient();
