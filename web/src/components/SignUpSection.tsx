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
import { PasswordStrengthIndicator } from "@/components/PasswordStrengthIndicator";
import { TimezoneSelect } from "@/components/TimezoneSelect";

interface SignUpSectionProps {
  signUpEmail: string;
  setSignUpEmail: (email: string) => void;
  username: string;
  setUsername: (username: string) => void;
  signUpPassword: string;
  setSignUpPassword: (password: string) => void;
  confirmPassword: string;
  setConfirmPassword: (password: string) => void;
  timezone: string;
  setTimezone: (timezone: string) => void;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => Promise<void>;
  isLoading: boolean;
  localError: string | null;
  authError: { message?: string } | null;
}

export function SignUpSection({
  signUpEmail,
  setSignUpEmail,
  username,
  setUsername,
  signUpPassword,
  setSignUpPassword,
  confirmPassword,
  setConfirmPassword,
  timezone,
  setTimezone,
  onSubmit,
  isLoading,
  localError,
  authError,
}: SignUpSectionProps) {
  const { t } = useTranslation();
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  return (
    <Card className="border-border bg-card/95 backdrop-blur">
      <CardHeader className="space-y-1">
        <CardTitle className="text-foreground">
          {t("login.createAccountTitle")}
        </CardTitle>
        <CardDescription>{t("login.createAccountDesc")}</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="space-y-4">
          {/* Error Alert */}
          {(localError || authError?.message) && (
            <div className="bg-red-500/10 border border-red-500 rounded-md p-3 text-sm text-red-600">
              {localError || authError?.message}
            </div>
          )}
          {/* Email Field */}
          <div className="space-y-2">
            <label
              htmlFor="signup-email"
              className="text-sm font-medium text-foreground"
            >
              {t("login.emailAddress")}
            </label>
            <Input
              id="signup-email"
              type="email"
              placeholder="timely@yours.com"
              value={signUpEmail}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setSignUpEmail(e.target.value)
              }
              className="bg-secondary/50 border-border text-foreground placeholder:text-muted-foreground"
              required
            />
          </div>

          {/* Username Field */}
          <div className="space-y-2">
            <label
              htmlFor="username"
              className="text-sm font-medium text-foreground"
            >
              {t("login.username")}
            </label>
            <Input
              id="username"
              type="text"
              placeholder={t("login.usernamePlaceholder")}
              value={username}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setUsername(e.target.value)
              }
              className="bg-secondary/50 border-border text-foreground placeholder:text-muted-foreground"
              required
            />
          </div>

          {/* Password Fields with Strength Indicator */}
          <div className="space-y-2">
            <label className="text-sm font-medium text-foreground">
              {t("login.password")}
            </label>
            <div className="flex gap-3">
              <div className="flex-1 space-y-6">
                {/* Password Field */}
                <div className="relative">
                  <Input
                    id="signup-password"
                    type={showPassword ? "text" : "password"}
                    placeholder={t("login.passwordPlaceholder")}
                    value={signUpPassword}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      setSignUpPassword(e.target.value)
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

                {/* Confirm Password Field */}
                <div className="relative">
                  <Input
                    id="confirm-password"
                    type={showConfirmPassword ? "text" : "password"}
                    placeholder={t("login.passwordPlaceholder")}
                    value={confirmPassword}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      setConfirmPassword(e.target.value)
                    }
                    className={`bg-secondary/50 border-border text-foreground placeholder:text-muted-foreground pr-10 ${
                      confirmPassword && signUpPassword !== confirmPassword
                        ? "border-red-500"
                        : ""
                    }`}
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-3 top-1/2 transform -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                  >
                    {showConfirmPassword ? (
                      <EyeOff className="w-4 h-4" />
                    ) : (
                      <Eye className="w-4 h-4" />
                    )}
                  </button>
                </div>
                {confirmPassword && signUpPassword !== confirmPassword && (
                  <p className="text-xs text-red-500">
                    {t("login.passwordsDontMatch")}
                  </p>
                )}
              </div>
              <div className="hidden lg:flex flex-1 pt-1">
                <PasswordStrengthIndicator password={signUpPassword} />
              </div>
            </div>
          </div>

          {/* Timezone Select */}
          <div className="space-y-2">
            <label
              htmlFor="timezone"
              className="text-sm font-medium text-foreground"
            >
              {t("login.timezone")}
            </label>
            <TimezoneSelect value={timezone} onChange={setTimezone} />
          </div>

          {/* Sign Up Button */}
          <Button
            type="submit"
            disabled={
              signUpPassword !== confirmPassword ||
              !signUpPassword ||
              !signUpEmail ||
              !username ||
              isLoading
            }
            className="w-full bg-accent hover:bg-accent/90 disabled:opacity-50 disabled:cursor-not-allowed text-accent-foreground font-semibold mt-6"
          >
            {isLoading ? t("login.creatingAccount") : t("login.createAccount")}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
