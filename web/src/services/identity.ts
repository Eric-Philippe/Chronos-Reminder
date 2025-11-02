import { accountService } from "./account";
import type { Account, AccountIdentity } from "./types";

/**
 * Identity Capabilities
 * Represents what features the user has access to based on their configured identities
 */
export interface IdentityCapabilities {
  hasDiscordIdentity: boolean;
  hasAppIdentity: boolean;
  account: Account | null;
}

/**
 * Identity Service
 * Handles identity-related operations and capability checks
 */
class IdentityService {
  /**
   * Get all user identities and determine what capabilities they have
   * @returns IdentityCapabilities including whether user has Discord and/or App identity
   */
  async getIdentityCapabilities(): Promise<IdentityCapabilities> {
    try {
      const account = await accountService.getAccount();

      if (!account) {
        return {
          hasDiscordIdentity: false,
          hasAppIdentity: false,
          account: null,
        };
      }

      // Check what identities the user has configured
      const identities = account.identities || [];
      const hasDiscordIdentity = identities.some(
        (identity: AccountIdentity) => identity.provider === "discord"
      );
      const hasAppIdentity = identities.some(
        (identity: AccountIdentity) => identity.provider === "app"
      );

      return {
        hasDiscordIdentity,
        hasAppIdentity,
        account,
      };
    } catch (error) {
      console.error("Failed to get identity capabilities:", error);
      return {
        hasDiscordIdentity: false,
        hasAppIdentity: false,
        account: null,
      };
    }
  }

  /**
   * Check if user has Discord identity
   */
  async hasDiscordIdentity(): Promise<boolean> {
    const capabilities = await this.getIdentityCapabilities();
    return capabilities.hasDiscordIdentity;
  }

  /**
   * Check if user has App identity
   */
  async hasAppIdentity(): Promise<boolean> {
    const capabilities = await this.getIdentityCapabilities();
    return capabilities.hasAppIdentity;
  }
}

// Export singleton instance
export const identityService = new IdentityService();
