import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useSearchParams } from "react-router-dom";
import { Eye, EyeOff } from "lucide-react";
import { LanguageSwitcher } from "@/components/common/language-switcher";
import { ModeToggle } from "@/components/common/mode-toggle";
import { useToast } from "@/hooks/useToast";
import { authService } from "@/services";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

type VerificationStatus = "verifying" | "valid" | "invalid" | "expired";

export function ResetPasswordPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { success, error: showError } = useToast();

  const email = searchParams.get("email");
  const token = searchParams.get("token");

  const [verificationStatus, setVerificationStatus] =
    useState<VerificationStatus>("verifying");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [localError, setLocalError] = useState<string | null>(null);
  const [isSuccess, setIsSuccess] = useState(false);

  // Verify token on mount
  useEffect(() => {
    if (!email || !token) {
      setVerificationStatus("invalid");
      return;
    }

    const verifyToken = async () => {
      try {
        const response = await authService.verifyResetToken({
          email,
          token,
        });

        if (response.valid) {
          setVerificationStatus("valid");
        } else {
          setVerificationStatus("expired");
        }
      } catch (err) {
        console.error("Token verification error:", err);
        setVerificationStatus("expired");
      }
    };

    verifyToken();
  }, [email, token]);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setLocalError(null);

    if (!password.trim()) {
      setLocalError(t("common.passwordRequired"));
      return;
    }

    if (password.length < 8) {
      setLocalError(t("validation.passwordTooShort"));
      return;
    }

    if (password !== confirmPassword) {
      setLocalError(t("common.passwordsDontMatch"));
      return;
    }

    if (!email || !token) {
      setLocalError(t("resetPassword.invalidRequest"));
      return;
    }

    setIsLoading(true);

    try {
      await authService.resetPassword({
        email,
        token,
        password,
      });

      // Show success message
      success(t("resetPassword.success") as string);
      setIsSuccess(true);

      // Redirect to login after 3 seconds
      setTimeout(() => {
        navigate("/login", { replace: true });
      }, 3000);
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : t("resetPassword.error");
      setLocalError(errorMsg);
      showError(errorMsg);
      setVerificationStatus("expired");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background-main to-background-secondary flex items-center justify-center p-4 relative">
      {/* Background decorative elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-0 right-0 w-72 h-72 bg-accent/10 rounded-full blur-3xl dark:bg-accent/5"></div>
        <div className="absolute bottom-0 left-0 w-96 h-96 bg-accent/10 rounded-full blur-3xl dark:bg-accent/5"></div>
      </div>

      {/* Theme and Language Controls */}
      <div className="absolute top-4 right-4 z-50 flex gap-2">
        <LanguageSwitcher />
        <ModeToggle />
      </div>

      {/* Back to Login Button */}
      <button
        onClick={() => navigate("/login")}
        className="absolute top-4 left-4 z-50 px-4 py-2 rounded-md text-foreground dark:text-white hover:text-amber-600 dark:hover:text-amber-400 hover:bg-amber-400/10 transition-colors border border-border dark:border-white/10"
      >
        ← {t("common.back") || "Back"}
      </button>

      <div className="w-full max-w-md relative z-10">
        <div className="flex flex-col justify-center items-center space-y-8">
          {/* Header */}
          <div className="text-center">
            <div className="flex items-center justify-center gap-3 mb-6">
              <img
                src="/logo_chronos.png"
                alt="Chronos Logo"
                className="w-12 h-12 rounded-full"
              />
              <span className="text-3xl font-bold text-accent">Chronos</span>
            </div>
            <h1 className="text-3xl font-bold text-foreground mb-2">
              {t("resetPassword.title")}
            </h1>
          </div>

          {/* Verifying State */}
          {verificationStatus === "verifying" && (
            <Card className="w-full border-border bg-card/95 backdrop-blur">
              <CardContent className="pt-6">
                <div className="text-center space-y-4">
                  <div className="animate-pulse">
                    {t("common.verifying")}...
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Valid Token - Form */}
          {verificationStatus === "valid" && !isSuccess && (
            <Card className="w-full border-border bg-card/95 backdrop-blur">
              <CardHeader className="space-y-1">
                <CardTitle className="text-foreground">
                  {t("resetPassword.enterNewPassword")}
                </CardTitle>
                <CardDescription>
                  {t("resetPassword.passwordRequirements")}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                  {/* Error Alert */}
                  {localError && (
                    <div className="bg-red-500/10 border border-red-500 rounded-md p-3 text-sm text-red-600">
                      {localError}
                    </div>
                  )}

                  {/* New Password Field */}
                  <div className="space-y-2">
                    <label
                      htmlFor="password"
                      className="text-sm font-medium text-foreground"
                    >
                      {t("common.newPassword")}
                    </label>
                    <div className="relative">
                      <Input
                        id="password"
                        type={showPassword ? "text" : "password"}
                        placeholder={t("common.passwordPlaceholder")}
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className="bg-secondary/50 border-border text-foreground placeholder:text-muted-foreground pr-10"
                        required
                        disabled={isLoading}
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

                  {/* Confirm Password Field */}
                  <div className="space-y-2">
                    <label
                      htmlFor="confirm-password"
                      className="text-sm font-medium text-foreground"
                    >
                      {t("common.confirmPassword")}
                    </label>
                    <div className="relative">
                      <Input
                        id="confirm-password"
                        type={showConfirmPassword ? "text" : "password"}
                        placeholder={t("common.passwordPlaceholder")}
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        className="bg-secondary/50 border-border text-foreground placeholder:text-muted-foreground pr-10"
                        required
                        disabled={isLoading}
                      />
                      <button
                        type="button"
                        onClick={() =>
                          setShowConfirmPassword(!showConfirmPassword)
                        }
                        className="absolute right-3 top-1/2 transform -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                      >
                        {showConfirmPassword ? (
                          <EyeOff className="w-4 h-4" />
                        ) : (
                          <Eye className="w-4 h-4" />
                        )}
                      </button>
                    </div>
                  </div>

                  {/* Submit Button */}
                  <Button
                    type="submit"
                    disabled={isLoading}
                    className="w-full bg-accent hover:bg-accent/90 disabled:opacity-50 disabled:cursor-not-allowed text-accent-foreground font-semibold mt-6"
                  >
                    {isLoading
                      ? t("common.resetting")
                      : t("resetPassword.resetPassword")}
                  </Button>
                </form>
              </CardContent>
            </Card>
          )}

          {/* Success State */}
          {isSuccess && (
            <Card className="w-full border-border bg-card/95 backdrop-blur border-green-500/30">
              <CardContent className="pt-6">
                <div className="text-center space-y-4">
                  <div className="flex justify-center">
                    <div className="w-16 h-16 rounded-full bg-green-500/20 flex items-center justify-center">
                      <span className="text-3xl">✓</span>
                    </div>
                  </div>
                  <h2 className="text-lg font-semibold text-foreground">
                    {t("resetPassword.passwordReset")}
                  </h2>
                  <p className="text-muted-foreground">
                    {t("resetPassword.youCanNowLogin")}
                  </p>
                  <Button
                    onClick={() => navigate("/login", { replace: true })}
                    className="w-full bg-accent hover:bg-accent/90 text-accent-foreground font-semibold mt-4"
                  >
                    {t("common.goToLogin")}
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Invalid/Expired Token */}
          {(verificationStatus === "invalid" ||
            verificationStatus === "expired") && (
            <Card className="w-full border-border bg-card/95 backdrop-blur border-red-500/30">
              <CardContent className="pt-6">
                <div className="text-center space-y-4">
                  <div className="flex justify-center">
                    <div className="w-16 h-16 rounded-full bg-red-500/20 flex items-center justify-center">
                      <span className="text-3xl">✕</span>
                    </div>
                  </div>
                  <h2 className="text-lg font-semibold text-foreground">
                    {verificationStatus === "invalid"
                      ? t("resetPassword.invalidLink")
                      : t("resetPassword.linkExpired")}
                  </h2>
                  <p className="text-muted-foreground">
                    {verificationStatus === "invalid"
                      ? t("resetPassword.invalidLinkMessage")
                      : t("resetPassword.linkExpiredMessage")}
                  </p>
                  <Button
                    onClick={() => navigate("/forgot-password")}
                    className="w-full bg-accent hover:bg-accent/90 text-accent-foreground font-semibold mt-4"
                  >
                    {t("resetPassword.requestNewLink")}
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}
