import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { Toaster } from "sonner";
import { ThemeProvider } from "./components/common/theme-provider";
import { AuthProvider } from "./hooks/AuthContext";
import { LoginPage } from "./pages/LoginPage";
import { DashboardPage } from "./pages/DashboardPage";
import { CreateReminderPage } from "./pages/CreateReminderPage";
import { OAuthCallbackPage } from "./pages/OAuthCallbackPage";
import { useAuth } from "./hooks/useAuth";
import "./i18n/config";

function AppRoutes() {
  const { isAuthenticated, isCheckingAuth } = useAuth();

  // Don't render routes while checking initial auth status
  // isCheckingAuth is only true during the initial mount auth check
  if (isCheckingAuth) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-pulse">Loading...</div>
      </div>
    );
  }

  return (
    <Routes>
      {/* OAuth Callback route: Public route for Discord OAuth callback */}
      <Route path="/auth/callback/discord" element={<OAuthCallbackPage />} />

      {/* Protected route: Dashboard page requires authentication */}
      <Route
        path="/dashboard"
        element={
          isAuthenticated ? <DashboardPage /> : <Navigate to="/login" replace />
        }
      />

      {/* Protected route: Create Reminder page requires authentication */}
      <Route
        path="/reminders/create"
        element={
          isAuthenticated ? (
            <CreateReminderPage />
          ) : (
            <Navigate to="/login" replace />
          )
        }
      />

      {/* Login route: Redirect to dashboard if already authenticated */}
      <Route
        path="/login"
        element={
          isAuthenticated ? <Navigate to="/dashboard" replace /> : <LoginPage />
        }
      />

      {/* Default route: Redirect to appropriate page based on auth status */}
      <Route
        path="/"
        element={
          <Navigate to={isAuthenticated ? "/dashboard" : "/login"} replace />
        }
      />

      {/* Catch all: Redirect to appropriate page */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <AuthProvider>
        <BrowserRouter>
          <AppRoutes />
          <Toaster richColors theme="dark" position="top-right" expand />
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
