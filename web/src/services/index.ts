/**
 * API Services Index
 * Centralized exports for all API services
 */

export { httpClient } from "./http";
export { authService } from "./auth";
export { remindersService } from "./reminders";
export { accountService } from "./account";

// Export all types
export type {
  ApiResponse,
  LoginResponse,
  RegisterResponse,
  SessionData,
  RegisterRequest,
  LoginRequest,
  Reminder,
  ReminderDestination,
  RemindersResponse,
  Account,
  AccountIdentity,
  AccountResponse,
  ReminderError,
  ReminderErrorsResponse,
} from "./types";
