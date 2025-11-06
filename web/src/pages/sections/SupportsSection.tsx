import { useTranslation } from "react-i18next";
import { Webhook, Key, Mail, Sparkles } from "lucide-react";
import { IntegrationCard } from "@/components/IntegrationCard";

export function SupportsSection() {
  const { t } = useTranslation();
  return (
    <section className="py-20 px-4 sm:px-6 lg:px-8 bg-secondary/30">
      <div className="max-w-7xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl sm:text-5xl font-bold text-foreground mb-4">
            {t("vitrine.supportsTitle")}
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            {t("vitrine.supportsSubtitle")}
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
          <IntegrationCard
            name={t("vitrine.supportsDiscord")}
            icon={
              <svg
                className="w-10 h-10 text-accent"
                viewBox="0 0 24 24"
                fill="currentColor"
              >
                <path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515a.074.074 0 0 0-.079.037c-.211.375-.444.864-.607 1.25a18.27 18.27 0 0 0-5.487 0c-.163-.386-.395-.875-.607-1.25a.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057a19.9 19.9 0 0 0 5.993 3.03a.078.078 0 0 0 .084-.028c.462-.63.873-1.295 1.226-1.994a.076.076 0 0 0-.042-.106a13.107 13.107 0 0 1-1.872-.892a.077.077 0 0 1 .008-.128c.125-.093.25-.19.371-.287a.074.074 0 0 1 .076-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .077.009c.12.098.246.195.371.288a.077.077 0 0 1 .009.127a13.073 13.073 0 0 1-1.872.892a.077.077 0 0 0-.041.107c.36.699.77 1.364 1.225 1.994a.076.076 0 0 0 .084.028a19.86 19.86 0 0 0 6.002-3.03a.077.077 0 0 0 .032-.057c.5-4.506-.838-8.42-3.549-11.59a.06.06 0 0 0-.031-.028zM8.02 15.33c-.999 0-1.823-.915-1.823-2.03c0-1.114.823-2.03 1.823-2.03c1.000 0 1.823.916 1.823 2.03c0 1.114-.823 2.03-1.823 2.03zm7.983 0c-.999 0-1.823-.915-1.823-2.03c0-1.114.824-2.03 1.823-2.03c1.000 0 1.823.916 1.823 2.03c0 1.114-.823 2.03-1.823 2.03z" />
              </svg>
            }
            description={t("vitrine.discordDesc")}
            features={[
              t("vitrine.discordFeature1"),
              t("vitrine.discordFeature2"),
              t("vitrine.discordFeature3"),
            ]}
          />

          <IntegrationCard
            name={t("vitrine.supportsSlack")}
            icon={
              <svg
                className="w-10 h-10 text-accent"
                viewBox="0 0 24 24"
                fill="currentColor"
              >
                <path d="M5.042 15.165a2.528 2.528 0 0 1-2.52 2.523A2.528 2.528 0 0 1 0 15.165a2.527 2.527 0 0 1 2.522-2.52h2.52v2.52zM6.313 15.165a2.527 2.527 0 0 1 2.521-2.52 2.528 2.528 0 0 1 2.524 2.52v6.31A2.529 2.529 0 0 1 8.834 24a2.529 2.529 0 0 1-2.521-2.525v-6.31zM8.834 5.042a2.528 2.528 0 0 1-2.521-2.52A2.528 2.528 0 0 1 8.834 0a2.528 2.528 0 0 1 2.521 2.522v2.52H8.834zM8.834 6.313a2.528 2.528 0 0 1 2.521 2.521 2.528 2.528 0 0 1-2.521 2.524H2.524A2.528 2.528 0 0 1 0 8.834a2.528 2.528 0 0 1 2.524-2.521h6.31zM18.958 8.834a2.528 2.528 0 0 1 2.521-2.521A2.528 2.528 0 0 1 24 8.834a2.528 2.528 0 0 1-2.521 2.524h-2.521v-2.524zM17.687 8.834a2.528 2.528 0 0 1-2.521 2.524 2.528 2.528 0 0 1-2.521-2.524V2.524A2.528 2.528 0 0 1 15.166 0a2.528 2.528 0 0 1 2.521 2.524v6.31zM15.166 18.958a2.528 2.528 0 0 1 2.521 2.521A2.528 2.528 0 0 1 15.166 24a2.527 2.527 0 0 1-2.521-2.521v-2.521h2.521zM15.166 17.687a2.528 2.528 0 0 1-2.521-2.521 2.528 2.528 0 0 1 2.521-2.524h6.31a2.527 2.527 0 0 1 2.521 2.524 2.527 2.527 0 0 1-2.521 2.521h-6.31z" />
              </svg>
            }
            description={t("vitrine.slackDesc")}
            features={[t("vitrine.slackFeature1"), t("vitrine.slackFeature2")]}
          />

          <IntegrationCard
            name={t("vitrine.supportsWebhook")}
            icon={<Webhook className="w-10 h-10 text-accent" />}
            description={t("vitrine.webhookDesc")}
            features={[
              t("vitrine.webhookFeature1"),
              t("vitrine.webhookFeature2"),
            ]}
          />

          <IntegrationCard
            name={t("vitrine.supportsApiKey")}
            icon={<Key className="w-10 h-10 text-accent" />}
            description={t("vitrine.apiKeyDesc")}
            features={[
              t("vitrine.apiKeyFeature1"),
              t("vitrine.apiKeyFeature2"),
            ]}
            isApiKeyCard={true}
          />

          <div className="relative h-full group">
            <div className="relative p-6 rounded-xl bg-background/50 border border-white/5 hover:border-accent/30 h-full flex flex-col items-center justify-center gap-4 cursor-not-allowed opacity-75 hover:opacity-90 transition-all duration-300">
              <div className="absolute top-3 right-3">
                <div className="flex items-center gap-1 px-3 py-1 rounded-full bg-accent/20 border border-accent/50">
                  <Sparkles className="w-3 h-3 text-accent" />
                  <span className="text-xs font-semibold text-accent">
                    Coming Soon
                  </span>
                </div>
              </div>
              <div className="w-16 h-16 rounded-lg bg-accent/20 flex items-center justify-center">
                <Mail className="w-10 h-10 text-accent" />
              </div>
              <h3 className="text-lg font-semibold text-foreground text-center">
                {t("vitrine.supportsEmail")}
              </h3>
              <p className="text-sm text-muted-foreground text-center leading-relaxed">
                {t("vitrine.emailComingSoonDesc")}
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
