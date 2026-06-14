import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { DiscordSetupSection } from "@/components/DiscordSetupSection";
import { authService } from "@/services";

interface SetupData {
  status: string;
  message: string;
  account_id: string;
  discord_email: string;
  discord_username: string;
  needs_setup: boolean;
}

interface MergeData {
  discord_username: string;
  merge_token: string;
}

export function OAuthCallbackPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { t } = useTranslation();
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [setupData, setSetupData] = useState<SetupData | null>(null);
  const [isCompletingSetup, setIsCompletingSetup] = useState(false);
  const [mergeData, setMergeData] = useState<MergeData | null>(null);
  const [isMerging, setIsMerging] = useState(false);

  useEffect(() => {
    // Prevent duplicate requests in development (React StrictMode calls effects twice)
    let isMounted = true;
    let hasProcessed = false;

    const handleCallback = async () => {
      const code = searchParams.get("code");
      const errorParam = searchParams.get("error");
      const state = searchParams.get("state");

      // Only process once
      if (hasProcessed) {
        return;
      }
      hasProcessed = true;

      if (errorParam) {
        if (isMounted) {
          setError(`Discord authentication failed: ${errorParam}`);
          setIsLoading(false);
        }
        return;
      }

      if (!code) {
        if (isMounted) {
          setError("No authorization code received from Discord");
          setIsLoading(false);
        }
        return;
      }

      const apiUrl =
        import.meta.env.VITE_API_URL || "https://api.chronosrmd.com";

      // "link" flow: an already-authenticated user is connecting Discord to
      // their existing account. Call the authenticated link endpoint instead
      // of the login/signup callback.
      if (state === "link") {
        try {
          const token = localStorage.getItem("auth_token");
          const response = await fetch(
            `${apiUrl}/api/account/identity/discord/link`,
            {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                ...(token ? { Authorization: `Bearer ${token}` } : {}),
              },
              credentials: "include",
              body: JSON.stringify({ code }),
            },
          );

          if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || "Failed to link Discord account");
          }

          const linkData = await response.json().catch(() => ({}));

          if (linkData.merge_required) {
            if (isMounted) {
              setMergeData({
                discord_username: linkData.discord_username || "",
                merge_token: linkData.merge_token || "",
              });
              setIsLoading(false);
            }
            return;
          }

          if (isMounted) {
            navigate("/account", { replace: true });
          }
        } catch (err) {
          const errorMessage =
            err instanceof Error
              ? err.message
              : "Failed to link Discord account";
          if (isMounted) {
            setError(errorMessage);
            setIsLoading(false);
          }
        }
        return;
      }

      try {
        // Send code to backend
        const response = await fetch(`${apiUrl}/api/auth/discord/callback`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({ code }),
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || "Authentication failed");
        }

        const data = await response.json();

        // Check if setup is required (user has only Discord identity, no app identity)
        if (data.status === "setup_required" && data.needs_setup) {
          if (isMounted) {
            setSetupData(data);
            setIsLoading(false);
          }
          return;
        }

        // Normal login flow
        localStorage.setItem("auth_token", data.token);

        // Calculate and store expires_at (30 days from now)
        const expiresAt = new Date();
        expiresAt.setDate(expiresAt.getDate() + 30);
        localStorage.setItem("token_expires_at", expiresAt.toISOString());

        const userData = {
          user_id: data.id,
          email: data.email,
          username: data.username,
          expires_at: expiresAt.toISOString(),
        };

        localStorage.setItem("user_data", JSON.stringify(userData));

        // Dispatch custom event to notify auth context
        window.dispatchEvent(new Event("auth-updated"));

        // Small delay to ensure localStorage is persisted
        await new Promise((resolve) => setTimeout(resolve, 100));

        // Navigate to welcome page
        if (isMounted) {
          navigate("/welcome", { replace: true });
        }
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Authentication failed";
        if (isMounted) {
          setError(errorMessage);
          setIsLoading(false);
        }
      }
    };

    handleCallback();

    return () => {
      isMounted = false;
    };
  }, [searchParams, navigate]);

  const handleSetupComplete = async (setupForm: {
    email: string;
    username: string;
    password: string;
    timezone: string;
  }) => {
    if (!setupData) {
      return;
    }

    try {
      setIsCompletingSetup(true);
      const apiUrl =
        import.meta.env.VITE_API_URL || "https://api.chronosrmd.com";

      const response = await fetch(`${apiUrl}/api/auth/discord/setup`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          account_id: setupData.account_id,
          email: setupForm.email,
          username: setupForm.username,
          password: setupForm.password,
          timezone: setupForm.timezone,
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Setup failed");
      }

      const data = await response.json();

      // Use the expires_at from backend response if available, otherwise calculate 30 days
      let expiresAtStr = data.expiresAt;
      if (!expiresAtStr) {
        const expiresAt = new Date();
        expiresAt.setDate(expiresAt.getDate() + 30);
        expiresAtStr = expiresAt.toISOString();
      }

      const userData = {
        user_id: data.id,
        email: data.email,
        username: data.username,
        expires_at: expiresAtStr,
      };

      // Use the ApiClient's method to set authentication (handles both localStorage and internal state)
      authService.setAuthentication(data.token, expiresAtStr, userData);

      // Dispatch custom event to notify auth context
      window.dispatchEvent(new Event("auth-updated"));

      // Small delay to ensure localStorage is persisted
      await new Promise((resolve) => setTimeout(resolve, 100));

      // Navigate to welcome page
      navigate("/welcome", { replace: true });
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Setup failed";
      setError(errorMessage);
      setIsCompletingSetup(false);
    }
  };

  const handleConfirmMerge = async () => {
    if (!mergeData) return;
    const apiUrl = import.meta.env.VITE_API_URL || "https://api.chronosrmd.com";
    try {
      setIsMerging(true);
      const token = localStorage.getItem("auth_token");
      const response = await fetch(`${apiUrl}/api/account/merge`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        credentials: "include",
        body: JSON.stringify({ merge_token: mergeData.merge_token }),
      });
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || t("account.merge.failed"));
      }
      navigate("/account", { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : t("account.merge.failed"));
      setMergeData(null);
    } finally {
      setIsMerging(false);
    }
  };

  // Show merge confirmation dialog
  if (mergeData && !error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background-main to-background-secondary flex items-center justify-center p-4">
        <div className="bg-card rounded-lg shadow-lg p-8 max-w-md w-full space-y-4">
          <h2 className="text-xl font-bold text-foreground">{t("account.merge.title")}</h2>
          <p className="text-muted-foreground text-sm">
            {t("account.merge.description", { username: mergeData.discord_username })}
          </p>
          <div className="flex gap-3 pt-2">
            <button
              onClick={handleConfirmMerge}
              disabled={isMerging}
              className="flex-1 px-4 py-2 bg-accent hover:bg-accent/90 disabled:opacity-50 text-accent-foreground rounded font-medium transition-colors"
            >
              {isMerging ? t("account.merge.merging") : t("account.merge.confirm")}
            </button>
            <button
              onClick={() => navigate("/account", { replace: true })}
              disabled={isMerging}
              className="flex-1 px-4 py-2 border border-border hover:bg-secondary/30 text-foreground rounded font-medium transition-colors"
            >
              {t("account.merge.cancel")}
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background-main to-background-secondary flex items-center justify-center">
        <div className="text-center space-y-4">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-accent"></div>
          <p className="text-foreground text-lg">
            Authenticating with Discord...
          </p>
        </div>
      </div>
    );
  }

  // Show setup form if setup is required
  if (setupData && !error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background-main to-background-secondary flex items-center justify-center p-4">
        <div className="bg-card rounded-lg shadow-lg p-8 max-w-md w-full">
          <DiscordSetupSection
            discordEmail={setupData.discord_email}
            discordUsername={setupData.discord_username}
            onSetupComplete={handleSetupComplete}
            isLoading={isCompletingSetup}
            error={error}
          />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background-main to-background-secondary flex items-center justify-center p-4">
        <div className="bg-card rounded-lg shadow-lg p-8 max-w-md w-full text-center space-y-4">
          <div className="text-red-600 text-xl">⚠️ Authentication Error</div>
          <p className="text-red-600 text-sm">{error}</p>
          <button
            onClick={() => navigate("/login", { replace: true })}
            className="px-4 py-2 bg-accent hover:bg-accent/90 text-accent-foreground rounded font-medium transition-colors w-full"
          >
            Back to Login
          </button>
        </div>
      </div>
    );
  }

  return null;
}
