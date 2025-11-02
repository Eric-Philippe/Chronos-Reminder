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

interface LoginSectionProps {
  email: string;
  setEmail: (email: string) => void;
  password: string;
  setPassword: (password: string) => void;
  rememberMe: boolean;
  setRememberMe: (rememberMe: boolean) => void;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => Promise<void>;
  isLoading: boolean;
  localError: string | null;
  authError: { message?: string } | null;
}

export function LoginSection({
  email,
  setEmail,
  password,
  setPassword,
  rememberMe,
  setRememberMe,
  onSubmit,
  isLoading,
  localError,
  authError,
}: LoginSectionProps) {
  const { t } = useTranslation();
  const [showPassword, setShowPassword] = useState(false);

  return (
    <Card className="border-border bg-card/95 backdrop-blur">
      <CardHeader className="space-y-1">
        <CardTitle className="text-foreground">
          {t("login.loginTitle")}
        </CardTitle>
        <CardDescription>{t("login.loginDesc")}</CardDescription>
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
              htmlFor="email"
              className="text-sm font-medium text-foreground"
            >
              {t("login.emailAddress")}
            </label>
            <Input
              id="email"
              type="email"
              placeholder="timely@yours.com"
              value={email}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setEmail(e.target.value)
              }
              className="bg-secondary/50 border-border text-foreground placeholder:text-muted-foreground"
              required
            />
          </div>

          {/* Password Field */}
          <div className="space-y-2">
            <label
              htmlFor="password"
              className="text-sm font-medium text-foreground"
            >
              {t("login.password")}
            </label>
            <div className="relative">
              <Input
                id="password"
                type={showPassword ? "text" : "password"}
                placeholder={t("login.passwordPlaceholder")}
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
          </div>

          {/* Remember Me & Forgot Password */}
          <div className="flex items-center justify-between text-sm">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={rememberMe}
                onChange={(e) => setRememberMe(e.target.checked)}
                className="w-4 h-4 rounded border-border bg-secondary/50 text-accent accent-accent"
              />
              <span className="text-muted-foreground">
                {t("login.rememberMe")}
              </span>
            </label>
            <a
              href="#"
              className="text-accent hover:text-accent/80 transition-colors"
            >
              {t("login.forgotPassword")}
            </a>
          </div>

          {/* Login Button */}
          <Button
            type="submit"
            disabled={isLoading}
            className="w-full bg-accent hover:bg-accent/90 disabled:opacity-50 disabled:cursor-not-allowed text-accent-foreground font-semibold mt-6"
          >
            {isLoading ? t("login.signingIn") : t("login.signIn")}
          </Button>
        </form>

        {/* Divider */}
        <div className="relative my-6">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-border"></div>
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="px-2 bg-card text-muted-foreground">
              {t("login.continueWith")}
            </span>
          </div>
        </div>

        {/* Discord OAuth Button */}
        <Button
          type="button"
          variant="outline"
          onClick={() => {
            const clientId = import.meta.env.VITE_DISCORD_CLIENT_ID;
            const redirectUri = import.meta.env.VITE_DISCORD_REDIRECT_URI;

            console.log("OAuth Client ID:", clientId);
            console.log("OAuth Redirect URI:", redirectUri);

            if (!clientId || !redirectUri) {
              console.error(
                "Discord OAuth configuration is missing. Please check your environment variables."
              );
              return;
            }

            const discordAuthUrl = `https://discord.com/api/oauth2/authorize?client_id=${clientId}&redirect_uri=${encodeURIComponent(
              redirectUri
            )}&response_type=code&scope=identify%20email%20guilds%20guilds.members.read`;
            window.location.href = discordAuthUrl;
          }}
          className="w-full border-border text-foreground hover:bg-secondary/50 hover:text-foreground"
        >
          <svg className="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 24 24">
            <path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515a.074.074 0 0 0-.079.037c-.211.375-.444.864-.607 1.25a18.27 18.27 0 0 0-5.487 0c-.163-.386-.395-.875-.607-1.25a.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057a19.9 19.9 0 0 0 5.993 3.03a.078.078 0 0 0 .084-.028c.462-.63.873-1.295 1.226-1.994a.076.076 0 0 0-.041-.106a13.107 13.107 0 0 1-1.872-.892a.077.077 0 0 1-.008-.128c.126-.094.252-.192.372-.29a.074.074 0 0 1 .076-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .076.01c.12.098.246.196.372.29a.077.077 0 0 1-.006.127a12.299 12.299 0 0 1-1.873.892a.077.077 0 0 0-.041.107c.36.698.77 1.363 1.225 1.993a.076.076 0 0 0 .084.028a19.963 19.963 0 0 0 6.002-3.03a.077.077 0 0 0 .032-.054c.5-4.467.151-8.343-.71-12.382a.06.06 0 0 0-.031-.028zM8.02 15.33c-1.183 0-2.157-.965-2.157-2.156c0-1.193.964-2.157 2.157-2.157c1.193 0 2.168.964 2.157 2.157c0 1.19-.964 2.156-2.157 2.156zm7.975 0c-1.183 0-2.157-.965-2.157-2.156c0-1.193.965-2.157 2.157-2.157c1.192 0 2.167.964 2.157 2.157c0 1.19-.965 2.156-2.157 2.156z" />
          </svg>
          {t("login.discord")}
        </Button>
      </CardContent>
    </Card>
  );
}
