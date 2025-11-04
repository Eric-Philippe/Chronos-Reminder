/**
 * Central routing configuration for the application
 * Each route defines:
 * - path: the URL path
 * - requiresAuth: whether authentication is required
 * - name: human-readable name for the route
 * - showInNav: whether to show in navigation (optional)
 * - group: menu group name for dropdown menus
 * - submenu: whether this is a submenu item
 */

export interface Route {
  path: string;
  requiresAuth: boolean;
  name: string;
  showInNav?: boolean;
  group?: string;
  submenu?: boolean;
}

export interface MenuGroup {
  name: string;
  label: string;
  requiresAuth: boolean;
  items: Route[];
  showcase?: {
    title?: string;
    items: Array<{
      icon: string; // emoji or icon name
      label?: string;
    }>;
  };
}

export const ROUTES = {
  HOME: {
    path: "/",
    requiresAuth: false,
    name: "home",
    showInNav: true,
  } as Route,

  // Reminders group
  REMINDERS: {
    path: "/reminders",
    requiresAuth: true,
    name: "myReminders",
    showInNav: true,
    group: "reminders",
  } as Route,

  REMINDERS_CREATE: {
    path: "/reminders/create",
    requiresAuth: true,
    name: "createReminder",
    submenu: true,
    group: "reminders",
  } as Route,

  REMINDER_DETAILS: {
    path: "/reminders/:reminderId",
    requiresAuth: true,
    name: "reminderDetails",
    showInNav: false,
  } as Route,

  // Resources group
  CHANGELOG: {
    path: "/changelog",
    requiresAuth: false,
    name: "changelog",
    submenu: true,
    group: "resources",
  } as Route,

  SELFHOST: {
    path: "/selfhost",
    requiresAuth: false,
    name: "selfHost",
    submenu: true,
    group: "resources",
  } as Route,

  // Help group
  CONTACT: {
    path: "/contact",
    requiresAuth: false,
    name: "contact",
    submenu: true,
    group: "help",
  } as Route,

  HELP: {
    path: "/help",
    requiresAuth: false,
    name: "help",
    submenu: true,
    group: "help",
  } as Route,

  STATUS: {
    path: "/status",
    requiresAuth: false,
    name: "status",
    submenu: true,
    group: "help",
  } as Route,

  // Account (Settings group)
  ACCOUNT: {
    path: "/account",
    requiresAuth: true,
    name: "myAccount",
    submenu: true,
    group: "settings",
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

// Menu groups configuration
export const MENU_GROUPS: MenuGroup[] = [
  {
    name: "reminders",
    label: "myReminders",
    requiresAuth: true,
    items: [ROUTES.REMINDERS, ROUTES.REMINDERS_CREATE],
  },
  {
    name: "resources",
    label: "resources",
    requiresAuth: false,
    items: [ROUTES.CHANGELOG, ROUTES.SELFHOST],
  },
  {
    name: "help",
    label: "help",
    requiresAuth: false,
    items: [ROUTES.CONTACT, ROUTES.HELP, ROUTES.STATUS],
  },
  {
    name: "settings",
    label: "settings",
    requiresAuth: true,
    items: [ROUTES.ACCOUNT],
  },
];

export const ROUTES_ARRAY: Route[] = Object.values(ROUTES);
