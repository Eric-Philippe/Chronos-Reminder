import { useState } from "react";
import { Trash2, Plus, MessageCircle, Megaphone, Link2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { useTranslation } from "react-i18next";
import type { ReminderFormData } from "@/pages/CreateReminderPage";

interface DestinationsStepProps {
  formData: ReminderFormData;
  onFormChange: (data: ReminderFormData) => void;
}

export function DestinationsStep({
  formData,
  onFormChange,
}: DestinationsStepProps) {
  const { t } = useTranslation();
  const [webhookUrl, setWebhookUrl] = useState("");

  const handleAddDiscordDM = () => {
    const newDest = {
      type: "discord_dm" as const,
      metadata: {},
    };
    onFormChange({
      ...formData,
      destinations: [...formData.destinations, newDest],
    });
  };

  const handleAddDiscordGuild = () => {
    const newDest = {
      type: "discord_channel" as const,
      metadata: {
        guild_id: "",
        channel_id: "",
      },
    };
    onFormChange({
      ...formData,
      destinations: [...formData.destinations, newDest],
    });
  };

  const handleAddWebhook = () => {
    if (!webhookUrl.trim()) return;
    const newDest = {
      type: "webhook" as const,
      metadata: {
        url: webhookUrl,
      },
    };
    onFormChange({
      ...formData,
      destinations: [...formData.destinations, newDest],
    });
    setWebhookUrl("");
  };

  const handleRemoveDestination = (index: number) => {
    onFormChange({
      ...formData,
      destinations: formData.destinations.filter((_, i) => i !== index),
    });
  };

  const handleUpdateWebhookMetadata = (index: number, url: string) => {
    const updated = [...formData.destinations];
    updated[index].metadata = { url };
    onFormChange({ ...formData, destinations: updated });
  };

  const handleUpdateGuildMetadata = (
    index: number,
    guildId: string,
    channelId: string
  ) => {
    const updated = [...formData.destinations];
    updated[index].metadata = { guild_id: guildId, channel_id: channelId };
    onFormChange({ ...formData, destinations: updated });
  };

  const hasDiscordIdentity = true; // TODO: Get from context/props

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold text-foreground mb-4">
          {t("reminderCreation.destinations.title")}
        </h2>
        <p className="text-sm text-muted-foreground mb-6">
          {t("reminderCreation.destinations.subtitle")}
        </p>
      </div>

      {/* Selected Destinations */}
      {formData.destinations.length > 0 && (
        <div className="space-y-3">
          <h3 className="text-sm font-semibold text-foreground">
            {t("reminderCreation.destinations.addedDestinations")} (
            {formData.destinations.length})
          </h3>
          {formData.destinations.map((dest, idx) => (
            <div
              key={idx}
              className="p-4 rounded-lg border border-border bg-secondary/20 space-y-3"
            >
              {dest.type === "discord_dm" && (
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <MessageCircle className="w-5 h-5 text-accent flex-shrink-0" />
                    <div>
                      <p className="text-sm font-semibold text-foreground">
                        {t("reminderCreation.destinations.discordDM")}
                      </p>
                      <p className="text-xs text-muted-foreground mt-1">
                        {t("reminderCreation.destinations.discordDMDesc")}
                      </p>
                    </div>
                  </div>
                  <Button
                    onClick={() => handleRemoveDestination(idx)}
                    variant="outline"
                    size="sm"
                    className="border-red-500/50 text-red-600 dark:text-red-400 hover:bg-red-500/10"
                  >
                    <Trash2 className="w-3 h-3" />
                  </Button>
                </div>
              )}

              {dest.type === "discord_channel" && (
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <Megaphone className="w-5 h-5 text-accent flex-shrink-0" />
                      <p className="text-sm font-semibold text-foreground">
                        {t("reminderCreation.destinations.discordGuild")}
                      </p>
                    </div>
                    <Button
                      onClick={() => handleRemoveDestination(idx)}
                      variant="outline"
                      size="sm"
                      className="border-red-500/50 text-red-600 dark:text-red-400 hover:bg-red-500/10"
                    >
                      <Trash2 className="w-3 h-3" />
                    </Button>
                  </div>
                  <div className="grid grid-cols-2 gap-3">
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
                      className="px-3 py-2 rounded border border-border bg-background text-foreground text-sm placeholder-muted-foreground"
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
                      className="px-3 py-2 rounded border border-border bg-background text-foreground text-sm placeholder-muted-foreground"
                    />
                  </div>
                </div>
              )}

              {dest.type === "webhook" && (
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <Link2 className="w-5 h-5 text-accent flex-shrink-0" />
                      <p className="text-sm font-semibold text-foreground">
                        {t("reminderCreation.destinations.webhook")}
                      </p>
                    </div>
                    <Button
                      onClick={() => handleRemoveDestination(idx)}
                      variant="outline"
                      size="sm"
                      className="border-red-500/50 text-red-600 dark:text-red-400 hover:bg-red-500/10"
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
                    className="w-full px-3 py-2 rounded border border-border bg-background text-foreground text-sm placeholder-muted-foreground"
                  />
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Add Destination Options */}
      <div className="space-y-3">
        <h3 className="text-sm font-semibold text-foreground">
          {t("reminderCreation.destinations.addDestination")}
        </h3>

        {/* Discord DM Option */}
        <Card className="border-border bg-secondary/20 hover:border-accent/50 cursor-pointer transition-colors">
          <div
            onClick={handleAddDiscordDM}
            className="p-4 flex items-center justify-between"
          >
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-accent/10 flex items-center justify-center">
                <MessageCircle className="w-5 h-5 text-accent" />
              </div>
              <div>
                <p className="text-sm font-semibold text-foreground">
                  {t("reminderCreation.destinations.discordDM")}
                </p>
                <p className="text-xs text-muted-foreground">
                  {hasDiscordIdentity
                    ? t("reminderCreation.destinations.sendAsDM")
                    : t("reminderCreation.destinations.connectDiscordFirst")}
                </p>
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
              className={`border-accent/50 ${
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
            className="p-4 flex items-center justify-between"
          >
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-accent/10 flex items-center justify-center">
                <Megaphone className="w-5 h-5 text-accent" />
              </div>
              <div>
                <p className="text-sm font-semibold text-foreground">
                  {t("reminderCreation.destinations.discordGuild")}
                </p>
                <p className="text-xs text-muted-foreground">
                  {t("reminderCreation.destinations.discordGuildDesc")}
                </p>
              </div>
            </div>
            <Button
              onClick={(e) => {
                e.stopPropagation();
                handleAddDiscordGuild();
              }}
              variant="outline"
              size="sm"
              className="border-accent/50 text-accent hover:bg-accent/10"
            >
              <Plus className="w-4 h-4" />
            </Button>
          </div>
        </Card>

        {/* Webhook Option */}
        <Card className="border-border bg-secondary/20">
          <div className="p-4 space-y-3">
            <div className="flex items-center gap-3 mb-3">
              <div className="w-10 h-10 rounded-lg bg-accent/10 flex items-center justify-center">
                <Link2 className="w-5 h-5 text-accent" />
              </div>
              <div>
                <p className="text-sm font-semibold text-foreground">
                  {t("reminderCreation.destinations.webhook")}
                </p>
                <p className="text-xs text-muted-foreground">
                  {t("reminderCreation.destinations.webhookDesc")}
                </p>
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
                className="bg-accent hover:bg-accent/90 text-accent-foreground"
              >
                <Plus className="w-4 h-4" />
              </Button>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
}
