import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { LanguageSwitcher } from "@/components/language-switcher";
import { ModeToggle } from "@/components/mode-toggle";
import { LoginSection } from "@/components/LoginSection";
import { SignUpSection } from "@/components/SignUpSection";
import { useAuth } from "@/hooks/useAuth";

export function LoginPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { login, register, isLoading, error } = useAuth();

  const [isSignUp, setIsSignUp] = useState(false);
  const [localError, setLocalError] = useState<string | null>(null);

  // Login state
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [rememberMe, setRememberMe] = useState(false);

  // Sign up state
  const [signUpEmail, setSignUpEmail] = useState("");
  const [username, setUsername] = useState("");
  const [signUpPassword, setSignUpPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [timezone, setTimezone] = useState("UTC");

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setLocalError(null);

    try {
      await login(email, password, rememberMe);
      // Navigate to welcome page on successful login
      navigate("/welcome", { replace: true });
    } catch (err) {
      // Error is already set in the hook's error state
      // The error will be displayed from the hook's error state via the error banner
      console.error("Login failed:", err);
    }
  };

  const handleSignUp = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setLocalError(null);

    if (signUpPassword !== confirmPassword) {
      setLocalError(t("login.passwordsDontMatch"));
      return;
    }

    try {
      await register(signUpEmail, username, signUpPassword, timezone);
      setLocalError(null);
      // Show success message and switch to login
      alert(t("login.registrationSuccess"));
      setIsSignUp(false);
      setSignUpEmail("");
      setUsername("");
      setSignUpPassword("");
      setConfirmPassword("");
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : "Registration failed";
      setLocalError(errorMsg);
      console.error("Registration failed:", err);
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

      <div className="w-full max-w-6xl relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 lg:gap-12 items-center min-h-[90vh]">
          {/* Left Column - Welcome Section (hidden on mobile) */}
          <div className="hidden lg:flex flex-col justify-center items-start space-y-8">
            <div>
              <div className="flex items-center gap-3 mb-6">
                <img
                  src="/logo_chronos.png"
                  alt="Chronos Logo"
                  className="w-12 h-12 rounded-full"
                />
                <span className="text-3xl font-bold text-accent">Chronos</span>
              </div>
              <h1 className="text-5xl font-bold text-foreground mb-4 leading-tight">
                {t("login.welcomeTitle")}
              </h1>
              <p className="text-lg text-muted-foreground mb-8">
                {t("login.welcomeSubtitle")}
              </p>
            </div>

            {/* Features List */}
            <div className="space-y-4">
              <div className="flex gap-3">
                <div className="flex-shrink-0">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-accent/20">
                    <span className="text-accent font-semibold">✓</span>
                  </div>
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">
                    {t("login.smartRemindersTitle")}
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    {t("login.smartRemindersDesc")}
                  </p>
                </div>
              </div>

              <div className="flex gap-3">
                <div className="flex-shrink-0">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-accent/20">
                    <span className="text-accent font-semibold">✓</span>
                  </div>
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">
                    {t("login.alwaysOnTimeTitle")}
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    {t("login.alwaysOnTimeDesc")}
                  </p>
                </div>
              </div>

              <div className="flex gap-3">
                <div className="flex-shrink-0">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-accent/20">
                    <span className="text-accent font-semibold">✓</span>
                  </div>
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">
                    {t("login.seamlessIntegrationTitle")}
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    {t("login.seamlessIntegrationDesc")}
                  </p>
                </div>
              </div>

              <div className="flex gap-3">
                <div className="flex-shrink-0">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-accent/20">
                    <span className="text-accent font-semibold">✓</span>
                  </div>
                </div>
                <div>
                  <h3 className="font-semibold text-foreground">
                    {t("login.allTimezoneTitle")}
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    {t("login.allTimezoneDesc")}
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Right Column - Login Form */}
          <div className="w-full">
            {/* Mobile Header (shown only on mobile) */}
            <div className="lg:hidden text-center mb-8">
              <div className="flex items-center justify-center gap-2 mb-4">
                <img
                  src="/logo_chronos.png"
                  alt="Chronos Logo"
                  className="w-10 h-10 rounded-full"
                />
                <span className="text-2xl font-bold text-accent">Chronos</span>
              </div>
              <h1 className="text-3xl font-bold text-foreground mb-2">
                {t("login.welcomeTitle")}
              </h1>
              <p className="text-muted-foreground">
                {t("login.mobileSubtitle")}
              </p>
            </div>

            {/* Tab Toggle */}
            <div className="flex gap-2 mb-6 bg-secondary/30 rounded-lg p-1">
              <button
                type="button"
                onClick={() => setIsSignUp(false)}
                className={`flex-1 py-2 px-4 rounded-md font-medium transition-colors ${
                  !isSignUp
                    ? "bg-accent text-accent-foreground"
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                {t("login.signIn")}
              </button>
              <button
                type="button"
                onClick={() => setIsSignUp(true)}
                className={`flex-1 py-2 px-4 rounded-md font-medium transition-colors ${
                  isSignUp
                    ? "bg-accent text-accent-foreground"
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                {t("login.signUp")}
              </button>
            </div>

            {/* Login Section */}
            {!isSignUp ? (
              <LoginSection
                email={email}
                setEmail={setEmail}
                password={password}
                setPassword={setPassword}
                rememberMe={rememberMe}
                setRememberMe={setRememberMe}
                onSubmit={handleLogin}
                isLoading={isLoading}
                localError={localError}
                authError={error}
              />
            ) : (
              <SignUpSection
                signUpEmail={signUpEmail}
                setSignUpEmail={setSignUpEmail}
                username={username}
                setUsername={setUsername}
                signUpPassword={signUpPassword}
                setSignUpPassword={setSignUpPassword}
                confirmPassword={confirmPassword}
                setConfirmPassword={setConfirmPassword}
                timezone={timezone}
                setTimezone={setTimezone}
                onSubmit={handleSignUp}
                isLoading={isLoading}
                localError={localError}
                authError={error}
              />
            )}

            {/* Footer */}
            <p className="text-center text-muted-foreground text-sm mt-6">
              {isSignUp
                ? t("login.alreadyHaveAccount")
                : t("login.dontHaveAccount")}{" "}
              <button
                type="button"
                onClick={() => setIsSignUp(!isSignUp)}
                className="text-accent hover:text-accent/80 transition-colors font-medium"
              >
                {isSignUp ? t("login.signIn") : t("login.signUp")}
              </button>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
