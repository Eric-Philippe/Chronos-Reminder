import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";

export function CTASection() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
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
            <a
              href="https://docs.chronosreminder.com"
              target="_blank"
              rel="noopener noreferrer"
            >
              {t("vitrine.viewDocs")}
            </a>
          </Button>
        </div>
      </div>
    </section>
  );
}
