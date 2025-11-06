import { useTranslation } from "react-i18next";
import { WorldClocks } from "@/components/Clock";

export function WorldClocksSection() {
  const { t } = useTranslation();
  return (
    <WorldClocks
      title={t("vitrine.worldClocksTitle")}
      subtitle={t("vitrine.worldClocksSubtitle")}
    />
  );
}
