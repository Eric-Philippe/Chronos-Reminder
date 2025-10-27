import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Menu, X, LogOut } from "lucide-react";
import { useTranslation } from "react-i18next";
import { ModeToggle } from "@/components/common/mode-toggle";
import { LanguageSwitcher } from "@/components/common/language-switcher";
import { useAuth } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
import { ROUTES, ROUTES_ARRAY } from "@/config/routes";

export function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { logout, isLoading, isAuthenticated } = useAuth();

  const handleLogout = async () => {
    try {
      await logout();
      setTimeout(() => {
        navigate(ROUTES.LOGIN.path, { replace: true });
      }, 100);
    } catch (err) {
      console.error("Logout failed:", err);
    }
  };

  // Only redirect if user was previously authenticated but is now logged out
  useEffect(() => {
    // Only perform redirect if we're no longer checking auth
    if (!isLoading && !isAuthenticated) {
      // Check if user was previously authenticated by checking sessionStorage
      const wasAuthenticated = sessionStorage.getItem("wasAuthenticated");
      if (wasAuthenticated) {
        navigate(ROUTES.LOGIN.path, { replace: true });
        sessionStorage.removeItem("wasAuthenticated");
      }
    }
  }, [isLoading, isAuthenticated, navigate]);

  // Mark that user was authenticated when they log in
  useEffect(() => {
    if (isAuthenticated) {
      sessionStorage.setItem("wasAuthenticated", "true");
    }
  }, [isAuthenticated]);

  // Authenticated header (Dashboard view)
  return (
    <header className="sticky top-0 z-50 flex justify-center pt-4 px-4">
      <div className="max-w-7xl w-full py-2 px-6 rounded-2xl backdrop-blur-2xl border shadow-2xl bg-background-secondary/60 dark:bg-[rgba(39,39,37,0.4)] border-border dark:border-white/10">
        <div className="flex items-center justify-between">
          {/* Logo and Brand */}
          <div
            className="flex items-center gap-2 cursor-pointer"
            onClick={() => navigate(ROUTES.VITRINE.path)}
          >
            <div className="w-10 h-10 rounded-full overflow-hidden flex items-center justify-center bg-amber-400/10">
              <img
                src="/logo_chronos.png"
                alt="Chronos Logo"
                className="w-11 h-11 object-cover"
              />
            </div>
            <h1 className="text-xl font-bold text-amber-600 dark:text-yellow-400">
              Chronos
            </h1>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center gap-8">
            {ROUTES_ARRAY.map((route) => {
              if (
                !route.showInNav ||
                (route.requiresAuth && !isAuthenticated)
              ) {
                return null;
              }
              return (
                <a
                  key={route.path}
                  href={route.path}
                  className="text-foreground hover:text-amber-600 dark:hover:text-amber-400 transition-colors"
                >
                  {t(`header.${route.name}`)}
                </a>
              );
            })}
          </nav>

          {/* Right Side - Language Switcher, Theme Toggle, Logout Button & Mobile Menu Button */}
          <div className="flex items-center gap-2">
            <LanguageSwitcher />
            <ModeToggle />
            {isAuthenticated ? (
              <Button
                onClick={handleLogout}
                disabled={isLoading}
                variant="ghost"
                className="text-foreground dark:text-white hover:text-red-600 dark:hover:text-red-400 hover:bg-red-400/10 transition-colors hidden sm:flex"
              >
                Logout
              </Button>
            ) : (
              <>
                <Button
                  onClick={() => navigate(ROUTES.LOGIN.path)}
                  variant="ghost"
                  className="text-foreground dark:text-white hover:text-amber-600 dark:hover:text-amber-400 hover:bg-amber-400/10 transition-colors hidden sm:flex border"
                >
                  {t("header.startNow")}
                </Button>
              </>
            )}

            {/* Hamburger Menu Button */}
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="md:hidden text-foreground dark:text-white hover:text-amber-600 dark:hover:text-amber-400 transition-colors"
            >
              {mobileMenuOpen ? (
                <X className="w-6 h-6" />
              ) : (
                <Menu className="w-6 h-6" />
              )}
            </button>
          </div>
        </div>

        {/* Mobile Navigation Menu */}
        {mobileMenuOpen && (
          <nav className="md:hidden mt-4 pt-4 border-t border-border dark:border-white/10 flex flex-col gap-3">
            {ROUTES_ARRAY.map((route) => {
              if (
                !route.showInNav ||
                (route.requiresAuth && !isAuthenticated)
              ) {
                return null;
              }
              return (
                <a
                  key={route.path}
                  href={route.path}
                  className="text-foreground dark:text-white hover:text-amber-600 dark:hover:text-amber-400 transition-colors py-2"
                >
                  {t(`header.${route.name}`)}
                </a>
              );
            })}
            <hr className="border-white/10 my-2" />
            <Button
              onClick={handleLogout}
              disabled={isLoading}
              variant="ghost"
              className="w-full justify-start text-red-400 hover:text-red-300 hover:bg-red-400/10 transition-colors"
            >
              <LogOut className="w-4 h-4 mr-2" />
              {t("header.logout") || "Logout"}
            </Button>
          </nav>
        )}
      </div>
    </header>
  );
}
