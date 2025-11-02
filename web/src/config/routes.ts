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
}

export const ROUTES = {
  VITRINE: {
    path: "/",
    requiresAuth: false,
    name: "home",
    showInNav: true,
  } as Route,

  DASHBOARD: {
    path: "/home",
    requiresAuth: true,
    name: "myReminders",
    showInNav: true,
  } as Route,

  CREATE_REMINDER: {
    path: "/reminders/create",
    requiresAuth: true,
    name: "createReminder",
    showInNav: true,
  } as Route,

  REMINDER_DETAILS: {
    path: "/reminders/:reminderId",
    requiresAuth: true,
    name: "reminderDetails",
    showInNav: false,
  } as Route,

  INSTALLATION: {
    path: "/installation",
    requiresAuth: false,
    name: "installation",
    showInNav: true,
  } as Route,

  CONTACT: {
    path: "/contact",
    requiresAuth: false,
    name: "contact",
    showInNav: true,
  } as Route,

  ACCOUNT: {
    path: "/account",
    requiresAuth: true,
    name: "myAccount",
    showInNav: true,
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
} as const;

export const ROUTES_ARRAY: Route[] = Object.values(ROUTES);
