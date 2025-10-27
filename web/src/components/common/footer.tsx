import { useTranslation } from "react-i18next";

export function Footer() {
  const { t } = useTranslation();

  return (
    <footer className="border-t border-border mt-16 bg-background-main dark:bg-background-main">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex flex-col md:flex-row items-center justify-between text-muted-foreground text-sm gap-4">
          <p>{t("footer.copyright")}</p>
          <div className="flex gap-6 mt-4 md:mt-0">
            <a href="#" className="hover:text-accent transition-colors">
              {t("footer.documentation")}
            </a>
            <a href="#" className="hover:text-accent transition-colors">
              {t("footer.support")}
            </a>
            <a href="#" className="hover:text-accent transition-colors">
              {t("footer.discord")}
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
