/**
 * API Response Wrapper
 */
export interface ApiResponse<T = unknown> {
  data?: T;
  error?: string;
  message?: string;
}

/**
 * Authentication Types
 */
export interface LoginResponse {
  id: string;
  email: string;
  username: string;
  token: string;
  expires_at: string;
  message: string;
}

export interface RegisterResponse {
  id: string;
  email: string;
  username: string;
  message: string;
}

export interface SessionData {
  user_id: string;
  email: string;
  username: string;
  expires_at: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
  timezone: string;
  [key: string]: unknown;
}

export interface LoginRequest {
  email: string;
  password: string;
  remember_me: boolean;
  [key: string]: unknown;
}

/**
 * Reminder Types
 */
export interface Reminder {
  id: string;
  account_id: string;
  remind_at_utc: string;
  snoozed_at_utc?: string | null;
  next_fire_utc?: string | null;
  message: string;
  created_at: string;
  recurrence_type: number;
  is_paused: boolean;
  destinations?: ReminderDestination[];
}

export interface ReminderDestination {
  id: string;
  reminder_id: string;
  type: "discord_dm" | "discord_channel" | "webhook";
  metadata: Record<string, unknown>;
}

export interface RemindersResponse {
  reminders: Reminder[];
  count: number;
}

/**
 * Account Types
 */
export interface Account {
  id: string;
  email: string;
  username: string;
  timezone: string;
  created_at: string;
  identities?: AccountIdentity[];
}

export interface AccountIdentity {
  id: string;
  account_id: string;
  provider: string;
  provider_id: string;
  created_at: string;
}

export interface AccountResponse {
  id: string;
  email: string;
  username: string;
  timezone: string;
  created_at: string;
  identities?: AccountIdentity[];
}

/**
 * Error Types
 */
export interface ReminderError {
  id: string;
  reminder_id: string;
  error_message: string;
  created_at: string;
}

export interface ReminderErrorsResponse {
  errors: ReminderError[];
  count: number;
}

/**
 * Discord Guild Types
 */
export interface DiscordGuild {
  id: string;
  name: string;
  icon: string;
  owner: boolean;
  permissions: number;
  features: string[];
}

export interface DiscordChannel {
  id: string;
  name: string;
  type: number;
  position: number;
  topic?: string | null;
}

export interface DiscordRole {
  id: string;
  name: string;
  color: number;
  position: number;
  permissions: number;
  managed: boolean;
  mentionable: boolean;
}

export interface GetUserGuildsResponse {
  guilds: DiscordGuild[];
  error?: string;
}

export interface GetGuildChannelsResponse {
  channels: DiscordChannel[];
  error?: string;
}

export interface GetGuildRolesResponse {
  roles: DiscordRole[];
  error?: string;
}
