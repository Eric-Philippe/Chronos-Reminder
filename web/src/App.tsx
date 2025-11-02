import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { Toaster } from "sonner";
import { ThemeProvider } from "./components/common/theme-provider";
import { AuthProvider } from "./hooks/AuthContext";
import { VitrinePage } from "./pages/VitrinePage";
import { LoginPage } from "./pages/LoginPage";
import { HomePage } from "./pages/HomePage";
import { CreateReminderPage } from "./pages/CreateReminderPage";
import { ReminderDetailsPage } from "./pages/ReminderDetailsPage";
import { AccountPage } from "./pages/AccountPage";
import { OAuthCallbackPage } from "./pages/OAuthCallbackPage";
import { useAuth } from "./hooks/useAuth";
import { ROUTES } from "./config/routes";
import "./i18n/config";

function AppRoutes() {
  const { isAuthenticated, isCheckingAuth } = useAuth();

  const clientId = import.meta.env.VITE_DISCORD_CLIENT_ID;
  const redirectUri = import.meta.env.VITE_DISCORD_REDIRECT_URI;
  const API_URL = import.meta.env.VITE_API_URL || "https://api.chronosrmdr.com";

  console.log("API URL:", API_URL);
  console.log("DISCORD CLIENT ID:", clientId);
  console.log("DISCORD REDIRECT URI:", redirectUri);

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
      <Route
        path={ROUTES.AUTH_CALLBACK_DISCORD.path}
        element={<OAuthCallbackPage />}
      />

      {/* Vitrine route: Public route at root for all users */}
      <Route path={ROUTES.VITRINE.path} element={<VitrinePage />} />

      {/* Protected route: Dashboard page requires authentication */}
      <Route
        path={ROUTES.DASHBOARD.path}
        element={
          isAuthenticated ? (
            <HomePage />
          ) : (
            <Navigate to={ROUTES.VITRINE.path} replace />
          )
        }
      />

      {/* Protected route: Create Reminder page requires authentication */}
      <Route
        path={ROUTES.CREATE_REMINDER.path}
        element={
          isAuthenticated ? (
            <CreateReminderPage />
          ) : (
            <Navigate to={ROUTES.VITRINE.path} replace />
          )
        }
      />

      {/* Protected route: Reminder Details page requires authentication */}
      <Route
        path={ROUTES.REMINDER_DETAILS.path}
        element={
          isAuthenticated ? (
            <ReminderDetailsPage />
          ) : (
            <Navigate to={ROUTES.VITRINE.path} replace />
          )
        }
      />

      {/* Protected route: Account page requires authentication */}
      <Route
        path={ROUTES.ACCOUNT.path}
        element={
          isAuthenticated ? (
            <AccountPage />
          ) : (
            <Navigate to={ROUTES.VITRINE.path} replace />
          )
        }
      />

      {/* Login route: Redirect to dashboard if already authenticated */}
      <Route
        path={ROUTES.LOGIN.path}
        element={
          isAuthenticated ? (
            <Navigate to={ROUTES.DASHBOARD.path} replace />
          ) : (
            <LoginPage />
          )
        }
      />

      {/* Catch all: Redirect to vitrine by default */}
      <Route
        path="*"
        element={
          <Navigate
            to={isAuthenticated ? ROUTES.DASHBOARD.path : ROUTES.VITRINE.path}
            replace
          />
        }
      />
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
