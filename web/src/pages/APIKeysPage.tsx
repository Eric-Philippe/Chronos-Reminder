import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { Header } from "@/components/common/header";
import { Footer } from "@/components/common/footer";
import { APIKeysSection } from "@/components/APIKeysSection";
import { APITestGuide } from "@/components/APITestGuide";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { ROUTES } from "@/config/routes";

export function APIKeysPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-background-main dark:bg-background-main">
      <Header />

      {/* Main Content */}
      <main className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-12 pt-24">
        {/* Page Header */}
        <div className="mb-8 flex items-center gap-4">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => navigate(ROUTES.ACCOUNT.path)}
            className="gap-2 h-9 px-3"
          >
            <ArrowLeft className="w-4 h-4" />
            {t("common.back") || "Back"}
          </Button>
        </div>

        <div className="mb-12">
          <h1 className="text-4xl font-bold text-foreground mb-2">
            {t("apiKeys.title") || "API Keys"}
          </h1>
          <p className="text-muted-foreground">
            {t("apiKeys.description") ||
              "Manage API keys for programmatic access to your reminders"}
          </p>
        </div>

        <div className="grid gap-8">
          {/* API Keys Management */}
          <div>
            <APIKeysSection />
          </div>

          {/* Quick Guide - Full Width */}
          <div>
            <APITestGuide />
          </div>
        </div>
      </main>

      {/* Footer */}
      <Footer />
    </div>
  );
}
