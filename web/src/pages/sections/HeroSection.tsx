import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { ArrowRight, CheckCircle2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Hourglass3D } from "@/components/Hourglass3D";

export function HeroSection() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <section className="relative overflow-hidden pt-4 pb-40 px-4 sm:px-6 lg:px-8">
      {/* Decorative Background Elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-20 right-10 w-72 h-72 bg-accent/10 rounded-full blur-3xl dark:bg-accent/5"></div>
        <div className="absolute bottom-0 left-0 w-96 h-96 bg-accent/15 rounded-full blur-3xl dark:bg-accent/10"></div>
      </div>

      <div className="max-w-7xl mx-auto relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-start">
          {/* Mobile Hourglass - shown on small screens, ordered first */}
          <div className="lg:hidden flex flex-col items-center justify-center mb-2 order-first">
            <div className="relative w-full h-[250px]">
              <Hourglass3D />
            </div>
          </div>

          {/* Left Column - Hero Text */}
          <div className="space-y-2 lg:pt-16 lg:mt-16">
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
                <span className="text-foreground">{t("vitrine.feature1")}</span>
              </div>
              <div className="flex items-center gap-3">
                <CheckCircle2 className="w-5 h-5 text-accent flex-shrink-0" />
                <span className="text-foreground">{t("vitrine.feature2")}</span>
              </div>
              <div className="flex items-center gap-3">
                <CheckCircle2 className="w-5 h-5 text-accent flex-shrink-0" />
                <span className="text-foreground">{t("vitrine.feature3")}</span>
              </div>
              <div className="flex items-center gap-3">
                <CheckCircle2 className="w-5 h-5 text-accent flex-shrink-0" />
                <span className="text-foreground">{t("vitrine.feature4")}</span>
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
          <div className="relative h-full hidden lg:flex items-center justify-center">
            <div className="relative w-full h-[800px]">
              {/* 3D Hourglass */}
              <div className="w-full h-full rounded-2xl overflow-hidden">
                <Hourglass3D />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
