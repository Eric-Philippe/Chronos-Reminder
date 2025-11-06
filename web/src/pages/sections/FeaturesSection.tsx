import { useTranslation } from "react-i18next";
import { Bell, Zap, Globe, Clock as ClockIcon } from "lucide-react";

export function FeaturesSection() {
  const { t } = useTranslation();

  return (
    <section
      id="features"
      className="py-20 px-4 sm:px-6 lg:px-8 bg-secondary/30"
    >
      <div className="max-w-7xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl sm:text-5xl font-bold text-foreground mb-4">
            {t("vitrine.featuresTitle")}
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            {t("vitrine.featuresSubtitle")}
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          <div className="group p-6 rounded-xl bg-background/50 border border-white/5 hover:border-accent/50 transition-all duration-300 hover:shadow-lg hover:shadow-accent/10">
            <div className="w-12 h-12 rounded-lg bg-accent/20 flex items-center justify-center mb-4 group-hover:bg-accent/30 transition-colors">
              <Bell className="w-6 h-6 text-accent" />
            </div>
            <h3 className="text-lg font-semibold text-foreground mb-2">
              {t("vitrine.smartRemindersTitle")}
            </h3>
            <p className="text-sm text-muted-foreground">
              {t("vitrine.smartRemindersDesc")}
            </p>
          </div>

          <div className="group p-6 rounded-xl bg-background/50 border border-white/5 hover:border-accent/50 transition-all duration-300 hover:shadow-lg hover:shadow-accent/10">
            <div className="w-12 h-12 rounded-lg bg-accent/20 flex items-center justify-center mb-4 group-hover:bg-accent/30 transition-colors">
              <Zap className="w-6 h-6 text-accent" />
            </div>
            <h3 className="text-lg font-semibold text-foreground mb-2">
              {t("vitrine.alwaysOnTimeTitle")}
            </h3>
            <p className="text-sm text-muted-foreground">
              {t("vitrine.alwaysOnTimeDesc")}
            </p>
          </div>

          <div className="group p-6 rounded-xl bg-background/50 border border-white/5 hover:border-accent/50 transition-all duration-300 hover:shadow-lg hover:shadow-accent/10">
            <div className="w-12 h-12 rounded-lg bg-accent/20 flex items-center justify-center mb-4 group-hover:bg-accent/30 transition-colors">
              <Globe className="w-6 h-6 text-accent" />
            </div>
            <h3 className="text-lg font-semibold text-foreground mb-2">
              {t("vitrine.multiPlatformTitle")}
            </h3>
            <p className="text-sm text-muted-foreground">
              {t("vitrine.multiPlatformDesc")}
            </p>
          </div>

          <div className="group p-6 rounded-xl bg-background/50 border border-white/5 hover:border-accent/50 transition-all duration-300 hover:shadow-lg hover:shadow-accent/10">
            <div className="w-12 h-12 rounded-lg bg-accent/20 flex items-center justify-center mb-4 group-hover:bg-accent/30 transition-colors">
              <ClockIcon className="w-6 h-6 text-accent" />
            </div>
            <h3 className="text-lg font-semibold text-foreground mb-2">
              {t("vitrine.timezoneTitle")}
            </h3>
            <p className="text-sm text-muted-foreground">
              {t("vitrine.timezoneDesc")}
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
