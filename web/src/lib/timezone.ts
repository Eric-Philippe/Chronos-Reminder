/**
 * Timezone-aware date/time helpers for the reminder forms.
 *
 * The browser's local Date methods (getHours, setHours, toISOString, ...)
 * always operate in the device's timezone. Reminders are scheduled against
 * the account's configured IANA timezone, which can differ from the
 * device's, so these helpers go through Intl.DateTimeFormat instead of
 * relying on the device clock.
 */

/** Formats an absolute instant as YYYY-MM-DD / HH:mm strings in the given IANA timezone. */
export function formatPartsInTimezone(
  date: Date,
  timeZone: string
): { dateStr: string; timeStr: string } {
  const formatter = new Intl.DateTimeFormat("en-CA", {
    timeZone,
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    hourCycle: "h23",
  });
  const parts = formatter.formatToParts(date).reduce((acc, part) => {
    acc[part.type] = part.value;
    return acc;
  }, {} as Record<string, string>);

  return {
    dateStr: `${parts.year}-${parts.month}-${parts.day}`,
    timeStr: `${parts.hour}:${parts.minute}`,
  };
}

/** Returns the date/time strings for "now + minutesAhead", as seen in the given timezone. */
export function getDefaultDateTime(
  timeZone: string,
  minutesAhead = 10
): { dateStr: string; timeStr: string } {
  const future = new Date(Date.now() + minutesAhead * 60 * 1000);
  return formatPartsInTimezone(future, timeZone);
}

/**
 * Checks whether a date+time (interpreted in the given timezone) is in the
 * past, by comparing wall-clock strings against "now" in that same timezone.
 * This avoids building a Date from the strings, which would otherwise be
 * interpreted in the browser's local timezone instead of the account's.
 */
export function isDateTimeInPast(
  dateStr: string,
  timeStr: string,
  timeZone: string
): boolean {
  if (!dateStr || !timeStr) return false;

  const { dateStr: nowDate, timeStr: nowTime } = formatPartsInTimezone(
    new Date(),
    timeZone
  );

  if (dateStr !== nowDate) return dateStr < nowDate;
  return timeStr <= nowTime;
}

/**
 * Parses a YYYY-MM-DD string into a Date using local calendar components.
 * Use this instead of `new Date(dateStr)` for display formatting: the native
 * parser treats date-only ISO strings as UTC midnight, which can roll over
 * to the previous/next day once re-displayed in the device's timezone.
 */
export function parseDateStrLocal(dateStr: string): Date {
  const [year, month, day] = dateStr.split("-").map(Number);
  return new Date(year, (month ?? 1) - 1, day ?? 1);
}

/**
 * Converts a wall-clock date+time (interpreted in the given IANA timezone)
 * into an absolute UTC Date instant. Used when a preview/computation needs
 * a real instant to compare against "now" or display, rather than just the
 * wall-clock strings handled by isDateTimeInPast/formatPartsInTimezone.
 */
export function zonedTimeToUtc(
  dateStr: string,
  timeStr: string,
  timeZone: string
): Date {
  const [year, month, day] = dateStr.split("-").map(Number);
  const [hour, minute] = timeStr.split(":").map(Number);

  // First guess: treat the wall-clock values as if they were UTC.
  const guess = new Date(
    Date.UTC(year, (month ?? 1) - 1, day ?? 1, hour ?? 0, minute ?? 0)
  );

  // Find what wall-clock time that UTC instant actually represents in the
  // target timezone, then correct the guess by the difference.
  const { dateStr: tzDateStr, timeStr: tzTimeStr } = formatPartsInTimezone(
    guess,
    timeZone
  );
  const [tzYear, tzMonth, tzDay] = tzDateStr.split("-").map(Number);
  const [tzHour, tzMinute] = tzTimeStr.split(":").map(Number);
  const tzAsUtc = Date.UTC(tzYear, tzMonth - 1, tzDay, tzHour, tzMinute);

  const offset = guess.getTime() - tzAsUtc;
  return new Date(guess.getTime() + offset);
}
