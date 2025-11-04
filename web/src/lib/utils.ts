import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Check if a date and time combination is in the past
 * @param date - Date object
 * @param time - Time string in HH:mm format
 * @returns true if the datetime is in the past, false otherwise
 */
export function isDateTimeInPast(date: Date | null, time: string): boolean {
  if (!date || !time) return false;

  // Parse the time
  const [hours, minutes] = time.split(":").map(Number);
  if (isNaN(hours) || isNaN(minutes)) return false;

  // Create a datetime object with the given date and time
  const selectedDateTime = new Date(date);
  selectedDateTime.setHours(hours, minutes, 0, 0);

  // Compare with current time
  return selectedDateTime < new Date();
}
