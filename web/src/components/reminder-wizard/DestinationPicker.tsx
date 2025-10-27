import { useState } from "react";
import { Trash2, Plus, MessageCircle, Megaphone, Link2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { useTranslation } from "react-i18next";

export type ReminderDestinationType =
  | "discord_dm"
  | "discord_channel"
  | "webhook";

export interface ReminderDestination {
  type: ReminderDestinationType;
  metadata: Record<string, unknown>;
}

interface DestinationPickerProps {
  destinations: ReminderDestination[];
  onDestinationsChange: (destinations: ReminderDestination[]) => void;
  showTitle?: boolean;
  showAddOptions?: boolean;
  compact?: boolean;
}

export function DestinationPicker({
  destinations,
  onDestinationsChange,
  showTitle = true,
  showAddOptions = true,
  compact = false,
}: DestinationPickerProps) {
  const { t } = useTranslation();
  const [webhookUrl, setWebhookUrl] = useState("");

  const handleAddDiscordDM = () => {
    const newDest = {
      type: "discord_dm" as const,
      metadata: {},
    };
    onDestinationsChange([...destinations, newDest]);
  };

  const handleAddDiscordGuild = () => {
    const newDest = {
      type: "discord_channel" as const,
      metadata: {
        guild_id: "",
        channel_id: "",
      },
    };
    onDestinationsChange([...destinations, newDest]);
  };

  const handleAddWebhook = () => {
    if (!webhookUrl.trim()) return;
    const newDest = {
      type: "webhook" as const,
      metadata: {
        url: webhookUrl,
      },
    };
    onDestinationsChange([...destinations, newDest]);
    setWebhookUrl("");
  };

  const handleRemoveDestination = (index: number) => {
    onDestinationsChange(destinations.filter((_, i) => i !== index));
  };

  const handleUpdateWebhookMetadata = (index: number, url: string) => {
    const updated = [...destinations];
    updated[index].metadata = { url };
    onDestinationsChange(updated);
  };

  const handleUpdateGuildMetadata = (
    index: number,
    guildId: string,
    channelId: string
  ) => {
    const updated = [...destinations];
    updated[index].metadata = { guild_id: guildId, channel_id: channelId };
    onDestinationsChange(updated);
  };

  const hasDiscordIdentity = true; // TODO: Get from context/props

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
                <div className="space-y-2">
                  <div className="flex items-center justify-between gap-2">
                    <div className="flex items-center gap-2 flex-1 min-w-0">
                      <Megaphone className="w-4 h-4 text-accent flex-shrink-0" />
                      <p
                        className={`font-semibold text-foreground truncate ${
                          compact ? "text-xs" : "text-sm"
                        }`}
                      >
                        {t("reminderCreation.destinations.discordGuild")}
                      </p>
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
                  <div
                    className={`grid gap-2 ${
                      compact ? "grid-cols-1" : "grid-cols-2"
                    }`}
                  >
                    <input
                      type="text"
                      placeholder={t("reminderCreation.destinations.guildId")}
                      value={(dest.metadata.guild_id as string) || ""}
                      onChange={(e) =>
                        handleUpdateGuildMetadata(
                          idx,
                          e.target.value,
                          (dest.metadata.channel_id as string) || ""
                        )
                      }
                      className={`px-3 py-2 rounded border border-border bg-background text-foreground placeholder-muted-foreground ${
                        compact ? "text-xs" : "text-sm"
                      }`}
                    />
                    <input
                      type="text"
                      placeholder={t("reminderCreation.destinations.channelId")}
                      value={(dest.metadata.channel_id as string) || ""}
                      onChange={(e) =>
                        handleUpdateGuildMetadata(
                          idx,
                          (dest.metadata.guild_id as string) || "",
                          e.target.value
                        )
                      }
                      className={`px-3 py-2 rounded border border-border bg-background text-foreground placeholder-muted-foreground ${
                        compact ? "text-xs" : "text-sm"
                      }`}
                    />
                  </div>
                </div>
              )}

              {dest.type === "webhook" && (
                <div className="space-y-2">
                  <div className="flex items-center justify-between gap-2">
                    <div className="flex items-center gap-2 flex-1 min-w-0">
                      <Link2 className="w-4 h-4 text-accent flex-shrink-0" />
                      <p
                        className={`font-semibold text-foreground truncate ${
                          compact ? "text-xs" : "text-sm"
                        }`}
                      >
                        {t("reminderCreation.destinations.webhook")}
                      </p>
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
                  <input
                    type="text"
                    placeholder={t(
                      "reminderCreation.destinations.webhookPlaceholder"
                    )}
                    value={(dest.metadata.url as string) || ""}
                    onChange={(e) =>
                      handleUpdateWebhookMetadata(idx, e.target.value)
                    }
                    className={`w-full px-3 py-2 rounded border border-border bg-background text-foreground placeholder-muted-foreground ${
                      compact ? "text-xs" : "text-sm"
                    }`}
                  />
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

          {/* Discord DM Option */}
          <Card
            className="border-border bg-secondary/20 hover:border-accent/50 cursor-pointer transition-colors"
            onClick={handleAddDiscordDM}
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
                      {hasDiscordIdentity
                        ? t("reminderCreation.destinations.sendAsDM")
                        : t(
                            "reminderCreation.destinations.connectDiscordFirst"
                          )}
                    </p>
                  )}
                </div>
              </div>
              <Button
                onClick={(e) => {
                  e.stopPropagation();
                  handleAddDiscordDM();
                }}
                variant="outline"
                size="sm"
                disabled={!hasDiscordIdentity}
                className={`border-accent/50 flex-shrink-0 ${
                  hasDiscordIdentity
                    ? "text-accent hover:bg-accent/10"
                    : "text-muted-foreground opacity-50 cursor-not-allowed"
                }`}
              >
                <Plus className="w-4 h-4" />
              </Button>
            </div>
          </Card>

          {/* Discord Guild Option */}
          <Card className="border-border bg-secondary/20 hover:border-accent/50 cursor-pointer transition-colors">
            <div
              onClick={handleAddDiscordGuild}
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
                      {t("reminderCreation.destinations.discordGuildDesc")}
                    </p>
                  )}
                </div>
              </div>
              <Button
                onClick={(e) => {
                  e.stopPropagation();
                  handleAddDiscordGuild();
                }}
                variant="outline"
                size="sm"
                className="border-accent/50 text-accent hover:bg-accent/10 flex-shrink-0"
              >
                <Plus className="w-4 h-4" />
              </Button>
            </div>
          </Card>

          {/* Webhook Option */}
          <Card className="border-border bg-secondary/20">
            <div className="p-3 space-y-2">
              <div className="flex items-center gap-2">
                <div className="w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center flex-shrink-0">
                  <Link2 className="w-4 h-4 text-accent" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-semibold text-foreground">
                    {t("reminderCreation.destinations.webhook")}
                  </p>
                  {!compact && (
                    <p className="text-xs text-muted-foreground">
                      {t("reminderCreation.destinations.webhookDesc")}
                    </p>
                  )}
                </div>
              </div>
              <div className="flex gap-2">
                <input
                  type="text"
                  placeholder={t(
                    "reminderCreation.destinations.webhookPlaceholder"
                  )}
                  value={webhookUrl}
                  onChange={(e) => setWebhookUrl(e.target.value)}
                  className="flex-1 px-3 py-2 rounded border border-border bg-background text-foreground text-sm placeholder-muted-foreground"
                />
                <Button
                  onClick={handleAddWebhook}
                  disabled={!webhookUrl.trim()}
                  className="bg-accent hover:bg-accent/90 text-accent-foreground flex-shrink-0"
                >
                  <Plus className="w-4 h-4" />
                </Button>
              </div>
            </div>
          </Card>
        </div>
      )}
    </div>
  );
}
