import { httpClient } from "./http";
import type { Timezone } from "./types";

/**
 * Timezone Service
 * Handles all timezone-related API calls
 */
class TimezoneService {
  /**
   * Fetch all available timezones
   */
  async getAvailableTimezones(): Promise<Timezone[]> {
    try {
      const response = await httpClient.get<Timezone[]>("/api/timezones");
      return response || [];
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error("Failed to fetch timezones");
    }
  }
}

// Export singleton instance
export const timezoneService = new TimezoneService();
