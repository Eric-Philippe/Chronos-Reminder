import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { Header } from "@/components/common/header";
import { Footer } from "@/components/common/footer";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert } from "@/components/ui/alert";
import {
  Eye,
  EyeOff,
  Loader2,
  Lock,
  Mail,
  CheckCircle2,
  AlertCircle,
  Trash2,
} from "lucide-react";
import { accountService, type Account } from "@/services";
import { PasswordStrengthIndicator } from "@/components/common/PasswordStrengthIndicator";
import { useToast } from "@/hooks/useToast";
import { useAuth } from "@/hooks/useAuth";
import { ROUTES } from "@/config/routes";

const TIMEZONES = [
  "UTC",
  "America/New_York",
  "America/Chicago",
  "America/Denver",
  "America/Los_Angeles",
  "Europe/London",
  "Europe/Paris",
  "Europe/Berlin",
  "Asia/Tokyo",
  "Asia/Shanghai",
  "Asia/Hong_Kong",
  "Asia/Singapore",
  "Asia/Dubai",
  "Asia/Kolkata",
  "Australia/Sydney",
  "Australia/Melbourne",
  "Pacific/Auckland",
  "America/Toronto",
  "America/Mexico_City",
  "America/Argentina/Buenos_Aires",
  "America/Sao_Paulo",
  "Africa/Cairo",
  "Africa/Johannesburg",
];

export function AccountPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { logout } = useAuth();
  const toast = useToast();

  // Account info state
  const [account, setAccount] = useState<Account | null>(null);
  const [isLoadingAccount, setIsLoadingAccount] = useState(true);
  const [accountError, setAccountError] = useState<string | null>(null);

  // Password change state
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isChangingPassword, setIsChangingPassword] = useState(false);
  const [passwordError, setPasswordError] = useState<string | null>(null);
  const [passwordSuccess, setPasswordSuccess] = useState(false);

  // Timezone state
  const [selectedTimezone, setSelectedTimezone] = useState("");
  const [isChangingTimezone, setIsChangingTimezone] = useState(false);
  const [timezoneError, setTimezoneError] = useState<string | null>(null);

  // Edit username state
  const [editingUsername, setEditingUsername] = useState(false);
  const [newUsername, setNewUsername] = useState("");
  const [isUpdatingUsername, setIsUpdatingUsername] = useState(false);
  const [usernameError, setUsernameError] = useState<string | null>(null);

  // Edit email state
  const [editingEmail, setEditingEmail] = useState(false);
  const [newEmail, setNewEmail] = useState("");
  const [isUpdatingEmail, setIsUpdatingEmail] = useState(false);
  const [emailError, setEmailError] = useState<string | null>(null);

  // Delete account state
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [deleteConfirmText, setDeleteConfirmText] = useState("");
  const [isDeletingAccount, setIsDeletingAccount] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  // Fetch account info
  useEffect(() => {
    const fetchAccount = async () => {
      try {
        setIsLoadingAccount(true);
        setAccountError(null);
        const accountData = await accountService.getAccount();
        if (accountData) {
          setAccount(accountData);
          setSelectedTimezone(accountData.timezone || "UTC");
        } else {
          setAccountError(t("account.loadingFailed"));
        }
      } catch (err) {
        const errorMsg =
          err instanceof Error ? err.message : "Failed to load account";
        setAccountError(errorMsg);
        console.error("Failed to fetch account:", err);
      } finally {
        setIsLoadingAccount(false);
      }
    };

    fetchAccount();
  }, [t]);

  // Validate password form
  const validatePasswordForm = (): boolean => {
    setPasswordError(null);

    if (!currentPassword.trim()) {
      setPasswordError(t("account.currentPasswordRequired"));
      return false;
    }

    if (!newPassword.trim()) {
      setPasswordError(t("account.newPasswordRequired"));
      return false;
    }

    if (newPassword.length < 8) {
      setPasswordError(t("account.passwordMinLength"));
      return false;
    }

    if (newPassword !== confirmPassword) {
      setPasswordError(t("account.passwordsDontMatch"));
      return false;
    }

    if (newPassword === currentPassword) {
      setPasswordError(t("account.passwordSameAsOld"));
      return false;
    }

    return true;
  };

  // Handle password change
  const handleChangePassword = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();

    if (!validatePasswordForm()) {
      return;
    }

    try {
      setIsChangingPassword(true);
      setPasswordError(null);
      setPasswordSuccess(false);

      await accountService.updateAppIdentityPassword(
        currentPassword,
        newPassword
      );

      setPasswordSuccess(true);
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");

      toast.success(t("account.passwordChangedSuccess"), {
        description: t("account.passwordChangedSuccessDesc"),
      });

      // Reset success message after 5 seconds
      setTimeout(() => setPasswordSuccess(false), 5000);
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : t("account.passwordChangeFailed");
      setPasswordError(errorMsg);
      toast.error(t("account.error"), {
        description: errorMsg,
      });
    } finally {
      setIsChangingPassword(false);
    }
  };

  // Handle timezone change
  const handleTimezoneChange = async (timezone: string) => {
    try {
      setIsChangingTimezone(true);
      setTimezoneError(null);
      setSelectedTimezone(timezone);

      await accountService.updateTimezone(timezone);

      toast.success(t("account.timezoneUpdated"), {
        description: t("account.timezoneUpdatedDesc"),
      });
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : t("account.timezoneUpdateFailed");
      setTimezoneError(errorMsg);
      setSelectedTimezone(account?.timezone || "UTC");
      toast.error(t("account.error"), {
        description: errorMsg,
      });
    } finally {
      setIsChangingTimezone(false);
    }
  };

  // Handle username update
  const handleUpdateUsername = async () => {
    if (!newUsername.trim()) {
      setUsernameError("Username is required");
      return;
    }

    try {
      setIsUpdatingUsername(true);
      setUsernameError(null);

      await accountService.updateAppIdentityUsername(newUsername);

      toast.success(t("account.usernameUpdated"), {
        description: t("account.usernameUpdatedDesc"),
      });

      if (account && appIdentity) {
        appIdentity.username = newUsername;
        setAccount({ ...account });
      }

      setEditingUsername(false);
      setNewUsername("");
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : t("account.usernameUpdateFailed");
      setUsernameError(errorMsg);
      toast.error(t("account.error"), {
        description: errorMsg,
      });
    } finally {
      setIsUpdatingUsername(false);
    }
  };

  // Handle email update
  const handleUpdateEmail = async () => {
    if (!newEmail.trim()) {
      setEmailError("Email is required");
      return;
    }

    // Basic email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(newEmail)) {
      setEmailError("Please enter a valid email address");
      return;
    }

    try {
      setIsUpdatingEmail(true);
      setEmailError(null);

      await accountService.updateAppIdentityEmail(newEmail);

      toast.success(t("account.emailUpdated"), {
        description: t("account.emailUpdatedDesc"),
      });

      if (account && appIdentity) {
        appIdentity.external_id = newEmail;
        setAccount({ ...account });
      }

      setEditingEmail(false);
      setNewEmail("");
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : t("account.emailUpdateFailed");
      setEmailError(errorMsg);
      toast.error(t("account.error"), {
        description: errorMsg,
      });
    } finally {
      setIsUpdatingEmail(false);
    }
  };

  // Handle delete account
  const handleDeleteAccount = async () => {
    if (deleteConfirmText !== "DELETE") {
      setDeleteError(t("account.deleteConfirmationMismatch"));
      return;
    }

    try {
      setIsDeletingAccount(true);
      setDeleteError(null);

      await accountService.deleteAccount();

      toast.success(t("account.accountDeleted"), {
        description: t("account.accountDeletedDesc"),
      });

      // Logout and redirect to home
      logout();
      setTimeout(() => {
        navigate(ROUTES.VITRINE.path);
      }, 1000);
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : t("account.deleteAccountFailed");
      setDeleteError(errorMsg);
      toast.error(t("account.error"), {
        description: errorMsg,
      });
    } finally {
      setIsDeletingAccount(false);
    }
  };

  // Get app identity
  const appIdentity = account?.identities?.find((id) => id.provider === "app");
  const discordIdentity = account?.identities?.find(
    (id) => id.provider === "discord"
  );

  return (
    <div className="min-h-screen bg-background-main dark:bg-background-main">
      <Header />

      {/* Main Content */}
      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Page Header */}
        <div className="mb-12">
          <h1 className="text-4xl font-bold text-foreground mb-2">
            {t("account.title")}
          </h1>
          <p className="text-muted-foreground">{t("account.subtitle")}</p>
        </div>

        {/* Error State */}
        {accountError && (
          <Alert className="mb-6 border-red-500/50 bg-red-500/10">
            <AlertCircle className="h-4 w-4 text-red-600 dark:text-red-400" />
            <div className="ml-3">
              <p className="text-red-600 dark:text-red-400">{accountError}</p>
            </div>
          </Alert>
        )}

        {/* Loading State */}
        {isLoadingAccount ? (
          <Card className="border-border bg-card/95 backdrop-blur">
            <CardContent className="pt-6 flex items-center justify-center py-12">
              <Loader2 className="w-8 h-8 text-muted-foreground animate-spin" />
              <p className="text-muted-foreground ml-3">
                {t("common.loading")}
              </p>
            </CardContent>
          </Card>
        ) : account ? (
          <div className="space-y-6">
            {/* Account Information Section */}
            <Card className="border-border bg-card/95 backdrop-blur">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Mail className="w-5 h-5 text-accent" />
                  {t("account.accountInfo")}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-6">
                {/* Username from app identity */}
                <div>
                  <Label className="text-muted-foreground">
                    {t("account.username")}
                  </Label>
                  {!editingUsername ? (
                    <div
                      className="mt-2 flex items-center gap-2 p-3 bg-secondary/30 rounded-md border border-border group hover:bg-secondary/50 transition-colors cursor-pointer"
                      onClick={() => {
                        setNewUsername(
                          appIdentity?.username || account.username
                        );
                        setEditingUsername(true);
                        setUsernameError(null);
                      }}
                    >
                      <p className="text-foreground flex-1">
                        {appIdentity?.username || account.username}
                      </p>
                      <Button
                        size="sm"
                        variant="ghost"
                        className="opacity-0 group-hover:opacity-100 transition-opacity h-7 w-7 p-0"
                      >
                        <Mail className="w-4 h-4 text-accent" />
                      </Button>
                    </div>
                  ) : (
                    <div className="mt-2 space-y-2">
                      {usernameError && (
                        <Alert className="border-red-500/50 bg-red-500/10">
                          <AlertCircle className="h-4 w-4 text-red-600 dark:text-red-400" />
                          <div className="ml-3">
                            <p className="text-red-600 dark:text-red-400">
                              {usernameError}
                            </p>
                          </div>
                        </Alert>
                      )}
                      <div className="flex gap-2 items-center">
                        <Input
                          value={newUsername}
                          onChange={(e) => setNewUsername(e.target.value)}
                          disabled={isUpdatingUsername}
                          placeholder="Enter new username"
                          className="flex-1"
                          autoFocus
                        />
                        <Button
                          onClick={handleUpdateUsername}
                          disabled={isUpdatingUsername}
                          size="sm"
                          className="bg-accent hover:bg-accent/90 h-9 px-3"
                        >
                          {isUpdatingUsername ? (
                            <Loader2 className="w-4 h-4 animate-spin" />
                          ) : (
                            <CheckCircle2 className="w-4 h-4" />
                          )}
                        </Button>
                        <Button
                          variant="ghost"
                          onClick={() => {
                            setEditingUsername(false);
                            setNewUsername("");
                            setUsernameError(null);
                          }}
                          disabled={isUpdatingUsername}
                          size="sm"
                          className="h-9 px-3"
                        >
                          ✕
                        </Button>
                      </div>
                    </div>
                  )}
                </div>

                {/* Email from app identity */}
                <div>
                  <Label className="text-muted-foreground">
                    {t("account.email")}
                  </Label>
                  {!editingEmail ? (
                    <div
                      className="mt-2 flex items-center gap-2 p-3 bg-secondary/30 rounded-md border border-border group hover:bg-secondary/50 transition-colors cursor-pointer"
                      onClick={() => {
                        setNewEmail(appIdentity?.external_id || account.email);
                        setEditingEmail(true);
                        setEmailError(null);
                      }}
                    >
                      <Mail className="w-4 h-4 text-muted-foreground flex-shrink-0" />
                      <p className="text-foreground flex-1">
                        {appIdentity?.external_id || account.email}
                      </p>
                      <Button
                        size="sm"
                        variant="ghost"
                        className="opacity-0 group-hover:opacity-100 transition-opacity h-7 w-7 p-0"
                      >
                        <Mail className="w-4 h-4 text-accent" />
                      </Button>
                    </div>
                  ) : (
                    <div className="mt-2 space-y-2">
                      {emailError && (
                        <Alert className="border-red-500/50 bg-red-500/10">
                          <AlertCircle className="h-4 w-4 text-red-600 dark:text-red-400" />
                          <div className="ml-3">
                            <p className="text-red-600 dark:text-red-400">
                              {emailError}
                            </p>
                          </div>
                        </Alert>
                      )}
                      <div className="flex gap-2 items-center">
                        <Input
                          type="email"
                          value={newEmail}
                          onChange={(e) => setNewEmail(e.target.value)}
                          disabled={isUpdatingEmail}
                          placeholder="Enter new email"
                          className="flex-1"
                          autoFocus
                        />
                        <Button
                          onClick={handleUpdateEmail}
                          disabled={isUpdatingEmail}
                          size="sm"
                          className="bg-accent hover:bg-accent/90 h-9 px-3"
                        >
                          {isUpdatingEmail ? (
                            <Loader2 className="w-4 h-4 animate-spin" />
                          ) : (
                            <CheckCircle2 className="w-4 h-4" />
                          )}
                        </Button>
                        <Button
                          variant="ghost"
                          onClick={() => {
                            setEditingEmail(false);
                            setNewEmail("");
                            setEmailError(null);
                          }}
                          disabled={isUpdatingEmail}
                          size="sm"
                          className="h-9 px-3"
                        >
                          ✕
                        </Button>
                      </div>
                    </div>
                  )}
                </div>

                {/* Account Created */}
                <div>
                  <Label className="text-muted-foreground">
                    {t("account.created")}
                  </Label>
                  <div className="mt-2 p-3 bg-secondary/30 rounded-md border border-border">
                    <p className="text-foreground font-medium">
                      {new Date(account.created_at).toLocaleDateString([], {
                        year: "numeric",
                        month: "long",
                        day: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Timezone Selector Section */}
            <Card className="border-border bg-card/95 backdrop-blur">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <CheckCircle2 className="w-5 h-5 text-accent" />
                  {t("account.timezoneSettings")}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {timezoneError && (
                  <Alert className="border-red-500/50 bg-red-500/10">
                    <AlertCircle className="h-4 w-4 text-red-600 dark:text-red-400" />
                    <div className="ml-3">
                      <p className="text-red-600 dark:text-red-400">
                        {timezoneError}
                      </p>
                    </div>
                  </Alert>
                )}

                <div>
                  <Label htmlFor="timezone-select" className="text-foreground">
                    {t("account.selectTimezone")}
                  </Label>
                  <select
                    id="timezone-select"
                    value={selectedTimezone}
                    onChange={(e) => handleTimezoneChange(e.target.value)}
                    disabled={isChangingTimezone}
                    className="mt-2 w-full px-3 py-2 bg-secondary/30 border border-border rounded-md text-foreground focus:outline-none focus:ring-2 focus:ring-accent disabled:opacity-50"
                  >
                    {TIMEZONES.map((tz) => (
                      <option key={tz} value={tz}>
                        {tz}
                      </option>
                    ))}
                  </select>
                </div>

                {isChangingTimezone && (
                  <div className="flex items-center gap-2 text-muted-foreground text-sm">
                    <Loader2 className="w-4 h-4 animate-spin" />
                    {t("account.updatingTimezone")}
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Identities Section */}
            <Card className="border-border bg-card/95 backdrop-blur">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <CheckCircle2 className="w-5 h-5 text-accent" />
                  {t("account.connectedIdentities")}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* App Identity */}
                <div className="p-4 bg-secondary/20 rounded-lg border border-border flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 sm:gap-3">
                  <div className="flex items-center gap-3 flex-1 min-w-0">
                    <div className="flex flex-shrink-0 h-10 w-10 items-center justify-center rounded-lg bg-accent/20">
                      <Lock className="w-5 h-5 flex-shrink-0 text-accent" />
                    </div>
                    <div>
                      <p className="font-semibold text-foreground">
                        {t("account.appIdentity")}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {t("account.appIdentityDesc")}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center justify-end flex-shrink-0">
                    {appIdentity ? (
                      <span className="inline-flex items-center px-2 py-0.5 sm:px-3 sm:py-1 rounded-full text-xs font-medium bg-green-500/20 text-green-600 dark:text-green-400 border border-green-500/30 whitespace-nowrap">
                        {t("account.connected")}
                      </span>
                    ) : (
                      <span className="inline-flex items-center px-2 sm:px-3 py-0.5 sm:py-1 rounded-full text-xs font-medium bg-gray-500/20 text-gray-600 dark:text-gray-400 border border-gray-500/30">
                        {t("account.notConnected")}
                      </span>
                    )}
                  </div>
                </div>

                {/* Discord Identity */}
                <div className="p-4 bg-secondary/20 rounded-lg border border-border flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 sm:gap-3">
                  <div className="flex items-center gap-3 flex-1 min-w-0">
                    <div className="flex flex-shrink-0 h-10 w-10 items-center justify-center rounded-lg bg-indigo-500/20">
                      <svg
                        className="w-5 h-5 flex-shrink-0 text-indigo-500"
                        fill="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path d="M20.317 4.37a19.791 19.791 0 00-4.885-1.515.074.074 0 00-.079.037c-.21.375-.444.864-.607 1.25a18.27 18.27 0 00-5.487 0c-.163-.386-.397-.875-.61-1.25a.077.077 0 00-.079-.037A19.736 19.736 0 003.677 4.37a.07.07 0 00-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 00.031.056 19.9 19.9 0 005.993 3.03.078.078 0 00.084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 00-.042-.106 13.107 13.107 0 01-1.872-.892.077.077 0 00-.008-.128 10.713 10.713 0 00.372-.294.075.075 0 00.03-.066c.329-.246.648-.5.954-.76a.07.07 0 00.076-.01 13.697 13.697 0 0011.086 0 .07.07 0 00.076.009c.305.26.625.514.954.759a.077.077 0 00.03.067c.12.088.246.177.371.294a.077.077 0 00-.006.127 13.227 13.227 0 01-1.873.892.076.076 0 00-.041.107c.352.699.764 1.364 1.225 1.994a.076.076 0 00.084.028 19.963 19.963 0 006.002-3.03.077.077 0 00.032-.054c.5-4.817-.838-9.033-3.55-12.765a.061.061 0 00-.031-.03zM8.02 15.33c-1.183 0-2.157-.969-2.157-2.156 0-1.193.974-2.157 2.157-2.157 1.193 0 2.156.964 2.156 2.157 0 1.187-.963 2.156-2.156 2.156zm7.975 0c-1.183 0-2.157-.969-2.157-2.156 0-1.193.974-2.157 2.157-2.157 1.193 0 2.157.964 2.157 2.157 0 1.187-.964 2.156-2.157 2.156z" />
                      </svg>
                    </div>
                    <div>
                      <p className="font-semibold text-foreground">
                        {t("account.discordIdentity")}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {t("account.discordIdentityDesc")}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center justify-end flex-shrink-0">
                    {discordIdentity ? (
                      <span className="inline-flex items-center px-2 py-0.5 sm:px-3 sm:py-1 rounded-full text-xs font-medium bg-green-500/20 text-green-600 dark:text-green-400 border border-green-500/30 whitespace-nowrap">
                        {t("account.connected")}
                      </span>
                    ) : (
                      <Button
                        type="button"
                        onClick={() => {
                          const clientId = import.meta.env
                            .VITE_DISCORD_CLIENT_ID;
                          const redirectUri = import.meta.env
                            .VITE_DISCORD_REDIRECT_URI;

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
                        variant="outline"
                        size="sm"
                        className="border-indigo-500/30 text-indigo-600 dark:text-indigo-400 hover:bg-indigo-500/10"
                      >
                        <svg
                          className="w-4 h-4 mr-2"
                          fill="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path d="M20.317 4.37a19.791 19.791 0 00-4.885-1.515.074.074 0 00-.079.037c-.21.375-.444.864-.607 1.25a18.27 18.27 0 00-5.487 0c-.163-.386-.397-.875-.61-1.25a.077.077 0 00-.079-.037A19.736 19.736 0 003.677 4.37a.07.07 0 00-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 00.031.056 19.9 19.9 0 005.993 3.03.078.078 0 00.084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 00-.042-.106 13.107 13.107 0 01-1.872-.892.077.077 0 00-.008-.128 10.713 10.713 0 00.372-.294.075.075 0 00.03-.066c.329-.246.648-.5.954-.76a.07.07 0 00.076-.01 13.697 13.697 0 0011.086 0 .07.07 0 00.076.009c.305.26.625.514.954.759a.077.077 0 00.03.067c.12.088.246.177.371.294a.077.077 0 00-.006.127 13.227 13.227 0 01-1.873.892.076.076 0 00-.041.107c.352.699.764 1.364 1.225 1.994a.076.076 0 00.084.028 19.963 19.963 0 006.002-3.03.077.077 0 00.032-.054c.5-4.817-.838-9.033-3.55-12.765a.061.061 0 00-.031-.03zM8.02 15.33c-1.183 0-2.157-.969-2.157-2.156 0-1.193.974-2.157 2.157-2.157 1.193 0 2.156.964 2.156 2.157 0 1.187-.963 2.156-2.156 2.156zm7.975 0c-1.183 0-2.157-.969-2.157-2.156 0-1.193.974-2.157 2.157-2.157 1.193 0 2.157.964 2.157 2.157 0 1.187-.964 2.156-2.157 2.156z" />
                        </svg>
                        {t("account.connectDiscord")}
                      </Button>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Change Password Section - Only show if app identity exists */}
            {appIdentity && (
              <Card className="border-border bg-card/95 backdrop-blur">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Lock className="w-5 h-5 text-accent" />
                    {t("account.changePassword")}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <form onSubmit={handleChangePassword} className="space-y-6">
                    {/* Success Message */}
                    {passwordSuccess && (
                      <Alert className="border-green-500/50 bg-green-500/10">
                        <CheckCircle2 className="h-4 w-4 text-green-600 dark:text-green-400" />
                        <div className="ml-3">
                          <p className="text-green-600 dark:text-green-400">
                            {t("account.passwordChangedSuccess")}
                          </p>
                        </div>
                      </Alert>
                    )}

                    {/* Error Message */}
                    {passwordError && (
                      <Alert className="border-red-500/50 bg-red-500/10">
                        <AlertCircle className="h-4 w-4 text-red-600 dark:text-red-400" />
                        <div className="ml-3">
                          <p className="text-red-600 dark:text-red-400">
                            {passwordError}
                          </p>
                        </div>
                      </Alert>
                    )}

                    {/* Current Password */}
                    <div>
                      <Label
                        htmlFor="current-password"
                        className="text-foreground"
                      >
                        {t("account.currentPassword")}
                      </Label>
                      <div className="mt-2 relative">
                        <Input
                          id="current-password"
                          type={showCurrentPassword ? "text" : "password"}
                          placeholder={t("account.currentPasswordPlaceholder")}
                          value={currentPassword}
                          onChange={(e) => setCurrentPassword(e.target.value)}
                          disabled={isChangingPassword}
                          className="pr-10"
                        />
                        <button
                          type="button"
                          onClick={() =>
                            setShowCurrentPassword(!showCurrentPassword)
                          }
                          className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                        >
                          {showCurrentPassword ? (
                            <EyeOff className="w-4 h-4" />
                          ) : (
                            <Eye className="w-4 h-4" />
                          )}
                        </button>
                      </div>
                    </div>

                    {/* New Password */}
                    <div>
                      <Label htmlFor="new-password" className="text-foreground">
                        {t("account.newPassword")}
                      </Label>
                      <div className="mt-2 relative">
                        <Input
                          id="new-password"
                          type={showNewPassword ? "text" : "password"}
                          placeholder={t("account.newPasswordPlaceholder")}
                          value={newPassword}
                          onChange={(e) => setNewPassword(e.target.value)}
                          disabled={isChangingPassword}
                          className="pr-10"
                        />
                        <button
                          type="button"
                          onClick={() => setShowNewPassword(!showNewPassword)}
                          className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                        >
                          {showNewPassword ? (
                            <EyeOff className="w-4 h-4" />
                          ) : (
                            <Eye className="w-4 h-4" />
                          )}
                        </button>
                      </div>
                      {newPassword && (
                        <PasswordStrengthIndicator password={newPassword} />
                      )}
                    </div>

                    {/* Confirm Password */}
                    <div>
                      <Label
                        htmlFor="confirm-password"
                        className="text-foreground"
                      >
                        {t("account.confirmPassword")}
                      </Label>
                      <div className="mt-2 relative">
                        <Input
                          id="confirm-password"
                          type={showConfirmPassword ? "text" : "password"}
                          placeholder={t("account.confirmPasswordPlaceholder")}
                          value={confirmPassword}
                          onChange={(e) => setConfirmPassword(e.target.value)}
                          disabled={isChangingPassword}
                          className="pr-10"
                        />
                        <button
                          type="button"
                          onClick={() =>
                            setShowConfirmPassword(!showConfirmPassword)
                          }
                          className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                        >
                          {showConfirmPassword ? (
                            <EyeOff className="w-4 h-4" />
                          ) : (
                            <Eye className="w-4 h-4" />
                          )}
                        </button>
                      </div>
                      {confirmPassword && newPassword === confirmPassword && (
                        <p className="text-xs text-green-600 dark:text-green-400 mt-1">
                          {t("account.passwordsMatch")}
                        </p>
                      )}
                    </div>

                    {/* Submit Button */}
                    <div className="flex gap-3 pt-4">
                      <Button
                        type="submit"
                        disabled={isChangingPassword}
                        className="bg-accent hover:bg-accent/90 text-accent-foreground font-semibold gap-2"
                      >
                        {isChangingPassword ? (
                          <>
                            <Loader2 className="w-4 h-4 animate-spin" />
                            {t("account.updatingPassword")}
                          </>
                        ) : (
                          <>
                            <Lock className="w-4 h-4" />
                            {t("account.updatePassword")}
                          </>
                        )}
                      </Button>
                      <Button
                        type="button"
                        variant="ghost"
                        onClick={() => {
                          setCurrentPassword("");
                          setNewPassword("");
                          setConfirmPassword("");
                          setPasswordError(null);
                          setPasswordSuccess(false);
                        }}
                        disabled={isChangingPassword}
                        className="text-muted-foreground hover:text-foreground"
                      >
                        {t("common.cancel")}
                      </Button>
                    </div>
                  </form>
                </CardContent>
              </Card>
            )}

            {/* Delete Account Section - Danger Zone */}
            <Card className="border-red-500/30 bg-red-500/5 backdrop-blur">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-red-600 dark:text-red-400">
                  <Trash2 className="w-5 h-5" />
                  {t("account.dangerZone")}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {!showDeleteConfirm ? (
                  <>
                    <div className="flex items-start gap-3 p-4 bg-red-500/10 rounded-lg border border-red-500/20">
                      <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" />
                      <div>
                        <p className="text-red-600 dark:text-red-400 font-medium">
                          {t("account.deleteAccountWarning")}
                        </p>
                        <p className="text-red-600/80 dark:text-red-400/80 text-sm mt-1">
                          {t("account.deleteAccountWarningDesc")}
                        </p>
                      </div>
                    </div>
                    <Button
                      onClick={() => setShowDeleteConfirm(true)}
                      className="w-full bg-red-600 hover:bg-red-700 text-white font-semibold gap-2"
                    >
                      <Trash2 className="w-4 h-4" />
                      {t("account.deleteAccountButton")}
                    </Button>
                  </>
                ) : (
                  <>
                    {deleteError && (
                      <Alert className="border-red-500/50 bg-red-500/10">
                        <AlertCircle className="h-4 w-4 text-red-600 dark:text-red-400" />
                        <div className="ml-3">
                          <p className="text-red-600 dark:text-red-400">
                            {deleteError}
                          </p>
                        </div>
                      </Alert>
                    )}

                    <div className="space-y-4 p-4 bg-red-500/10 rounded-lg border border-red-500/20">
                      <div>
                        <p className="text-foreground font-medium mb-2">
                          {t("account.deleteConfirmationPrompt")}
                        </p>
                        <p className="text-muted-foreground text-sm mb-3">
                          {t("account.deleteConfirmationDesc")}
                        </p>
                        <Input
                          type="text"
                          placeholder="DELETE"
                          value={deleteConfirmText}
                          onChange={(e) => setDeleteConfirmText(e.target.value)}
                          disabled={isDeletingAccount}
                          className="font-mono"
                        />
                      </div>

                      <div className="flex gap-3">
                        <Button
                          onClick={handleDeleteAccount}
                          disabled={
                            isDeletingAccount || deleteConfirmText !== "DELETE"
                          }
                          className="flex-1 bg-red-600 hover:bg-red-700 text-white font-semibold gap-2"
                        >
                          {isDeletingAccount ? (
                            <>
                              <Loader2 className="w-4 h-4 animate-spin" />
                              {t("account.deletingAccount")}
                            </>
                          ) : (
                            <>
                              <Trash2 className="w-4 h-4" />
                              {t("account.confirmDelete")}
                            </>
                          )}
                        </Button>
                        <Button
                          type="button"
                          variant="ghost"
                          onClick={() => {
                            setShowDeleteConfirm(false);
                            setDeleteConfirmText("");
                            setDeleteError(null);
                          }}
                          disabled={isDeletingAccount}
                          className="text-muted-foreground hover:text-foreground"
                        >
                          {t("common.cancel")}
                        </Button>
                      </div>
                    </div>
                  </>
                )}
              </CardContent>
            </Card>
          </div>
        ) : null}
      </main>

      {/* Footer */}
      <Footer />
    </div>
  );
}
