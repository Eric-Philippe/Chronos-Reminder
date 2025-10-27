import { httpClient } from "./http";
import type { Account, AccountResponse, ApiResponse } from "./types";

/**
 * Account Service
 * Handles all account-related API calls
 */
class AccountService {
  /**
   * Normalize account data from API response
   * Ensures timezone is a string, not an object
   */
  private normalizeAccount(
    account: Record<string, unknown> & Partial<Account>
  ): Account {
    // Handle timezone field - it might be an object or a string
    let timezone = "UTC";
    if (account.timezone) {
      if (typeof account.timezone === "string") {
        timezone = account.timezone;
      } else if (
        typeof account.timezone === "object" &&
        account.timezone !== null
      ) {
        // If it's an object like {id, name, gmt_offset, iana_location}, use the name or iana_location
        const tzObj = account.timezone as Record<string, unknown>;
        timezone =
          (tzObj.iana_location as string) || (tzObj.name as string) || "UTC";
      }
    }

    return {
      id: String(account.id || ""),
      email: String(account.email || ""),
      username: String(account.username || ""),
      timezone,
      created_at: String(account.created_at || ""),
      identities: Array.isArray(account.identities) ? account.identities : [],
    };
  }

  /**
   * Fetch the authenticated user's account information
   */
  async getAccount(): Promise<Account | null> {
    try {
      const response = await httpClient.get<ApiResponse<AccountResponse>>(
        "/api/account"
      );
      const account = (response.data || response) as Account;
      return this.normalizeAccount(
        account as Record<string, unknown> & Partial<Account>
      );
    } catch (error) {
      console.error("Failed to fetch account:", error);
      return null;
    }
  }
}

// Export singleton instance
export const accountService = new AccountService();
