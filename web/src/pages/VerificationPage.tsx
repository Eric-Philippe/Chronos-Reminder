import { useState, useEffect, useCallback, useRef } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Loader2, CheckCircle2, XCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/useToast";
import { httpClient } from "@/services/http";
import type { VerifyEmailResponse } from "@/services/types";

export function VerificationPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { success, error } = useToast();
  const verificationAttemptedRef = useRef(false);

  const [email, setEmail] = useState("");
  const [code, setCode] = useState("");
  const [status, setStatus] = useState<
    "pending" | "loading" | "success" | "error"
  >("pending");
  const [errorMessage, setErrorMessage] = useState("");
  const [isManualVerification, setIsManualVerification] = useState(false);

  const verifyEmailWithCode = useCallback(
    async (verifyEmail: string, verifyCode: string) => {
      try {
        setStatus("loading");
        const response = await httpClient.post<VerifyEmailResponse>(
          "/api/auth/verify",
          {
            email: verifyEmail,
            code: verifyCode,
          }
        );

        // Handle the response based on the actual return type
        const verifyResponse = (response.data ||
          response) as VerifyEmailResponse;

        // Set authentication
        if (verifyResponse.token && verifyResponse.id) {
          httpClient.setToken(
            verifyResponse.token,
            new Date(verifyResponse.expires_at)
          );

          const userData = {
            user_id: verifyResponse.id,
            email: verifyResponse.email,
            username: verifyResponse.username,
            expires_at: verifyResponse.expires_at,
          };

          localStorage.setItem("user_data", JSON.stringify(userData));

          // Dispatch event to notify AuthContext of successful login
          window.dispatchEvent(new Event("auth-updated"));

          success(t("verification.success") as string);

          setStatus("success");

          // Redirect to home after brief delay
          setTimeout(() => {
            navigate("/welcome", { replace: true });
          }, 1500);
        }
      } catch (err) {
        const errorMsg =
          (err instanceof Error ? err.message : null) ||
          t("verification.invalidCode");
        setErrorMessage(errorMsg as string);
        setStatus("error");

        error(errorMsg as string);
      }
    },
    [t, navigate, success, error]
  );

  // Check if we have email and code from URL params
  useEffect(() => {
    const urlEmail = searchParams.get("email");
    const urlCode = searchParams.get("code");

    // Only attempt auto-verification once
    if (urlEmail && urlCode && !verificationAttemptedRef.current) {
      verificationAttemptedRef.current = true;
      setEmail(urlEmail);
      setCode(urlCode);
      verifyEmailWithCode(urlEmail, urlCode);
    } else if (urlEmail) {
      setEmail(urlEmail);
      setIsManualVerification(true);
    }
  }, [searchParams, verifyEmailWithCode]);

  const handleManualVerification = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !code) {
      setErrorMessage(t("verification.fillFields"));
      return;
    }
    await verifyEmailWithCode(email, code);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background-main to-background-secondary flex items-center justify-center p-4">
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-0 right-0 w-72 h-72 bg-accent/10 rounded-full blur-3xl dark:bg-accent/5"></div>
        <div className="absolute bottom-0 left-0 w-96 h-96 bg-accent/10 rounded-full blur-3xl dark:bg-accent/5"></div>
      </div>

      <div className="relative z-10 w-full max-w-md">
        <Card className="border-border bg-card/95 backdrop-blur">
          <CardHeader className="space-y-1 text-center">
            {status === "success" ? (
              <>
                <div className="flex justify-center mb-4">
                  <CheckCircle2 className="w-12 h-12 text-green-500" />
                </div>
                <CardTitle className="text-foreground">
                  {t("verification.success")}
                </CardTitle>
                <CardDescription>
                  {t("verification.successDesc")}
                </CardDescription>
              </>
            ) : status === "error" && !isManualVerification ? (
              <>
                <div className="flex justify-center mb-4">
                  <XCircle className="w-12 h-12 text-red-500" />
                </div>
                <CardTitle className="text-foreground">
                  {t("verification.failed")}
                </CardTitle>
                <CardDescription>{errorMessage}</CardDescription>
              </>
            ) : status === "loading" && !isManualVerification ? (
              <>
                <div className="flex justify-center mb-4">
                  <Loader2 className="w-12 h-12 text-blue-500 animate-spin" />
                </div>
                <CardTitle className="text-foreground">
                  {t("verification.verifying")}
                </CardTitle>
                <CardDescription>
                  {t("verification.verifyingDesc")}
                </CardDescription>
              </>
            ) : (
              <>
                <CardTitle className="text-foreground">
                  {t("verification.enterCode")}
                </CardTitle>
                <CardDescription>
                  {t("verification.enterCodeDesc")}
                </CardDescription>
              </>
            )}
          </CardHeader>

          <CardContent>
            {status === "success" ? (
              <div className="space-y-4">
                <p className="text-sm text-muted-foreground text-center">
                  {t("verification.redirecting")}
                </p>
              </div>
            ) : status === "error" && !isManualVerification ? (
              <Button
                onClick={() => {
                  setStatus("pending");
                  setIsManualVerification(true);
                  setErrorMessage("");
                }}
                className="w-full"
                variant="outline"
              >
                {t("verification.tryAgain")}
              </Button>
            ) : (
              <form onSubmit={handleManualVerification} className="space-y-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-foreground">
                    {t("verification.email")}
                  </label>
                  <Input
                    type="email"
                    placeholder="your@email.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    disabled={status === "loading"}
                    className="bg-secondary/50 border-border text-foreground"
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium text-foreground">
                    {t("verification.code")}
                  </label>
                  <Input
                    type="text"
                    placeholder="000000"
                    value={code}
                    onChange={(e) => setCode(e.target.value.slice(0, 6))}
                    disabled={status === "loading"}
                    maxLength={6}
                    className="bg-secondary/50 border-border text-foreground tracking-widest text-center text-lg"
                  />
                </div>

                {errorMessage && (
                  <div className="bg-red-500/10 border border-red-500 rounded-md p-3 text-sm text-red-600">
                    {errorMessage}
                  </div>
                )}

                <Button
                  type="submit"
                  disabled={status === "loading" || !email || !code}
                  className="w-full"
                >
                  {status === "loading" ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      {t("verification.verifying")}
                    </>
                  ) : (
                    t("verification.verify")
                  )}
                </Button>

                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => navigate("/login?mode=login")}
                  className="w-full"
                >
                  {t("verification.backToLogin")}
                </Button>
              </form>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
