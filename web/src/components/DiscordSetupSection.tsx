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
import { TimezoneSelect } from "@/components/TimezoneSelect";
import { PasswordStrengthIndicator } from "@/components/PasswordStrengthIndicator";

interface DiscordSetupSectionProps {
  discordEmail: string;
  discordUsername: string;
  onSetupComplete: (data: {
    email: string;
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
  const [password, setPassword] = useState("");
  const [timezone, setTimezone] = useState("UTC");
  const [localError, setLocalError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setLocalError(null);

    if (!email.trim()) {
      setLocalError(t("login.emailAddress") + " is required");
      return;
    }

    if (!password.trim()) {
      setLocalError(t("login.password") + " is required");
      return;
    }

    if (password.length < 8) {
      setLocalError("Password must be at least 8 characters");
      return;
    }

    try {
      await onSetupComplete({ email, password, timezone });
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : "Setup failed";
      setLocalError(errorMsg);
    }
  };

  return (
    <Card className="border-border bg-card/95 backdrop-blur">
      <CardHeader className="space-y-1">
        <CardTitle className="text-foreground">Complete Your Profile</CardTitle>
        <CardDescription>
          Welcome, {discordUsername}! Set up your app identity to secure your
          account.
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
              Email Address
            </label>
            <Input
              id="setup-email"
              type="email"
              placeholder="your@email.com"
              value={email}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setEmail(e.target.value)
              }
              className="bg-secondary/50 border-border text-foreground placeholder:text-muted-foreground"
              required
            />
            <p className="text-xs text-muted-foreground">
              Pre-filled from your Discord account. You can change this if
              needed.
            </p>
          </div>

          {/* Password Field */}
          <div className="space-y-2">
            <label
              htmlFor="setup-password"
              className="text-sm font-medium text-foreground"
            >
              Password
            </label>
            <div className="relative">
              <Input
                id="setup-password"
                type={showPassword ? "text" : "password"}
                placeholder="Create a strong password"
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
              Timezone
            </label>
            <TimezoneSelect value={timezone} onChange={setTimezone} />
          </div>

          {/* Submit Button */}
          <Button
            type="submit"
            disabled={isLoading}
            className="w-full bg-accent hover:bg-accent/90 disabled:opacity-50 disabled:cursor-not-allowed text-accent-foreground font-semibold mt-6"
          >
            {isLoading ? "Completing Setup..." : "Complete Setup"}
          </Button>
        </form>

        {/* Info Box */}
        <div className="mt-6 p-4 bg-accent/10 rounded-lg border border-accent/20">
          <p className="text-sm text-foreground">
            <strong>💡 Tip:</strong> By setting up a password, you can log in
            using either your Discord account or your email address.
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
