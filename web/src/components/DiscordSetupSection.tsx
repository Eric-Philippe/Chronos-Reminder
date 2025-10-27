import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Eye, EyeOff } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { TimezoneSelect } from "@/components/common/TimezoneSelect";
import { PasswordStrengthIndicator } from "@/components/common/PasswordStrengthIndicator";

interface DiscordSetupSectionProps {
  discordEmail: string;
  discordUsername: string;
  onSetupComplete: (data: {
    email: string;
    username: string;
    password: string;
    timezone: string;
  }) => Promise<void>;
  isLoading: boolean;
  error: string | null;
}

export function DiscordSetupSection({
  discordEmail,
  discordUsername,
  onSetupComplete,
  isLoading,
  error,
}: DiscordSetupSectionProps) {
  const { t } = useTranslation();
  const [showPassword, setShowPassword] = useState(false);
  const [email, setEmail] = useState(discordEmail);
  const [username, setUsername] = useState(discordUsername);
  const [password, setPassword] = useState("");
  const [timezone, setTimezone] = useState("UTC");
  const [localError, setLocalError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setLocalError(null);

    if (!email.trim()) {
      setLocalError(
        t("login.discordSetup.emailLabel") +
          " " +
          t("login.discordSetup.passwordRequired")
      );
      return;
    }

    if (!username.trim()) {
      setLocalError(t("login.discordSetup.usernameRequired"));
      return;
    }

    if (!password.trim()) {
      setLocalError(t("login.discordSetup.passwordRequired"));
      return;
    }

    if (password.length < 8) {
      setLocalError(t("login.discordSetup.passwordMin"));
      return;
    }

    try {
      await onSetupComplete({ email, username, password, timezone });
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : "Setup failed";
      setLocalError(errorMsg);
    }
  };

  return (
    <Card className="border-border bg-card/95 backdrop-blur">
      <CardHeader className="space-y-1">
        <CardTitle className="text-foreground">
          {t("login.discordSetup.title")}
        </CardTitle>
        <CardDescription>
          {t("login.discordSetup.subtitle", { username: discordUsername })}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Error Alert */}
          {(localError || error) && (
            <div className="bg-red-500/10 border border-red-500 rounded-md p-3 text-sm text-red-600">
              {localError || error}
            </div>
          )}

          {/* Email Field (Pre-filled from Discord) */}
          <div className="space-y-2">
            <label
              htmlFor="setup-email"
              className="text-sm font-medium text-foreground"
            >
              {t("login.discordSetup.emailLabel")}
            </label>
            <Input
              id="setup-email"
              type="email"
              placeholder={t("login.discordSetup.emailPlaceholder")}
              value={email}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setEmail(e.target.value)
              }
              className="bg-secondary/50 border-border text-foreground placeholder:text-muted-foreground"
              required
            />
            <p className="text-xs text-muted-foreground">
              {t("login.discordSetup.emailHelper")}
            </p>
          </div>

          {/* Username Field (App Identity) */}
          <div className="space-y-2">
            <label
              htmlFor="setup-username"
              className="text-sm font-medium text-foreground"
            >
              {t("login.discordSetup.usernameLabel")}
            </label>
            <Input
              id="setup-username"
              type="text"
              placeholder={t("login.discordSetup.usernamePlaceholder")}
              value={username}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setUsername(e.target.value)
              }
              className="bg-secondary/50 border-border text-foreground placeholder:text-muted-foreground"
              required
            />
            <p className="text-xs text-muted-foreground">
              {t("login.discordSetup.usernameHelper")}
            </p>
          </div>

          {/* Password Field */}
          <div className="space-y-2">
            <label
              htmlFor="setup-password"
              className="text-sm font-medium text-foreground"
            >
              {t("login.discordSetup.passwordLabel")}
            </label>
            <div className="relative">
              <Input
                id="setup-password"
                type={showPassword ? "text" : "password"}
                placeholder={t("login.discordSetup.passwordPlaceholder")}
                value={password}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                  setPassword(e.target.value)
                }
                className="bg-secondary/50 border-border text-foreground placeholder:text-muted-foreground pr-10"
                required
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
              >
                {showPassword ? (
                  <EyeOff className="w-4 h-4" />
                ) : (
                  <Eye className="w-4 h-4" />
                )}
              </button>
            </div>

            {/* Password Strength Indicator */}
            {password && <PasswordStrengthIndicator password={password} />}
          </div>

          {/* Timezone Field */}
          <div className="space-y-2">
            <label
              htmlFor="setup-timezone"
              className="text-sm font-medium text-foreground"
            >
              {t("login.discordSetup.timezoneLabel")}
            </label>
            <TimezoneSelect
              value={timezone}
              onChange={setTimezone}
              searchPlaceholder={t("login.discordSetup.searchTimezone")}
              noResultsText={t("login.discordSetup.noTimezones")}
            />
          </div>

          {/* Submit Button */}
          <Button
            type="submit"
            disabled={isLoading}
            className="w-full bg-accent hover:bg-accent/90 disabled:opacity-50 disabled:cursor-not-allowed text-accent-foreground font-semibold mt-6"
          >
            {isLoading
              ? t("login.discordSetup.buttonCompleting")
              : t("login.discordSetup.buttonComplete")}
          </Button>
        </form>

        {/* Info Box */}
        <div className="mt-6 p-4 bg-accent/10 rounded-lg border border-accent/20">
          <p className="text-sm text-foreground">
            <strong>💡 {t("login.discordSetup.tip")}</strong>
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
