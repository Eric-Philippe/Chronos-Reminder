import { useState } from "react";
import { toast } from "sonner";
import {
  ChevronLeft,
  Zap,
  MessageCircle,
  Megaphone,
  Link2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Header } from "@/components/common/header";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { ReminderDetailsStep } from "@/components/reminder-wizard/ReminderDetailsStep";
import { DestinationsStep } from "@/components/reminder-wizard/DestinationsStep";
import { remindersService } from "@/services";
import { getRecurrenceTypeI18nKey } from "@/lib/recurrenceUtils";

export type ReminderStep = "details" | "destinations" | "review";

export interface ReminderFormData {
  date: Date | null;
  time: string; // HH:mm format
  message: string;
  recurrence: number;
  destinations: Array<{
    type: "discord_dm" | "discord_channel" | "webhook";
    metadata: Record<string, unknown>;
  }>;
}

export function CreateReminderPage() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [currentStep, setCurrentStep] = useState<ReminderStep>("details");

  // Initialize with current date and time
  const getInitialFormData = (): ReminderFormData => {
    const now = new Date();
    const hours = String(now.getHours()).padStart(2, "0");
    const minutes = String(now.getMinutes()).padStart(2, "0");

    return {
      date: now,
      time: `${hours}:${minutes}`,
      message: "",
      recurrence: 0, // RecurrenceOnce
      destinations: [],
    };
  };

  const [formData, setFormData] = useState<ReminderFormData>(
    getInitialFormData()
  );
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleNext = () => {
    if (currentStep === "details") {
      if (!formData.date || !formData.message.trim()) {
        return; // Validation will be handled in the step component
      }
      setCurrentStep("destinations");
    } else if (currentStep === "destinations") {
      setCurrentStep("review");
    }
  };

  const handleBack = () => {
    if (currentStep === "destinations") {
      setCurrentStep("details");
    } else if (currentStep === "review") {
      setCurrentStep("destinations");
    }
  };

  const handleCreate = async () => {
    setIsLoading(true);
    setError(null);
    try {
      if (!formData.date) {
        setError(t("reminderCreation.errors.selectDate"));
        setIsLoading(false);
        return;
      }

      // Convert the Date object to local date string (YYYY-MM-DD)
      // This ensures we get the date in the user's local timezone
      const dateStr = formData.date.toLocaleDateString("en-CA"); // "en-CA" gives YYYY-MM-DD format

      // Call API to create reminder
      // The backend will parse this date and time using the user's timezone
      const result = await remindersService.createReminder({
        date: dateStr,
        time: formData.time,
        message: formData.message,
        recurrence: formData.recurrence,
        destinations: formData.destinations,
      });

      if (result) {
        // Show success toast
        toast.success("Reminder created successfully!", {
          description: `"${
            formData.message
          }" will remind you on ${formData.date?.toLocaleDateString()} at ${
            formData.time
          }`,
          duration: 3000,
        });

        // After successful creation, navigate back to dashboard
        navigate("/dashboard");
      } else {
        setError(t("reminderCreation.errors.failed"));
        toast.error(t("reminderCreation.errors.failed"));
      }
    } catch (err) {
      console.error("Failed to create reminder:", err);
      const errorMessage =
        err instanceof Error
          ? err.message
          : t("reminderCreation.errors.failed");
      setError(errorMessage);
      toast.error("Failed to create reminder", {
        description: errorMessage,
      });
    } finally {
      setIsLoading(false);
    }
  };

  const getStepNumber = (): number => {
    switch (currentStep) {
      case "details":
        return 1;
      case "destinations":
        return 2;
      case "review":
        return 3;
      default:
        return 1;
    }
  };

  return (
    <div className="min-h-screen bg-background-main dark:bg-background-main">
      <Header />

      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Back Button */}
        <Button
          onClick={() => navigate("/dashboard")}
          variant="ghost"
          className="mb-8 text-foreground hover:bg-secondary/50 gap-2"
        >
          <ChevronLeft className="w-4 h-4" />
          {t("reminderCreation.backToDashboard")}
        </Button>

        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-foreground mb-2">
            {t("reminderCreation.title")}
          </h1>
          <p className="text-muted-foreground">
            {t("reminderCreation.step", {
              current: getStepNumber(),
              total: 3,
            })}
          </p>
        </div>

        {/* Progress Indicator */}
        <div className="mb-8 flex gap-2">
          {(["details", "destinations", "review"] as const).map((step, idx) => (
            <div
              key={step}
              className={`flex-1 h-2 rounded-full transition-colors ${
                idx < getStepNumber()
                  ? "bg-accent"
                  : idx === getStepNumber() - 1
                  ? "bg-accent/50"
                  : "bg-secondary/50"
              }`}
            />
          ))}
        </div>

        {/* Error Alert */}
        {error && (
          <Card className="border-red-500/50 bg-red-500/10 backdrop-blur mb-8">
            <CardContent className="pt-6">
              <p className="text-red-600 dark:text-red-400">{error}</p>
            </CardContent>
          </Card>
        )}

        {/* Step Content */}
        <Card className="border-border bg-card/95 backdrop-blur mb-8">
          <CardContent className="pt-6">
            {currentStep === "details" && (
              <ReminderDetailsStep
                formData={formData}
                onFormChange={setFormData}
              />
            )}

            {currentStep === "destinations" && (
              <DestinationsStep
                formData={formData}
                onFormChange={setFormData}
              />
            )}

            {currentStep === "review" && (
              <div className="space-y-6">
                <div>
                  <h2 className="text-xl font-bold text-foreground mb-4">
                    {t("reminderCreation.review.title")}
                  </h2>
                  <div className="space-y-4">
                    {/* Date & Time */}
                    <div className="p-4 rounded-lg border border-border bg-secondary/20">
                      <p className="text-sm text-muted-foreground mb-1">
                        {t("reminderCreation.review.dateTime")}
                      </p>
                      <p className="text-foreground font-semibold">
                        {formData.date?.toLocaleDateString()} at {formData.time}
                      </p>
                    </div>

                    {/* Message */}
                    <div className="p-4 rounded-lg border border-border bg-secondary/20">
                      <p className="text-sm text-muted-foreground mb-1">
                        {t("reminderCreation.review.message")}
                      </p>
                      <p className="text-foreground">{formData.message}</p>
                    </div>

                    {/* Recurrence */}
                    <div className="p-4 rounded-lg border border-border bg-secondary/20">
                      <p className="text-sm text-muted-foreground mb-1">
                        {t("reminderCreation.review.recurrence")}
                      </p>
                      <p className="text-foreground font-semibold">
                        {t(getRecurrenceTypeI18nKey(formData.recurrence))}
                      </p>
                    </div>

                    {/* Destinations */}
                    <div className="p-4 rounded-lg border border-border bg-secondary/20">
                      <p className="text-sm text-muted-foreground mb-3">
                        {t("reminderCreation.review.destinations")} (
                        {formData.destinations.length})
                      </p>
                      {formData.destinations.length > 0 ? (
                        <div className="space-y-2">
                          {formData.destinations.map((dest, idx) => (
                            <div
                              key={idx}
                              className="text-sm text-foreground px-3 py-2 rounded bg-accent/10 border border-accent/20 flex items-center gap-2"
                            >
                              {dest.type === "discord_dm" && (
                                <>
                                  <MessageCircle className="w-4 h-4 flex-shrink-0" />
                                  {t("reminderCreation.destinations.discordDM")}
                                </>
                              )}
                              {dest.type === "discord_channel" && (
                                <>
                                  <Megaphone className="w-4 h-4 flex-shrink-0" />
                                  {t(
                                    "reminderCreation.destinations.discordGuild"
                                  )}
                                </>
                              )}
                              {dest.type === "webhook" && (
                                <>
                                  <Link2 className="w-4 h-4 flex-shrink-0" />
                                  {t("reminderCreation.destinations.webhook")}
                                </>
                              )}
                            </div>
                          ))}
                        </div>
                      ) : (
                        <p className="text-sm text-muted-foreground italic">
                          {t("reminderCreation.review.noDestinations")}
                        </p>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Navigation Buttons */}
        <div className="flex gap-3">
          {currentStep !== "details" && (
            <Button
              onClick={handleBack}
              variant="outline"
              className="flex-1 border-border text-foreground hover:bg-secondary/50"
            >
              <ChevronLeft className="w-4 h-4 mr-2" />
              {t("reminderCreation.buttons.back")}
            </Button>
          )}

          {currentStep !== "review" && (
            <Button
              onClick={handleNext}
              className="flex-1 bg-accent hover:bg-accent/90 text-accent-foreground font-semibold"
              disabled={
                currentStep === "details" &&
                (!formData.date || !formData.message.trim())
              }
            >
              {t("reminderCreation.buttons.next")}
              <ChevronLeft className="w-4 h-4 ml-2 rotate-180" />
            </Button>
          )}

          {currentStep === "review" && (
            <Button
              onClick={handleCreate}
              className="flex-1 bg-accent hover:bg-accent/90 text-accent-foreground font-semibold"
              disabled={isLoading}
            >
              {isLoading
                ? t("reminderCreation.buttons.creating")
                : t("reminderCreation.buttons.create")}
              <Zap className="w-4 h-4 ml-2" />
            </Button>
          )}
        </div>
      </main>
    </div>
  );
}
