import { Calendar, Clock, MessageSquare, Zap } from "lucide-react";
import { Button } from "@/components/ui/button";
import type { ReminderFormData } from "@/pages/CreateReminderPage";
import { useTranslation } from "react-i18next";
import {
  RecurrenceDailyStr,
  RecurrenceHourlyStr,
  RecurrenceMonthlyStr,
  RecurrenceOnceStr,
  RecurrenceWorkdaysStr,
  RecurrenceWeekendStr,
  RecurrenceWeeklyStr,
  RecurrenceYearlyStr,
  getRecurrenceTypeI18nKeyFromString,
} from "@/lib/recurrenceUtils";
import { isDateTimeInPast } from "@/lib/utils";

interface ReminderDetailsStepProps {
  formData: ReminderFormData;
  onFormChange: (data: ReminderFormData) => void;
}

export function ReminderDetailsStep({
  formData,
  onFormChange,
}: ReminderDetailsStepProps) {
  const { t } = useTranslation();

  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const dateStr = e.target.value;
    const date = dateStr ? new Date(dateStr) : null;
    onFormChange({ ...formData, date });
  };

  const handleTimeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onFormChange({ ...formData, time: e.target.value });
  };

  const handleMessageChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    onFormChange({ ...formData, message: e.target.value });
  };

  const handleRecurrenceChange = (recurrence: string) => {
    onFormChange({ ...formData, recurrence });
  };

  const recurrenceOptions = [
    {
      value: RecurrenceOnceStr,
      label: t(getRecurrenceTypeI18nKeyFromString(RecurrenceOnceStr)),
    },
    {
      value: RecurrenceDailyStr,
      label: t(getRecurrenceTypeI18nKeyFromString(RecurrenceDailyStr)),
    },
    {
      value: RecurrenceWeeklyStr,
      label: t(getRecurrenceTypeI18nKeyFromString(RecurrenceWeeklyStr)),
    },
    {
      value: RecurrenceMonthlyStr,
      label: t(getRecurrenceTypeI18nKeyFromString(RecurrenceMonthlyStr)),
    },
    {
      value: RecurrenceYearlyStr,
      label: t(getRecurrenceTypeI18nKeyFromString(RecurrenceYearlyStr)),
    },
    {
      value: RecurrenceHourlyStr,
      label: t(getRecurrenceTypeI18nKeyFromString(RecurrenceHourlyStr)),
    },
    {
      value: RecurrenceWorkdaysStr,
      label: t(getRecurrenceTypeI18nKeyFromString(RecurrenceWorkdaysStr)),
    },
    {
      value: RecurrenceWeekendStr,
      label: t(getRecurrenceTypeI18nKeyFromString(RecurrenceWeekendStr)),
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold text-foreground mb-4">
          {t("reminderCreation.details.title")}
        </h2>
        <p className="text-sm text-muted-foreground mb-6">
          {t("reminderCreation.details.subtitle")}
        </p>
      </div>

      {/* Date & Time Section */}
      <div className="space-y-4">
        <div className="p-4 rounded-lg border border-border bg-secondary/20">
          <div className="flex items-center gap-2 mb-3">
            <Calendar className="w-5 h-5 text-accent" />
            <label className="text-sm font-semibold text-foreground">
              {t("reminderCreation.details.date")}
            </label>
          </div>
          <input
            type="date"
            value={
              formData.date ? formData.date.toISOString().split("T")[0] : ""
            }
            onChange={handleDateChange}
            className="w-full px-3 py-2 rounded border border-border bg-background text-foreground"
          />
          {!formData.date && (
            <p className="text-xs text-red-500 mt-2">
              {t("reminderCreation.details.dateRequired")}
            </p>
          )}
        </div>

        <div className="p-4 rounded-lg border border-border bg-secondary/20">
          <div className="flex items-center gap-2 mb-3">
            <Clock className="w-5 h-5 text-accent" />
            <label className="text-sm font-semibold text-foreground">
              {t("reminderCreation.details.time")}
            </label>
          </div>
          <input
            type="time"
            value={formData.time}
            onChange={handleTimeChange}
            className="w-full px-3 py-2 rounded border border-border bg-background text-foreground"
          />
          {formData.date && isDateTimeInPast(formData.date, formData.time) && (
            <p className="text-xs text-red-500 mt-2">
              {t("reminderCreation.errors.dateTimeInPast")}
            </p>
          )}
        </div>
      </div>

      {/* Message Section */}
      <div className="p-4 rounded-lg border border-border bg-secondary/20">
        <div className="flex items-center gap-2 mb-3">
          <MessageSquare className="w-5 h-5 text-accent" />
          <label className="text-sm font-semibold text-foreground">
            {t("reminderCreation.details.message")}
          </label>
        </div>
        <textarea
          value={formData.message}
          onChange={handleMessageChange}
          placeholder={t("reminderCreation.details.messagePlaceholder")}
          className="w-full px-3 py-2 rounded border border-border bg-background text-foreground placeholder-muted-foreground resize-none h-24"
        />
        {!formData.message.trim() && (
          <p className="text-xs text-red-500 mt-2">
            {t("reminderCreation.details.messageRequired")}
          </p>
        )}
      </div>

      {/* Recurrence Section */}
      <div className="p-4 rounded-lg border border-border bg-secondary/20">
        <div className="flex items-center gap-2 mb-4">
          <Zap className="w-5 h-5 text-accent" />
          <label className="text-sm font-semibold text-foreground">
            {t("reminderCreation.details.recurrence")}
          </label>
        </div>
        <div className="grid grid-cols-2 gap-2">
          {recurrenceOptions.map((option) => (
            <Button
              key={option.value}
              onClick={() => handleRecurrenceChange(option.value)}
              variant={
                formData.recurrence === option.value ? "default" : "outline"
              }
              className={
                formData.recurrence === option.value
                  ? "bg-accent text-accent-foreground"
                  : "border-border text-foreground hover:bg-secondary/50"
              }
            >
              {option.label}
            </Button>
          ))}
        </div>
      </div>
    </div>
  );
}
