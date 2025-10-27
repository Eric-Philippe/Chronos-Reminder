import {
  createContext,
  useContext,
  useCallback,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import {
  authService,
  type LoginRequest,
  type RegisterRequest,
  type SessionData,
  type LoginResponse,
  type RegisterResponse,
} from "@/services";

interface UseAuthReturn {
  isAuthenticated: boolean;
  isLoading: boolean;
  isCheckingAuth: boolean;
  error: Error | null;
  user: SessionData | null;
  login: (
    email: string,
    password: string,
    rememberMe: boolean
  ) => Promise<void>;
  register: (
    email: string,
    username: string,
    password: string,
    timezone: string
  ) => Promise<void>;
  logout: () => Promise<void>;
  clearError: () => void;
}

const AuthContext = createContext<UseAuthReturn | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [isCheckingAuth, setIsCheckingAuth] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [user, setUser] = useState<SessionData | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(() => {
    // Initialize with actual auth status from authService
    return authService.isAuthenticated();
  });

  // Check authentication status on mount and set user data
  useEffect(() => {
    const checkAuth = () => {
      const authenticated = authService.isAuthenticated();
      setIsAuthenticated(authenticated);

      if (authenticated) {
        const userData = authService.getUserData();
        setUser(userData);
      } else {
        setUser(null);
      }

      // Mark auth check as complete
      setIsCheckingAuth(false);
    };

    checkAuth();

    // Listen for storage changes (from other tabs or from localStorage updates)
    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === "auth_token" || e.key === "user_data") {
        console.log("[AUTH_CONTEXT] Storage changed, rechecking auth...");
        checkAuth();
      }
    };

    // Also listen for custom events from the same tab
    const handleAuthUpdate = () => {
      console.log(
        "[AUTH_CONTEXT] Auth update event received, rechecking auth..."
      );
      checkAuth();
    };

    window.addEventListener("storage", handleStorageChange);
    window.addEventListener("auth-updated", handleAuthUpdate);

    return () => {
      window.removeEventListener("storage", handleStorageChange);
      window.removeEventListener("auth-updated", handleAuthUpdate);
    };
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  const login = useCallback(
    async (email: string, password: string, rememberMe: boolean) => {
      setIsLoading(true);
      setError(null);

      try {
        const loginRequest: LoginRequest = {
          email,
          password,
          remember_me: rememberMe,
        };

        const response: LoginResponse = await authService.login(loginRequest);

        const userData: SessionData = {
          user_id: response.id,
          email: response.email,
          username: response.username,
          expires_at: response.expires_at,
        };

        // Update state - these updates should trigger a re-render
        setUser(userData);
        setIsAuthenticated(true);
        setError(null);
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Login failed";
        const errorObj = new Error(errorMessage);
        setError(errorObj);
        setUser(null);
        setIsAuthenticated(false);
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    []
  );

  const register = useCallback(
    async (
      email: string,
      username: string,
      password: string,
      timezone: string
    ) => {
      setIsLoading(true);
      setError(null);

      try {
        const registerRequest: RegisterRequest = {
          email,
          username,
          password,
          timezone,
        };

        const response: RegisterResponse = await authService.register(
          registerRequest
        );

        // After registration, user needs to log in
        console.log("Registration successful:", response);
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Registration failed";
        setError(new Error(errorMessage));
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    []
  );

  const logout = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      await authService.logout();
      setUser(null);
      setIsAuthenticated(false);
      setError(null);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Logout failed";
      setError(new Error(errorMessage));
      // Still clear auth even if logout request fails
      setUser(null);
      setIsAuthenticated(false);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const value: UseAuthReturn = {
    isAuthenticated,
    isLoading,
    isCheckingAuth,
    error,
    user,
    login,
    register,
    logout,
    clearError,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth(): UseAuthReturn {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}

export type AuthProvider = typeof AuthProvider;
