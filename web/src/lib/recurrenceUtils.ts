/**
 * Recurrence type constants - matching the backend implementation
 *
 * The backend now provides recurrence_type and is_paused as separate fields
 * in API responses, so the frontend doesn't need to handle bit manipulation.
 * These constants represent only the recurrence type values (0-7).
 */

// Recurrence type constants
export const RecurrenceOnce = 0;
export const RecurrenceYearly = 1;
export const RecurrenceMonthly = 2;
export const RecurrenceWeekly = 3;
export const RecurrenceDaily = 4;
export const RecurrenceHourly = 5;
export const RecurrenceWorkdays = 6;
export const RecurrenceWeekend = 7;

/**
 * Gets the i18n translation key for a recurrence type label
 * Use with useTranslation() to get the translated label
 */
export function getRecurrenceTypeI18nKey(recurrenceType: number): string {
  const keys: Record<number, string> = {
    [RecurrenceOnce]: "recurrence.once",
    [RecurrenceYearly]: "recurrence.yearly",
    [RecurrenceMonthly]: "recurrence.monthly",
    [RecurrenceWeekly]: "recurrence.weekly",
    [RecurrenceDaily]: "recurrence.daily",
    [RecurrenceHourly]: "recurrence.hourly",
    [RecurrenceWorkdays]: "recurrence.workdays",
    [RecurrenceWeekend]: "recurrence.weekend",
  };
  return keys[recurrenceType] || "recurrence.unknown";
}

/**
 * Gets the i18n translation key for a recurrence type from string
 * Maps uppercase recurrence names to i18n keys
 */
export function getRecurrenceTypeI18nKeyFromString(
  recurrenceType: string
): string {
  const keys: Record<string, string> = {
    ONCE: "recurrence.once",
    YEARLY: "recurrence.yearly",
    MONTHLY: "recurrence.monthly",
    WEEKLY: "recurrence.weekly",
    DAILY: "recurrence.daily",
    HOURLY: "recurrence.hourly",
    WORKDAYS: "recurrence.workdays",
    WEEKEND: "recurrence.weekend",
  };
  return keys[recurrenceType.toUpperCase()] || "recurrence.unknown";
}

/**
 * Gets the label for a recurrence type (fallback, non-translated)
 * Prefer using getRecurrenceTypeI18nKey with useTranslation() for translated labels
 */
export function getRecurrenceTypeLabel(recurrenceType: number): string {
  const labels: Record<number, string> = {
    [RecurrenceOnce]: "Once",
    [RecurrenceYearly]: "Yearly",
    [RecurrenceMonthly]: "Monthly",
    [RecurrenceWeekly]: "Weekly",
    [RecurrenceDaily]: "Daily",
    [RecurrenceHourly]: "Hourly",
    [RecurrenceWorkdays]: "Workdays",
    [RecurrenceWeekend]: "Weekend",
  };
  return labels[recurrenceType] || "Unknown";
}

/**
 * Gets the name for a recurrence type (uppercase)
 */
export function getRecurrenceTypeName(recurrenceType: number): string {
  const names: Record<number, string> = {
    [RecurrenceOnce]: "ONCE",
    [RecurrenceYearly]: "YEARLY",
    [RecurrenceMonthly]: "MONTHLY",
    [RecurrenceWeekly]: "WEEKLY",
    [RecurrenceDaily]: "DAILY",
    [RecurrenceHourly]: "HOURLY",
    [RecurrenceWorkdays]: "WORKDAYS",
    [RecurrenceWeekend]: "WEEKEND",
  };
  return names[recurrenceType] || "UNKNOWN";
}

/**
 * Recurrence option for UI display
 */
export interface RecurrenceOption {
  value: number;
  label: string;
}
