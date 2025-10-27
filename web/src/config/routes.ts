/**
 * Central routing configuration for the application
 * Each route defines:
 * - path: the URL path
 * - requiresAuth: whether authentication is required
 * - name: human-readable name for the route
 * - showInNav: whether to show in navigation (optional)
 * - navGroup: 'auth' for authenticated-only nav, 'public' for public nav (optional)
 */

export interface Route {
  path: string;
  requiresAuth: boolean;
  name: string;
  showInNav?: boolean;
  navGroup?: "auth" | "public";
}

export const ROUTES = {
  // Public routes
  VITRINE: {
    path: "/",
    requiresAuth: false,
    name: "home",
    showInNav: false,
  } as Route,

  LOGIN: {
    path: "/login",
    requiresAuth: false,
    name: "signIn",
    showInNav: false,
  } as Route,

  AUTH_CALLBACK_DISCORD: {
    path: "/auth/callback/discord",
    requiresAuth: false,
    name: "Discord OAuth Callback",
    showInNav: false,
  } as Route,

  // Protected routes
  DASHBOARD: {
    path: "/home",
    requiresAuth: true,
    name: "myReminders",
    showInNav: true,
    navGroup: "auth",
  } as Route,

  CREATE_REMINDER: {
    path: "/reminders/create",
    requiresAuth: true,
    name: "createReminder",
    showInNav: false,
  } as Route,
} as const;

// Helper function to get the appropriate redirect based on auth status
export const getDefaultRoute = (isAuthenticated: boolean): string => {
  return isAuthenticated ? ROUTES.DASHBOARD.path : ROUTES.LOGIN.path;
};

// Helper function to check if a route requires authentication
export const requiresAuthentication = (path: string): boolean => {
  const route = Object.values(ROUTES).find((r) => r.path === path);
  return route?.requiresAuth ?? false;
};

// Helper function to get navigation routes for a specific group
export const getNavRoutes = (navGroup: "auth" | "public"): Route[] => {
  return Object.values(ROUTES).filter(
    (route) => route.showInNav && route.navGroup === navGroup
  );
};
