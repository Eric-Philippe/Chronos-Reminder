import { httpClient } from "./http";
import type {
  DiscordGuild,
  DiscordChannel,
  DiscordRole,
  GetUserGuildsResponse,
  GetGuildChannelsResponse,
  GetGuildRolesResponse,
} from "./types";

/**
 * Discord Service
 * Handles all Discord-related API calls for guilds, channels, and roles
 */
class DiscordService {
  /**
   * Fetch all Discord guilds for the authenticated user
   */
  async getUserGuilds(accountId: string): Promise<DiscordGuild[]> {
    try {
      const response = await httpClient.post<GetUserGuildsResponse>(
        "/api/discord/guilds",
        {
          account_id: accountId,
        }
      );

      if (response.error) {
        throw new Error(response.error);
      }

      return response.guilds || [];
    } catch (error) {
      console.error("Failed to fetch user guilds:", error);
      throw error;
    }
  }

  /**
   * Fetch all channels for a specific Discord guild
   */
  async getGuildChannels(
    accountId: string,
    guildId: string
  ): Promise<DiscordChannel[]> {
    try {
      const response = await httpClient.post<GetGuildChannelsResponse>(
        "/api/discord/guilds/channels",
        {
          account_id: accountId,
          guild_id: guildId,
        }
      );

      if (response.error) {
        throw new Error(response.error);
      }

      return response.channels || [];
    } catch (error) {
      console.error("Failed to fetch guild channels:", error);
      throw error;
    }
  }

  /**
   * Fetch all roles for a specific Discord guild
   */
  async getGuildRoles(
    accountId: string,
    guildId: string
  ): Promise<DiscordRole[]> {
    try {
      const response = await httpClient.post<GetGuildRolesResponse>(
        "/api/discord/guilds/roles",
        {
          account_id: accountId,
          guild_id: guildId,
        }
      );

      if (response.error) {
        throw new Error(response.error);
      }

      return response.roles || [];
    } catch (error) {
      console.error("Failed to fetch guild roles:", error);
      throw error;
    }
  }

  /**
   * Generate Discord bot invite URL for a specific guild
   * Uses the VITE_DISCORD_CLIENT_ID environment variable
   */
  getBotInviteUrl(guildId: string): string {
    const botClientId = import.meta.env.VITE_DISCORD_CLIENT_ID;
    if (!botClientId) {
      throw new Error(
        "Discord Client ID is not configured. Set VITE_DISCORD_CLIENT_ID in your environment variables."
      );
    }
    const permissions = "2147483648"; // Manage Channels + Send Messages + Mention Everyone
    return `https://discord.com/api/oauth2/authorize?client_id=${botClientId}&permissions=${permissions}&scope=bot&guild_id=${guildId}`;
  }
}

// Export singleton instance
export const discordService = new DiscordService();
