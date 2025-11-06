import { useTranslation } from "react-i18next";

export function FooterSection() {
  const { t } = useTranslation();
  return (
    <footer className="border-t border-white/5 py-12 px-4 sm:px-6 lg:px-8 bg-background/50">
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
          <div>
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
          <div>
            <h4 className="font-semibold text-foreground mb-4">
              {t("vitrine.footerProduct")}
            </h4>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li>
                <a href="#" className="hover:text-accent transition-colors">
                  {t("vitrine.footerFeatures")}
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-accent transition-colors">
                  {t("vitrine.footerPricing")}
                </a>
              </li>
            </ul>
          </div>
          <div>
            <h4 className="font-semibold text-foreground mb-4">
              {t("vitrine.footerResources")}
            </h4>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li>
                <a href="#" className="hover:text-accent transition-colors">
                  {t("footer.documentation")}
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-accent transition-colors">
                  {t("footer.support")}
                </a>
              </li>
            </ul>
          </div>
          <div>
            <h4 className="font-semibold text-foreground mb-4">
              {t("vitrine.footerCommunity")}
            </h4>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li>
                <a href="#" className="hover:text-accent transition-colors">
                  {t("footer.discord")}
                </a>
              </li>
              <li>
                <a href="#" className="hover:text-accent transition-colors">
                  {t("vitrine.footerTwitter")}
                </a>
              </li>
            </ul>
          </div>
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
