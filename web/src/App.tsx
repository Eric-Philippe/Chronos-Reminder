import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { ThemeProvider } from "./components/theme-provider";
import { AuthProvider } from "./hooks/AuthContext";
import { LoginPage } from "./pages/LoginPage";
import { WelcomePage } from "./pages/WelcomePage";
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
      {/* Protected route: Welcome page requires authentication */}
      <Route
        path="/welcome"
        element={
          isAuthenticated ? <WelcomePage /> : <Navigate to="/login" replace />
        }
      />

      {/* Login route: Redirect to welcome if already authenticated */}
      <Route
        path="/login"
        element={
          isAuthenticated ? <Navigate to="/welcome" replace /> : <LoginPage />
        }
      />

      {/* Default route: Redirect to appropriate page based on auth status */}
      <Route
        path="/"
        element={
          <Navigate to={isAuthenticated ? "/welcome" : "/login"} replace />
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
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
