import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Menu, X, LogOut } from "lucide-react";
import { useTranslation } from "react-i18next";
import { ModeToggle } from "@/components/common/mode-toggle";
import { LanguageSwitcher } from "@/components/common/language-switcher";
import { useAuth } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";

export function Header() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { logout, isLoading, isAuthenticated } = useAuth();

  const handleLogout = async () => {
    try {
      await logout();
      // Force a small delay to ensure state updates propagate
      // Then navigate to login
      setTimeout(() => {
        navigate("/login", { replace: true });
      }, 100);
    } catch (err) {
      console.error("Logout failed:", err);
    }
  };

  // Watch for auth state changes and redirect if logged out
  useEffect(() => {
    if (!isAuthenticated && !isLoading) {
      navigate("/login", { replace: true });
    }
  }, [isAuthenticated, isLoading, navigate]);

  return (
    <header className="sticky top-0 z-50 flex justify-center pt-4 px-4">
      <div
        className="max-w-7xl w-full py-2 px-6 rounded-2xl backdrop-blur-2xl border border-white/10 shadow-2xl"
        style={{ backgroundColor: "rgba(39, 39, 37, 0.4)" }}
      >
        <div className="flex items-center justify-between">
          {/* Logo and Brand */}
          <div className="flex items-center gap-2">
            <div className="w-12 h-12 rounded-full overflow-hidden flex items-center justify-center bg-amber-400/10">
              <img
                src="/logo_chronos.png"
                alt="Chronos Logo"
                className="w-12 h-12 object-cover"
              />
            </div>
            <h1 className="text-xl font-bold text-yellow-400">Chronos</h1>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center gap-8">
            <a
              href="#"
              className="text-white hover:text-amber-400 transition-colors"
            >
              {t("header.myReminders")}
            </a>
            <a
              href="#"
              className="text-white hover:text-amber-400 transition-colors"
            >
              {t("header.integration")}
            </a>
            <a
              href="#"
              className="text-white hover:text-amber-400 transition-colors"
            >
              {t("header.settings")}
            </a>
            <a
              href="#"
              className="text-white hover:text-amber-400 transition-colors"
            >
              {t("header.myAccount")}
            </a>
          </nav>

          {/* Right Side - Language Switcher, Theme Toggle, Logout Button & Mobile Menu Button */}
          <div className="flex items-center gap-2">
            <LanguageSwitcher />
            <ModeToggle />
            <Button
              onClick={handleLogout}
              disabled={isLoading}
              variant="ghost"
              className="text-white hover:text-red-400 hover:bg-red-400/10 transition-colors hidden sm:flex"
              title={t("header.logout") || "Logout"}
            >
              <LogOut className="w-5 h-5" />
            </Button>
            {/* Hamburger Menu Button */}
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="md:hidden text-white hover:text-amber-400 transition-colors"
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
          <nav className="md:hidden mt-4 pt-4 border-t border-white/10 flex flex-col gap-3">
            <a
              href="#"
              className="text-white hover:text-amber-400 transition-colors py-2"
            >
              {t("header.myReminders")}
            </a>
            <a
              href="#"
              className="text-white hover:text-amber-400 transition-colors py-2"
            >
              {t("header.integration")}
            </a>
            <a
              href="#"
              className="text-white hover:text-amber-400 transition-colors py-2"
            >
              {t("header.settings")}
            </a>
            <a
              href="#"
              className="text-white hover:text-amber-400 transition-colors py-2"
            >
              {t("header.myAccount")}
            </a>
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
