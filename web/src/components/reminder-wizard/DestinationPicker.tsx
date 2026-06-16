import React, { useState, useEffect } from "react";
import { Trash2, Plus, MessageCircle, Megaphone, Link2, Mail, Smartphone } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { useTranslation } from "react-i18next";
import { identityService } from "@/services/identity";
import { DiscordGuildSelectionModal } from "./DiscordGuildSelectionModal";

export type ReminderDestinationType =
  | "discord_dm"
  | "discord_channel"
  | "webhook"
  | "email"
  | "android_push";

export type WebhookPlatform = "generic" | "discord" | "slack";

export interface ReminderDestination {
  type: ReminderDestinationType;
  metadata: Record<string, unknown>;
}

interface DestinationPickerProps {
  destinations: ReminderDestination[];
  onDestinationsChange: (destinations: ReminderDestination[]) => void;
  recurrence?: string;
  showTitle?: boolean;
  showAddOptions?: boolean;
  compact?: boolean;
}

export function DestinationPicker({
  destinations,
  onDestinationsChange,
  recurrence,
  showTitle = true,
  showAddOptions = true,
  compact = false,
}: DestinationPickerProps) {
  const { t } = useTranslation();
  const [webhookUrl, setWebhookUrl] = useState("");
  const [webhookPlatform, setWebhookPlatform] =
    useState<WebhookPlatform>("generic");
  const [webhookUsername, setWebhookUsername] = useState("");
  const [isGuildModalOpen, setIsGuildModalOpen] = useState(false);
  const [hasDiscordIdentity, setHasDiscordIdentity] = useState(false);
  const [hasEmail, setHasEmail] = useState(false);
  const [userEmail, setUserEmail] = useState<string | null>(null);
  const [hasAndroidPush, setHasAndroidPush] = useState(false);

  // Load user identities on component mount
  useEffect(() => {
    const loadIdentities = async () => {
      try {
        const capabilities = await identityService.getIdentityCapabilities();
        setHasDiscordIdentity(capabilities.hasDiscordIdentity);
        setHasEmail(capabilities.hasEmail);
        setUserEmail(capabilities.userEmail);
        setHasAndroidPush(capabilities.hasAndroidPush);
      } catch (error) {
        console.error("Failed to load identity capabilities:", error);
        setHasDiscordIdentity(false);
        setHasEmail(false);
        setUserEmail(null);
        setHasAndroidPush(false);
      }
    };

    loadIdentities();
  }, []);

  // Count destinations by type
  const dmCount = destinations.filter((d) => d.type === "discord_dm").length;
  const channelCount = destinations.filter(
    (d) => d.type === "discord_channel"
  ).length;
  const webhookCount = destinations.filter((d) => d.type === "webhook").length;
  const emailCount = destinations.filter((d) => d.type === "email").length;
  const pushCount = destinations.filter((d) => d.type === "android_push").length;

  // Destination limits
  const MAX_DM = 1;
  const MAX_CHANNELS = 5;
  const MAX_WEBHOOKS = 5;
  const MAX_EMAILS = 1;
  const MAX_PUSH = 1;

  const isHourlyRecurrence = recurrence === "HOURLY";

  const handleAddPush = () => {
    if (!hasAndroidPush || pushCount >= MAX_PUSH) return;
    onDestinationsChange([
      ...destinations,
      { type: "android_push" as const, metadata: {} },
    ]);
  };

  const handleAddEmail = () => {
    if (!hasEmail || emailCount >= MAX_EMAILS || isHourlyRecurrence) return;
    onDestinationsChange([
      ...destinations,
      { type: "email" as const, metadata: { email: userEmail } },
    ]);
  };

  const handleAddDiscordDM = () => {
    if (dmCount >= MAX_DM) return;
    const newDest = {
      type: "discord_dm" as const,
      metadata: {},
    };
    onDestinationsChange([...destinations, newDest]);
  };

  const handleAddDiscordGuild = () => {
    if (channelCount >= MAX_CHANNELS) return;
    setIsGuildModalOpen(true);
  };

  const handleGuildModalConfirm = (
    guildId: string,
    channelId: string,
    roleId?: string
  ) => {
    const metadata: Record<string, unknown> = {
      guild_id: guildId,
      channel_id: channelId,
    };

    if (roleId) {
      metadata.mention_role_id = roleId;
    }

    const newDest = {
      type: "discord_channel" as const,
      metadata,
    };
    onDestinationsChange([...destinations, newDest]);
  };

  const handleAddWebhook = () => {
    if (!webhookUrl.trim()) return;
    if (webhookCount >= MAX_WEBHOOKS) return;

    const metadata: Record<string, unknown> = {
      url: webhookUrl,
      platform: webhookPlatform,
    };

    // Add optional fields based on platform
    if (webhookUsername.trim()) {
      metadata.username = webhookUsername;
    }

    const newDest = {
      type: "webhook" as const,
      metadata,
    };
    onDestinationsChange([...destinations, newDest]);
    setWebhookUrl("");
    setWebhookPlatform("generic");
    setWebhookUsername("");
  };

  const handleRemoveDestination = (index: number) => {
    onDestinationsChange(destinations.filter((_, i) => i !== index));
  };

  const handleUpdateWebhookMetadata = (
    index: number,
    url: string,
    platform?: WebhookPlatform,
    username?: string
  ) => {
    const updated = [...destinations];
    const currentMetadata = updated[index].metadata;
    updated[index].metadata = {
      ...currentMetadata,
      url,
      ...(platform && { platform }),
      ...(username !== undefined && { username }),
    };
    onDestinationsChange(updated);
  };

  return (
    <div className="space-y-4">
      {showTitle && (
        <div>
          <h3 className="text-sm font-semibold text-foreground mb-4">
            {t("reminderCreation.destinations.title")}
          </h3>
        </div>
      )}

      {/* Selected Destinations */}
      {destinations.length > 0 && (
        <div className={compact ? "space-y-2" : "space-y-3"}>
          {compact && (
            <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">
              {t("reminderCreation.destinations.addedDestinations")} (
              {destinations.length})
            </p>
          )}
          {!compact && (
            <h4 className="text-sm font-semibold text-foreground">
              {t("reminderCreation.destinations.addedDestinations")} (
              {destinations.length})
            </h4>
          )}
          {destinations.map((dest, idx) => (
            <div
              key={idx}
              className={`p-3 rounded-lg border border-border bg-secondary/20 space-y-3 ${
                compact ? "p-2" : ""
              }`}
            >
              {dest.type === "discord_dm" && (
                <div className="flex items-center justify-between gap-2">
                  <div className="flex items-center gap-2 flex-1 min-w-0">
                    <MessageCircle className="w-4 h-4 text-accent flex-shrink-0" />
                    <div className="min-w-0 flex-1">
                      <p
                        className={`font-semibold text-foreground truncate ${
                          compact ? "text-xs" : "text-sm"
                        }`}
                      >
                        {t("reminderCreation.destinations.discordDM")}
                      </p>
                      {!compact && (
                        <p className="text-xs text-muted-foreground mt-1">
                          {t("reminderCreation.destinations.discordDMDesc")}
                        </p>
                      )}
                    </div>
                  </div>
                  <Button
                    onClick={() => handleRemoveDestination(idx)}
                    variant="outline"
                    size="sm"
                    className="border-red-500/50 text-red-600 dark:text-red-400 hover:bg-red-500/10 flex-shrink-0"
                  >
                    <Trash2 className="w-3 h-3" />
                  </Button>
                </div>
              )}

              {dest.type === "discord_channel" && (
                <div className="flex items-center justify-between gap-2">
                  <div className="flex items-center gap-2 flex-1 min-w-0">
                    <Megaphone className="w-4 h-4 text-accent flex-shrink-0" />
                    <div className="min-w-0 flex-1">
                      <p
                        className={`font-semibold text-foreground truncate ${
                          compact ? "text-xs" : "text-sm"
                        }`}
                      >
                        {t("reminderCreation.destinations.discordGuild")}
                      </p>
                      {!compact && (
                        <div className="text-xs text-muted-foreground space-y-0.5">
                          <p>
                            Channel:{" "}
                            <span className="text-foreground">
                              #{dest.metadata.channel_id as string}
                            </span>
                          </p>
                          {(dest.metadata.mention_role_id as string) && (
                            <p>
                              Role:{" "}
                              <span className="text-foreground">
                                @{dest.metadata.mention_role_id as string}
                              </span>
                            </p>
                          )}
                        </div>
                      )}
                    </div>
                  </div>
                  <Button
                    onClick={() => handleRemoveDestination(idx)}
                    variant="outline"
                    size="sm"
                    className="border-red-500/50 text-red-600 dark:text-red-400 hover:bg-red-500/10 flex-shrink-0"
                  >
                    <Trash2 className="w-3 h-3" />
                  </Button>
                </div>
              )}

              {dest.type === "email" && (
                <div className="flex items-center justify-between gap-2">
                  <div className="flex items-center gap-2 flex-1 min-w-0">
                    <Mail className="w-4 h-4 text-accent flex-shrink-0" />
                    <div className="min-w-0 flex-1">
                      <p
                        className={`font-semibold text-foreground truncate ${
                          compact ? "text-xs" : "text-sm"
                        }`}
                      >
                        {t("reminderCreation.destinations.email")}
                      </p>
                      {!compact && (
                        <p className="text-xs text-muted-foreground mt-1 truncate">
                          {dest.metadata.email as string}
                        </p>
                      )}
                    </div>
                  </div>
                  <Button
                    onClick={() => handleRemoveDestination(idx)}
                    variant="outline"
                    size="sm"
                    className="border-red-500/50 text-red-600 dark:text-red-400 hover:bg-red-500/10 flex-shrink-0"
                  >
                    <Trash2 className="w-3 h-3" />
                  </Button>
                </div>
              )}

              {dest.type === "android_push" && (
                <div className="flex items-center justify-between gap-2">
                  <div className="flex items-center gap-2 flex-1 min-w-0">
                    <Smartphone className="w-4 h-4 text-accent flex-shrink-0" />
                    <div className="min-w-0 flex-1">
                      <p
                        className={`font-semibold text-foreground truncate ${
                          compact ? "text-xs" : "text-sm"
                        }`}
                      >
                        {t("reminderCreation.destinations.androidPush")}
                      </p>
                      {!compact && (
                        <p className="text-xs text-muted-foreground mt-1">
                          {t("reminderCreation.destinations.androidPushDesc")}
                        </p>
                      )}
                    </div>
                  </div>
                  <Button
                    onClick={() => handleRemoveDestination(idx)}
                    variant="outline"
                    size="sm"
                    className="border-red-500/50 text-red-600 dark:text-red-400 hover:bg-red-500/10 flex-shrink-0"
                  >
                    <Trash2 className="w-3 h-3" />
                  </Button>
                </div>
              )}

              {dest.type === "webhook" && (
                <div className="space-y-2">
                  <div className="flex items-center justify-between gap-2">
                    <div className="flex items-center gap-2 flex-1 min-w-0">
                      <Link2 className="w-4 h-4 text-accent flex-shrink-0" />
                      <div className="min-w-0 flex-1">
                        <p
                          className={`font-medium text-foreground truncate ${
                            compact ? "text-xs" : "text-sm"
                          }`}
                        >
                          {t("reminderCreation.destinations.webhook")}
                        </p>
                        {(dest.metadata.platform as string) &&
                          dest.metadata.platform !== "generic" && (
                            <p className="text-xs text-muted-foreground capitalize">
                              {dest.metadata.platform as string}
                            </p>
                          )}
                      </div>
                      <span className="text-xs text-muted-foreground flex-shrink-0">
                        {idx + 1}/{webhookCount}
                      </span>
                    </div>
                    <Button
                      onClick={() => handleRemoveDestination(idx)}
                      variant="outline"
                      size="sm"
                      className="border-red-500/50 text-red-600 dark:text-red-400 hover:bg-red-500/10 flex-shrink-0"
                    >
                      <Trash2 className="w-3 h-3" />
                    </Button>
                  </div>
                  <div className="space-y-2">
                    <select
                      value={(dest.metadata.platform as string) || "generic"}
                      onChange={(e) =>
                        handleUpdateWebhookMetadata(
                          idx,
                          (dest.metadata.url as string) || "",
                          e.target.value as WebhookPlatform,
                          (dest.metadata.username as string) || ""
                        )
                      }
                      className={`w-full px-3 py-2 rounded border border-border bg-background text-foreground ${
                        compact ? "text-xs" : "text-sm"
                      }`}
                    >
                      <option value="generic">
                        {t(
                          "reminderCreation.destinations.webhookPlatforms.generic"
                        )}
                      </option>
                      <option value="discord">
                        {t(
                          "reminderCreation.destinations.webhookPlatforms.discord"
                        )}
                      </option>
                      <option value="slack">
                        {t(
                          "reminderCreation.destinations.webhookPlatforms.slack"
                        )}
                      </option>
                    </select>
                    <input
                      type="text"
                      placeholder={t(
                        "reminderCreation.destinations.webhookPlaceholder"
                      )}
                      value={(dest.metadata.url as string) || ""}
                      onChange={(e) =>
                        handleUpdateWebhookMetadata(
                          idx,
                          e.target.value,
                          (dest.metadata.platform as WebhookPlatform) ||
                            "generic",
                          (dest.metadata.username as string) || ""
                        )
                      }
                      className={`w-full px-3 py-2 rounded border border-border bg-background text-foreground placeholder-muted-foreground ${
                        compact ? "text-xs" : "text-sm"
                      }`}
                    />
                    {((dest.metadata.platform as string) === "discord" ||
                      (dest.metadata.platform as string) === "slack") && (
                      <input
                        type="text"
                        placeholder={t(
                          "reminderCreation.destinations.webhookUsername"
                        )}
                        value={(dest.metadata.username as string) || ""}
                        onChange={(e) =>
                          handleUpdateWebhookMetadata(
                            idx,
                            (dest.metadata.url as string) || "",
                            (dest.metadata.platform as WebhookPlatform) ||
                              "generic",
                            e.target.value
                          )
                        }
                        className={`w-full px-3 py-2 rounded border border-border bg-background text-foreground placeholder-muted-foreground ${
                          compact ? "text-xs" : "text-sm"
                        }`}
                      />
                    )}
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Add Destination Options */}
      {showAddOptions && (
        <div className={`space-y-2 ${compact ? "mt-2" : "mt-4"}`}>
          {!compact && (
            <h4 className="text-sm font-semibold text-foreground">
              {t("reminderCreation.destinations.addDestination")}
            </h4>
          )}

              {/* Destination options sorted: available first, unavailable last */}
          {((): { available: boolean; key: string; element: React.ReactNode }[] => {
            const dmAvailable = hasDiscordIdentity && dmCount < MAX_DM;
            const channelAvailable = hasDiscordIdentity && channelCount < MAX_CHANNELS;
            const webhookAvailable = webhookCount < MAX_WEBHOOKS;
            const emailAvailable = hasEmail && emailCount < MAX_EMAILS && !isHourlyRecurrence;
            const pushAvailable = hasAndroidPush && pushCount < MAX_PUSH;

            const options = [
              {
                key: "discord_dm",
                available: dmAvailable,
                element: (
                  <Card
                    key="discord_dm"
                    className={`border-border bg-secondary/20 transition-colors ${
                      dmAvailable ? "hover:border-accent/50 cursor-pointer" : "opacity-60 cursor-not-allowed"
                    }`}
                    onClick={dmAvailable ? handleAddDiscordDM : undefined}
                  >
                    <div className="p-3 flex items-center justify-between">
                      <div className="flex items-center gap-2 flex-1 min-w-0">
                        <div className="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center flex-shrink-0">
                          <MessageCircle className="w-4 h-4 text-accent" />
                        </div>
                        <div className="min-w-0 flex-1">
                          <p className="text-sm font-semibold text-foreground">
                            {t("reminderCreation.destinations.discordDM")}
                          </p>
                          {!compact && (
                            <p className="text-xs text-muted-foreground">
                              {!hasDiscordIdentity
                                ? t("reminderCreation.destinations.connectDiscordFirst")
                                : dmCount >= MAX_DM
                                ? "Limit reached"
                                : t("reminderCreation.destinations.sendAsDM")}
                            </p>
                          )}
                        </div>
                        <span className="text-xs text-muted-foreground flex-shrink-0 mr-2">
                          {dmCount}/{MAX_DM}
                        </span>
                      </div>
                      <Button
                        onClick={(e) => { e.stopPropagation(); handleAddDiscordDM(); }}
                        variant="outline"
                        size="sm"
                        disabled={!dmAvailable}
                        className={`border-accent/50 flex-shrink-0 ${dmAvailable ? "text-accent hover:bg-accent/10" : "text-muted-foreground opacity-50 cursor-not-allowed"}`}
                      >
                        <Plus className="w-4 h-4" />
                      </Button>
                    </div>
                  </Card>
                ),
              },
              {
                key: "discord_channel",
                available: channelAvailable,
                element: (
                  <Card
                    key="discord_channel"
                    className={`border-border bg-secondary/20 transition-colors ${
                      channelAvailable ? "hover:border-accent/50 cursor-pointer" : "opacity-60 cursor-not-allowed"
                    }`}
                  >
                    <div
                      onClick={channelAvailable ? handleAddDiscordGuild : undefined}
                      className="p-3 flex items-center justify-between"
                    >
                      <div className="flex items-center gap-2 flex-1 min-w-0">
                        <div className="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center flex-shrink-0">
                          <Megaphone className="w-4 h-4 text-accent" />
                        </div>
                        <div className="min-w-0 flex-1">
                          <p className="text-sm font-semibold text-foreground">
                            {t("reminderCreation.destinations.discordGuild")}
                          </p>
                          {!compact && (
                            <p className="text-xs text-muted-foreground">
                              {!hasDiscordIdentity
                                ? t("reminderCreation.destinations.connectDiscordFirst")
                                : channelCount >= MAX_CHANNELS
                                ? "Limit reached"
                                : t("reminderCreation.destinations.discordGuildDesc")}
                            </p>
                          )}
                        </div>
                        <span className="text-xs text-muted-foreground flex-shrink-0 mr-2">
                          {channelCount}/{MAX_CHANNELS}
                        </span>
                      </div>
                      <Button
                        onClick={(e) => { e.stopPropagation(); handleAddDiscordGuild(); }}
                        variant="outline"
                        size="sm"
                        disabled={!channelAvailable}
                        className={`border-accent/50 flex-shrink-0 ${channelAvailable ? "text-accent hover:bg-accent/10" : "text-muted-foreground opacity-50 cursor-not-allowed"}`}
                      >
                        <Plus className="w-4 h-4" />
                      </Button>
                    </div>
                  </Card>
                ),
              },
              {
                key: "webhook",
                available: webhookAvailable,
                element: (
                  <Card key="webhook" className={`border-border bg-secondary/20 ${!webhookAvailable ? "opacity-60" : ""}`}>
                    <div className="p-3 space-y-3">
                      <div className="flex items-center gap-2">
                        <Link2 className="w-4 h-4 text-accent flex-shrink-0" />
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium text-foreground">
                            {t("reminderCreation.destinations.webhook")}
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {t("reminderCreation.destinations.webhookDesc")}
                          </p>
                        </div>
                        <span className="text-xs text-muted-foreground flex-shrink-0">
                          {webhookCount}/{MAX_WEBHOOKS}
                        </span>
                      </div>
                      <div className="space-y-2">
                        <select
                          value={webhookPlatform}
                          onChange={(e) => setWebhookPlatform(e.target.value as WebhookPlatform)}
                          disabled={!webhookAvailable}
                          className="w-full px-3 py-2 rounded border border-border bg-background text-foreground text-sm disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                          <option value="generic">{t("reminderCreation.destinations.webhookPlatforms.generic")}</option>
                          <option value="discord">{t("reminderCreation.destinations.webhookPlatforms.discord")}</option>
                          <option value="slack">{t("reminderCreation.destinations.webhookPlatforms.slack")}</option>
                        </select>
                        <input
                          type="text"
                          placeholder={t("reminderCreation.destinations.webhookPlaceholder")}
                          value={webhookUrl}
                          onChange={(e) => setWebhookUrl(e.target.value)}
                          disabled={!webhookAvailable}
                          className="w-full px-3 py-2 rounded border border-border bg-background text-foreground text-sm placeholder-muted-foreground disabled:opacity-50 disabled:cursor-not-allowed"
                        />
                        {(webhookPlatform === "discord" || webhookPlatform === "slack") && (
                          <input
                            type="text"
                            placeholder={t("reminderCreation.destinations.webhookUsername")}
                            value={webhookUsername}
                            onChange={(e) => setWebhookUsername(e.target.value)}
                            disabled={!webhookAvailable}
                            className="w-full px-3 py-2 rounded border border-border bg-background text-foreground text-sm placeholder-muted-foreground disabled:opacity-50 disabled:cursor-not-allowed"
                          />
                        )}
                        <Button
                          onClick={handleAddWebhook}
                          variant="outline"
                          size="sm"
                          disabled={!webhookUrl.trim() || !webhookAvailable}
                          className="w-full border-accent/50 text-accent hover:bg-accent/10 disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                          <Plus className="w-4 h-4 mr-2" />
                          {t("reminderCreation.destinations.addWebhook")}
                        </Button>
                      </div>
                    </div>
                  </Card>
                ),
              },
              {
                key: "email",
                available: emailAvailable,
                element: (
                  <Card
                    key="email"
                    className={`border-border bg-secondary/20 transition-colors ${
                      emailAvailable ? "hover:border-accent/50 cursor-pointer" : "opacity-60 cursor-not-allowed"
                    }`}
                    onClick={emailAvailable ? handleAddEmail : undefined}
                  >
                    <div className="p-3 flex items-center justify-between">
                      <div className="flex items-center gap-2 flex-1 min-w-0">
                        <div className="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center flex-shrink-0">
                          <Mail className="w-4 h-4 text-accent" />
                        </div>
                        <div className="min-w-0 flex-1">
                          <p className="text-sm font-semibold text-foreground">
                            {t("reminderCreation.destinations.email")}
                          </p>
                          {!compact && (
                            <p className="text-xs text-muted-foreground">
                              {!hasEmail
                                ? t("reminderCreation.destinations.noEmailConfigured")
                                : isHourlyRecurrence
                                ? t("reminderCreation.destinations.emailHourlyDisabled")
                                : emailCount >= MAX_EMAILS
                                ? t("reminderCreation.destinations.limitReached")
                                : userEmail ?? ""}
                            </p>
                          )}
                        </div>
                        <span className="text-xs text-muted-foreground flex-shrink-0 mr-2">
                          {emailCount}/{MAX_EMAILS}
                        </span>
                      </div>
                      <Button
                        onClick={(e) => { e.stopPropagation(); handleAddEmail(); }}
                        variant="outline"
                        size="sm"
                        disabled={!emailAvailable}
                        className={`border-accent/50 flex-shrink-0 ${emailAvailable ? "text-accent hover:bg-accent/10" : "text-muted-foreground opacity-50 cursor-not-allowed"}`}
                      >
                        <Plus className="w-4 h-4" />
                      </Button>
                    </div>
                  </Card>
                ),
              },
              {
                key: "android_push",
                available: pushAvailable,
                element: (
                  <Card
                    key="android_push"
                    className={`border-border bg-secondary/20 transition-colors ${
                      pushAvailable ? "hover:border-accent/50 cursor-pointer" : "opacity-60 cursor-not-allowed"
                    }`}
                    onClick={pushAvailable ? handleAddPush : undefined}
                  >
                    <div className="p-3 flex items-center justify-between">
                      <div className="flex items-center gap-2 flex-1 min-w-0">
                        <div className="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center flex-shrink-0">
                          <Smartphone className="w-4 h-4 text-accent" />
                        </div>
                        <div className="min-w-0 flex-1">
                          <p className="text-sm font-semibold text-foreground">
                            {t("reminderCreation.destinations.androidPush")}
                          </p>
                          {!compact && (
                            <p className="text-xs text-muted-foreground">
                              {!hasAndroidPush
                                ? t("reminderCreation.destinations.androidPushNoApp")
                                : pushCount >= MAX_PUSH
                                ? t("reminderCreation.destinations.limitReached")
                                : t("reminderCreation.destinations.androidPushDesc")}
                            </p>
                          )}
                        </div>
                        <span className="text-xs text-muted-foreground flex-shrink-0 mr-2">
                          {pushCount}/{MAX_PUSH}
                        </span>
                      </div>
                      <Button
                        onClick={(e) => { e.stopPropagation(); handleAddPush(); }}
                        variant="outline"
                        size="sm"
                        disabled={!pushAvailable}
                        className={`border-accent/50 flex-shrink-0 ${pushAvailable ? "text-accent hover:bg-accent/10" : "text-muted-foreground opacity-50 cursor-not-allowed"}`}
                      >
                        <Plus className="w-4 h-4" />
                      </Button>
                    </div>
                  </Card>
                ),
              },
            ];

            return [...options.filter((o) => o.available), ...options.filter((o) => !o.available)];
          })().map((o) => o.element)}
        </div>
      )}

      {/* Discord Guild Selection Modal */}
      <DiscordGuildSelectionModal
        open={isGuildModalOpen}
        onOpenChange={setIsGuildModalOpen}
        onConfirm={handleGuildModalConfirm}
      />
    </div>
  );
}
