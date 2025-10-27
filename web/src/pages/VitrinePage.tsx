import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import {
  Clock,
  Bell,
  Zap,
  Globe,
  ArrowRight,
  CheckCircle2,
} from "lucide-react";
import { Header } from "@/components/common/header";
import { Button } from "@/components/ui/button";

export function VitrinePage() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-gradient-to-br from-background-main to-background-secondary">
      <Header />

      <main className="flex-1">
        {/* Hero Section */}
        <section className="relative overflow-hidden pt-12 pb-20 px-4 sm:px-6 lg:px-8">
          {/* Decorative Background Elements */}
          <div className="absolute inset-0 overflow-hidden pointer-events-none">
            <div className="absolute top-20 right-10 w-72 h-72 bg-accent/10 rounded-full blur-3xl dark:bg-accent/5"></div>
            <div className="absolute bottom-0 left-0 w-96 h-96 bg-accent/15 rounded-full blur-3xl dark:bg-accent/10"></div>
          </div>

          <div className="max-w-7xl mx-auto relative z-10">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center min-h-[70vh]">
              {/* Left Column - Hero Text */}
              <div className="space-y-8">
                <div>
                  <h1 className="text-5xl sm:text-6xl font-bold text-foreground leading-tight mb-6">
                    {t("vitrine.heroTitle")}
                  </h1>
                  <p className="text-lg sm:text-xl text-muted-foreground mb-8 leading-relaxed">
                    {t("vitrine.heroSubtitle")}
                  </p>
                </div>

                <div className="space-y-3">
                  <div className="flex items-center gap-3">
                    <CheckCircle2 className="w-5 h-5 text-accent flex-shrink-0" />
                    <span className="text-foreground">
                      {t("vitrine.feature1")}
                    </span>
                  </div>
                  <div className="flex items-center gap-3">
                    <CheckCircle2 className="w-5 h-5 text-accent flex-shrink-0" />
                    <span className="text-foreground">
                      {t("vitrine.feature2")}
                    </span>
                  </div>
                  <div className="flex items-center gap-3">
                    <CheckCircle2 className="w-5 h-5 text-accent flex-shrink-0" />
                    <span className="text-foreground">
                      {t("vitrine.feature3")}
                    </span>
                  </div>
                  <div className="flex items-center gap-3">
                    <CheckCircle2 className="w-5 h-5 text-accent flex-shrink-0" />
                    <span className="text-foreground">
                      {t("vitrine.feature4")}
                    </span>
                  </div>
                </div>

                <div className="flex flex-col sm:flex-row gap-4 pt-4">
                  <Button
                    onClick={() => navigate("/login")}
                    className="bg-accent hover:bg-accent/90 text-accent-foreground px-8 py-6 rounded-lg font-semibold flex items-center justify-center gap-2 text-lg"
                  >
                    {t("vitrine.getStarted")}
                    <ArrowRight className="w-5 h-5" />
                  </Button>
                  <Button
                    onClick={() => {
                      document
                        .getElementById("features")
                        ?.scrollIntoView({ behavior: "smooth" });
                    }}
                    variant="outline"
                    className="px-8 py-6 rounded-lg font-semibold text-lg"
                  >
                    {t("vitrine.learnMore")}
                  </Button>
                </div>
              </div>

              {/* Right Column - Hero Image/Illustration */}
              <div className="relative h-full min-h-[400px] hidden lg:flex items-center justify-center">
                <div className="relative w-full h-full">
                  {/* Decorative Box with Clock Icon */}
                  <div className="absolute inset-0 rounded-2xl bg-gradient-to-br from-accent/10 to-accent/5 border border-accent/20 flex items-center justify-center">
                    <Clock className="w-32 h-32 text-accent/40" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Features Section */}
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
              {/* Feature 1: Smart Reminders */}
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

              {/* Feature 2: Always On Time */}
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

              {/* Feature 3: Multi-Platform */}
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

              {/* Feature 4: Timezone Support */}
              <div className="group p-6 rounded-xl bg-background/50 border border-white/5 hover:border-accent/50 transition-all duration-300 hover:shadow-lg hover:shadow-accent/10">
                <div className="w-12 h-12 rounded-lg bg-accent/20 flex items-center justify-center mb-4 group-hover:bg-accent/30 transition-colors">
                  <Clock className="w-6 h-6 text-accent" />
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

        {/* Benefits Section */}
        <section className="py-20 px-4 sm:px-6 lg:px-8">
          <div className="max-w-7xl mx-auto">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
              {/* Left - Content */}
              <div className="space-y-8">
                <div>
                  <h2 className="text-4xl font-bold text-foreground mb-4">
                    {t("vitrine.benefitsTitle")}
                  </h2>
                  <p className="text-lg text-muted-foreground">
                    {t("vitrine.benefitsSubtitle")}
                  </p>
                </div>

                <div className="space-y-6">
                  <div className="flex gap-4">
                    <div className="flex-shrink-0 w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center">
                      <span className="text-accent font-bold">1</span>
                    </div>
                    <div>
                      <h3 className="text-lg font-semibold text-foreground mb-1">
                        {t("vitrine.benefit1Title")}
                      </h3>
                      <p className="text-muted-foreground">
                        {t("vitrine.benefit1Desc")}
                      </p>
                    </div>
                  </div>

                  <div className="flex gap-4">
                    <div className="flex-shrink-0 w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center">
                      <span className="text-accent font-bold">2</span>
                    </div>
                    <div>
                      <h3 className="text-lg font-semibold text-foreground mb-1">
                        {t("vitrine.benefit2Title")}
                      </h3>
                      <p className="text-muted-foreground">
                        {t("vitrine.benefit2Desc")}
                      </p>
                    </div>
                  </div>

                  <div className="flex gap-4">
                    <div className="flex-shrink-0 w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center">
                      <span className="text-accent font-bold">3</span>
                    </div>
                    <div>
                      <h3 className="text-lg font-semibold text-foreground mb-1">
                        {t("vitrine.benefit3Title")}
                      </h3>
                      <p className="text-muted-foreground">
                        {t("vitrine.benefit3Desc")}
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              {/* Right - Image/Illustration */}
              <div className="relative h-full min-h-[400px] hidden lg:flex items-center justify-center">
                <div className="relative w-full h-full">
                  <div className="absolute inset-0 rounded-2xl bg-gradient-to-br from-accent/5 to-accent/10 border border-accent/20 flex items-center justify-center">
                    <Bell className="w-32 h-32 text-accent/40" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="py-20 px-4 sm:px-6 lg:px-8 bg-secondary/50">
          <div className="max-w-4xl mx-auto text-center">
            <h2 className="text-4xl sm:text-5xl font-bold text-foreground mb-6">
              {t("vitrine.ctaTitle")}
            </h2>
            <p className="text-lg text-muted-foreground mb-8 max-w-2xl mx-auto">
              {t("vitrine.ctaSubtitle")}
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button
                onClick={() => navigate("/login")}
                className="bg-accent hover:bg-accent/90 text-accent-foreground px-8 py-6 rounded-lg font-semibold text-lg"
              >
                {t("vitrine.startFree")}
              </Button>
              <Button
                variant="outline"
                className="px-8 py-6 rounded-lg font-semibold text-lg"
              >
                {t("vitrine.viewDocs")}
              </Button>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
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
    </div>
  );
}
