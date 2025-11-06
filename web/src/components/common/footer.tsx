import { useTranslation } from "react-i18next";
import { MENU_GROUPS, ROUTES } from "@/config/routes";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";

export function Footer() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();

  const handleNavigation = (path: string, external?: boolean) => {
    if (external) {
      window.open(path, "_blank");
    } else {
      navigate(path);
    }
  };

  return (
    <footer className="border-t border-white/5 py-12 px-4 sm:px-6 lg:px-8 bg-background/50">
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-6 gap-4 mb-8">
          {/* Logo Section */}
          <div className="md:col-span-1.5 mr-8">
            <div className="flex items-center gap-2 mb-4">
              <img
                src="/logo_chronos.png"
                alt="Chronos"
                className="w-8 h-8 rounded-full"
              />
              <span className="text-lg font-bold text-accent">Chronos</span>
            </div>
            <p className="text-sm text-muted-foreground">
              {t("vitrine.footerDesc")}
            </p>
          </div>

          {/* Dynamic Menu Groups - Filter by requiresAuth */}
          {MENU_GROUPS.filter(
            (group) => !group.requiresAuth || isAuthenticated
          ).map((group) => (
            <div key={group.name}>
              <h4 className="font-semibold text-foreground mb-4">
                {t(`header.${group.label}`)}
              </h4>
              <ul className="space-y-1 text-sm text-muted-foreground">
                {group.items.map((item) => (
                  <li key={item.path}>
                    <button
                      onClick={() => handleNavigation(item.path, item.external)}
                      className="hover:text-accent transition-colors cursor-pointer text-left"
                    >
                      {t(`header.${item.name}`)}
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          ))}

          {/* Authentication Section - Only show if not authenticated */}
          {!isAuthenticated && (
            <div>
              <h4 className="font-semibold text-foreground mb-4">
                {t("header.account") || "Account"}
              </h4>
              <ul className="space-y-1 text-sm text-muted-foreground">
                <li>
                  <button
                    onClick={() => navigate(`${ROUTES.LOGIN.path}?mode=login`)}
                    className="hover:text-accent transition-colors cursor-pointer text-left"
                  >
                    {t("header.signIn") || "Log In"}
                  </button>
                </li>
                <li>
                  <button
                    onClick={() => navigate(`${ROUTES.LOGIN.path}?mode=signup`)}
                    className="hover:text-accent transition-colors cursor-pointer text-left"
                  >
                    {t("header.signUp") || "Create account"}
                  </button>
                </li>
              </ul>
            </div>
          )}
        </div>
        <div className="border-t border-white/5 pt-8">
          <p className="text-center text-sm text-muted-foreground">
            {t("footer.copyright")}
          </p>
        </div>
      </div>
    </footer>
  );
}
