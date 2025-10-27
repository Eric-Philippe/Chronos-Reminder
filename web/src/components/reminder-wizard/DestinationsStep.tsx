import { useTranslation } from "react-i18next";
import { DestinationPicker } from "./DestinationPicker";
import type { ReminderFormData } from "@/pages/CreateReminderPage";

interface DestinationsStepProps {
  formData: ReminderFormData;
  onFormChange: (data: ReminderFormData) => void;
}

export function DestinationsStep({
  formData,
  onFormChange,
}: DestinationsStepProps) {
  const { t } = useTranslation();

  const handleDestinationsChange = (
    destinations: typeof formData.destinations
  ) => {
    onFormChange({
      ...formData,
      destinations,
    });
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold text-foreground mb-4">
          {t("reminderCreation.destinations.title")}
        </h2>
        <p className="text-sm text-muted-foreground mb-6">
          {t("reminderCreation.destinations.subtitle")}
        </p>
      </div>

      <DestinationPicker
        destinations={formData.destinations}
        onDestinationsChange={handleDestinationsChange}
        showTitle={false}
        showAddOptions={true}
      />
    </div>
  );
}
