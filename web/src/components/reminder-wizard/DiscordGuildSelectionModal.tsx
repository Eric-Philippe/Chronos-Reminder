import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import {
  Server,
  Hash,
  Users,
  ExternalLink,
  Loader2,
  AlertCircle,
} from "lucide-react";
import { discordService } from "@/services/discord";
import type {
  DiscordGuild,
  DiscordChannel,
  DiscordRole,
} from "@/services/types";
import { useAuth } from "@/hooks/useAuth";

interface DiscordGuildSelectionModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: (guildId: string, channelId: string, roleId?: string) => void;
}

export function DiscordGuildSelectionModal({
  open,
  onOpenChange,
  onConfirm,
}: DiscordGuildSelectionModalProps) {
  const { t } = useTranslation();
  const { user } = useAuth();

  const [step, setStep] = useState<"guild" | "channel" | "role">("guild");
  const [guilds, setGuilds] = useState<DiscordGuild[]>([]);
  const [channels, setChannels] = useState<DiscordChannel[]>([]);
  const [roles, setRoles] = useState<DiscordRole[]>([]);

  const [selectedGuildId, setSelectedGuildId] = useState<string>("");
  const [selectedChannelId, setSelectedChannelId] = useState<string>("");
  const [selectedRoleId, setSelectedRoleId] = useState<string>("");

  const [isLoadingGuilds, setIsLoadingGuilds] = useState(false);
  const [isLoadingChannels, setIsLoadingChannels] = useState(false);
  const [isLoadingRoles, setIsLoadingRoles] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [botNotInGuild, setBotNotInGuild] = useState(false);

  const loadGuilds = async () => {
    if (!user?.user_id) return;

    setIsLoadingGuilds(true);
    setError(null);
    try {
      const userGuilds = await discordService.getUserGuilds(user.user_id);
      setGuilds(userGuilds);
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : "Failed to load guilds";
      setError(errorMsg);
    } finally {
      setIsLoadingGuilds(false);
    }
  };

  // Load guilds when modal opens
  useEffect(() => {
    if (open && user?.user_id) {
      loadGuilds();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, user?.user_id]);

  // Reset state when modal closes
  useEffect(() => {
    if (!open) {
      setStep("guild");
      setSelectedGuildId("");
      setSelectedChannelId("");
      setSelectedRoleId("");
      setChannels([]);
      setRoles([]);
      setError(null);
      setBotNotInGuild(false);
    }
  }, [open]);

  const loadChannels = async (guildId: string) => {
    if (!user?.user_id) return;

    setIsLoadingChannels(true);
    setError(null);
    setBotNotInGuild(false);
    try {
      const guildChannels = await discordService.getGuildChannels(
        user.user_id,
        guildId
      );
      setChannels(guildChannels);
      setStep("channel");
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : "Failed to load channels";

      // Check if error is due to bot not being in guild
      if (errorMsg.includes("Bot is not a member of this guild")) {
        setBotNotInGuild(true);
      }
      setError(errorMsg);
    } finally {
      setIsLoadingChannels(false);
    }
  };

  const loadRoles = async (guildId: string) => {
    if (!user?.user_id) return;

    setIsLoadingRoles(true);
    setError(null);
    try {
      const guildRoles = await discordService.getGuildRoles(
        user.user_id,
        guildId
      );
      // Filter out @everyone role and managed roles
      const filterRoles = guildRoles.filter(
        (role) => role.name !== "@everyone" && !role.managed
      );
      setRoles(filterRoles);
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : "Failed to load roles";
      setError(errorMsg);
    } finally {
      setIsLoadingRoles(false);
    }
  };

  const handleGuildSelect = async (guildId: string) => {
    setSelectedGuildId(guildId);
    await loadChannels(guildId);
    // Also preload roles
    await loadRoles(guildId);
  };

  const handleChannelSelect = (channelId: string) => {
    setSelectedChannelId(channelId);
  };

  const handleProceedToRole = () => {
    setStep("role");
  };

  const handleSkipRole = () => {
    onConfirm(selectedGuildId, selectedChannelId);
    onOpenChange(false);
  };

  const handleConfirm = () => {
    if (step === "role" && selectedRoleId) {
      onConfirm(selectedGuildId, selectedChannelId, selectedRoleId);
    } else {
      onConfirm(selectedGuildId, selectedChannelId);
    }
    onOpenChange(false);
  };

  const handleBack = () => {
    if (step === "role") {
      setStep("channel");
      setSelectedRoleId("");
    } else if (step === "channel") {
      setStep("guild");
      setSelectedChannelId("");
      setSelectedRoleId("");
      setChannels([]);
      setRoles([]);
      setBotNotInGuild(false);
    }
  };

  const handleInviteBot = () => {
    if (selectedGuildId) {
      const inviteUrl = discordService.getBotInviteUrl();
      window.open(inviteUrl, "_blank", "noopener,noreferrer");
    }
  };

  const getSelectedGuild = () => guilds.find((g) => g.id === selectedGuildId);
  const getSelectedChannel = () =>
    channels.find((c) => c.id === selectedChannelId);

  const canProceed = () => {
    if (step === "guild") return selectedGuildId !== "";
    if (step === "channel") return selectedChannelId !== "";
    if (step === "role") return true; // Role is optional
    return false;
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            {step === "guild" && (
              <>
                <Server className="w-5 h-5" />
                {t("reminderCreation.discord.selectGuild")}
              </>
            )}
            {step === "channel" && (
              <>
                <Hash className="w-5 h-5" />
                {t("reminderCreation.discord.selectChannel")}
              </>
            )}
            {step === "role" && (
              <>
                <Users className="w-5 h-5" />
                {t("reminderCreation.discord.selectRole")}
              </>
            )}
          </DialogTitle>
          <DialogDescription>
            {step === "guild" &&
              t("reminderCreation.discord.selectGuildDescription")}
            {step === "channel" &&
              t("reminderCreation.discord.selectChannelDescription")}
            {step === "role" &&
              t("reminderCreation.discord.selectRoleDescription")}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {/* Error Alert */}
          {error && !botNotInGuild && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          {/* Bot Not In Guild Alert */}
          {botNotInGuild && (
            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertDescription className="space-y-3">
                <p>{t("reminderCreation.discord.botNotInGuild")}</p>
                <Button
                  onClick={handleInviteBot}
                  variant="outline"
                  size="sm"
                  className="w-full gap-2"
                >
                  <ExternalLink className="w-4 h-4" />
                  {t("reminderCreation.discord.inviteBot")}
                </Button>
              </AlertDescription>
            </Alert>
          )}

          {/* Guild Selection */}
          {step === "guild" && (
            <div className="space-y-2">
              <Label htmlFor="guild">
                {t("reminderCreation.discord.guild")}
              </Label>
              {isLoadingGuilds ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
                </div>
              ) : (
                <Select
                  value={selectedGuildId}
                  onValueChange={handleGuildSelect}
                  disabled={isLoadingChannels}
                >
                  <SelectTrigger id="guild">
                    <SelectValue
                      placeholder={t(
                        "reminderCreation.discord.selectGuildPlaceholder"
                      )}
                    />
                  </SelectTrigger>
                  <SelectContent>
                    {guilds.map((guild) => (
                      <SelectItem key={guild.id} value={guild.id}>
                        <div className="flex items-center gap-2">
                          {guild.icon ? (
                            <img
                              src={`https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png`}
                              alt={guild.name}
                              className="w-5 h-5 rounded-full"
                            />
                          ) : (
                            <Server className="w-5 h-5" />
                          )}
                          <span>{guild.name}</span>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            </div>
          )}

          {/* Channel Selection */}
          {step === "channel" && !botNotInGuild && (
            <div className="space-y-2">
              <Label htmlFor="channel">
                {t("reminderCreation.discord.channel")}
              </Label>
              <div className="text-xs text-muted-foreground mb-2">
                {t("reminderCreation.discord.selectedGuild")}:{" "}
                <strong>{getSelectedGuild()?.name}</strong>
              </div>
              {isLoadingChannels ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
                </div>
              ) : (
                <Select
                  value={selectedChannelId}
                  onValueChange={handleChannelSelect}
                >
                  <SelectTrigger id="channel">
                    <SelectValue
                      placeholder={t(
                        "reminderCreation.discord.selectChannelPlaceholder"
                      )}
                    />
                  </SelectTrigger>
                  <SelectContent>
                    {channels
                      .filter((channel) => channel.type === 0) // Only text channels
                      .map((channel) => (
                        <SelectItem key={channel.id} value={channel.id}>
                          <div className="flex items-center gap-2">
                            <Hash className="w-4 h-4" />
                            <span>{channel.name}</span>
                          </div>
                        </SelectItem>
                      ))}
                  </SelectContent>
                </Select>
              )}
            </div>
          )}

          {/* Role Selection */}
          {step === "role" && (
            <div className="space-y-2">
              <Label htmlFor="role">
                {t("reminderCreation.discord.role")}{" "}
                <span className="text-muted-foreground">
                  ({t("reminderCreation.discord.optional")})
                </span>
              </Label>
              <div className="text-xs text-muted-foreground mb-2">
                {t("reminderCreation.discord.selectedChannel")}:{" "}
                <strong>#{getSelectedChannel()?.name}</strong>
              </div>
              {isLoadingRoles ? (
                <div className="flex items-center justify-center py-8">
                  <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
                </div>
              ) : (
                <Select
                  value={selectedRoleId}
                  onValueChange={setSelectedRoleId}
                >
                  <SelectTrigger id="role">
                    <SelectValue
                      placeholder={t(
                        "reminderCreation.discord.selectRolePlaceholder"
                      )}
                    />
                  </SelectTrigger>
                  <SelectContent>
                    {roles.map((role) => (
                      <SelectItem key={role.id} value={role.id}>
                        <div className="flex items-center gap-2">
                          <div
                            className="w-3 h-3 rounded-full"
                            style={{
                              backgroundColor:
                                role.color !== 0
                                  ? `#${role.color
                                      .toString(16)
                                      .padStart(6, "0")}`
                                  : "#99aab5",
                            }}
                          />
                          <span>{role.name}</span>
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            </div>
          )}
        </div>

        <DialogFooter className="gap-2 sm:gap-2">
          {step !== "guild" && (
            <Button
              type="button"
              variant="outline"
              onClick={handleBack}
              disabled={isLoadingChannels || isLoadingRoles}
            >
              {t("reminderCreation.buttons.back")}
            </Button>
          )}

          {step === "channel" && selectedChannelId && (
            <>
              <Button
                type="button"
                variant="outline"
                onClick={handleSkipRole}
                disabled={isLoadingChannels}
              >
                {t("reminderCreation.discord.skipRole")}
              </Button>
              <Button
                type="button"
                onClick={handleProceedToRole}
                disabled={!canProceed() || isLoadingChannels}
              >
                {t("reminderCreation.discord.addRole")}
              </Button>
            </>
          )}

          {step === "role" && (
            <Button
              type="button"
              onClick={handleConfirm}
              disabled={isLoadingRoles}
            >
              {t("reminderCreation.buttons.create")}
            </Button>
          )}

          {step === "guild" && !botNotInGuild && (
            <Button
              type="button"
              onClick={() => selectedGuildId && loadChannels(selectedGuildId)}
              disabled={!canProceed() || isLoadingGuilds || isLoadingChannels}
            >
              {isLoadingChannels && (
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              )}
              {t("reminderCreation.buttons.next")}
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
