import { Plus, Bell, CheckCircle2, Clock } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Header } from "@/components/common/header";
import { Calendar } from "@/components/Calendar";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import {
  remindersService,
  accountService,
  type Reminder,
  type Account,
} from "@/services";
import { Footer } from "@/components/common/footer";

export function RemindersPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [reminders, setReminders] = useState<Reminder[]>([]);
  const [account, setAccount] = useState<Account | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Fetch reminders and account data
  useEffect(() => {
    const fetchData = async () => {
      try {
        setIsLoading(true);
        setError(null);

        // Fetch reminders and account in parallel
        const [fetchedReminders, fetchedAccount] = await Promise.all([
          remindersService.getReminders(),
          accountService.getAccount(),
        ]);

        setReminders(fetchedReminders);
        setAccount(fetchedAccount);
      } catch (err) {
        console.error("Failed to fetch data:", err);
        setError(err instanceof Error ? err.message : "Failed to fetch data");
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  // Calculate statistics
  const totalReminders = reminders.length;
  const activeReminders = reminders.filter((r) => {
    const reminderDate = new Date(r.remind_at_utc);
    return reminderDate > new Date();
  }).length;

  const handleAddReminder = () => {
    navigate("/reminders/create");
  };

  return (
    <div className="min-h-screen bg-background-main dark:bg-background-main">
      <Header />

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Welcome Section */}
        <div className="mb-12">
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 mb-4">
            <div>
              <h2 className="text-3xl sm:text-4xl font-bold text-foreground">
                {t("welcome.title")}
              </h2>
              <p className="text-muted-foreground text-base sm:text-lg mt-2">
                {t("welcome.subtitle")}
              </p>
            </div>
            <Button
              onClick={() => navigate("/reminders/create")}
              className="bg-accent hover:bg-accent/90 text-accent-foreground font-semibold w-full sm:w-auto gap-2"
            >
              <Plus className="w-4 h-4" />
              {t("welcome.newReminder")}
            </Button>
          </div>
        </div>

        {/* Error State */}
        {error && (
          <Card className="border-red-500/50 bg-red-500/10 backdrop-blur mb-6">
            <CardContent className="pt-6">
              <p className="text-red-600 dark:text-red-400">{error}</p>
            </CardContent>
          </Card>
        )}

        {/* Loading State */}
        {isLoading ? (
          <Card className="border-border bg-card/95 backdrop-blur text-center py-12">
            <Clock className="w-12 h-12 text-muted-foreground mx-auto mb-4 animate-spin" />
            <p className="text-muted-foreground">Loading your reminders...</p>
          </Card>
        ) : (
          <>
            {/* Account Overview Cards */}
            {totalReminders > 0 && (
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-12">
                {/* Total Reminders Card */}
                <Card className="border-border bg-card/95 backdrop-blur hover:border-accent/50 transition-colors">
                  <CardHeader className="pb-3">
                    <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                      <Bell className="w-4 h-4" />
                      {t("overview.totalReminders")}
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-3xl font-bold text-foreground">
                      {totalReminders}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      {t("overview.active")}: {activeReminders}
                    </p>
                  </CardContent>
                </Card>

                {/* Active Reminders Card */}
                <Card className="border-border bg-card/95 backdrop-blur hover:border-accent/50 transition-colors">
                  <CardHeader className="pb-3">
                    <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                      <Clock className="w-4 h-4 text-accent" />
                      {t("overview.activeReminders")}
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-3xl font-bold text-accent">
                      {activeReminders}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      {t("overview.accountStatus")}:{" "}
                      <span className="text-accent font-semibold">
                        {t("overview.active")}
                      </span>
                    </p>
                  </CardContent>
                </Card>

                {/* Timezone Card */}
                <Card className="border-border bg-card/95 backdrop-blur hover:border-accent/50 transition-colors">
                  <CardHeader className="pb-3">
                    <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                      <CheckCircle2 className="w-4 h-4" />
                      {t("overview.timezone")}
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-3xl font-bold text-foreground">
                      {account?.timezone || "UTC"}
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      {new Date().toLocaleDateString()}
                    </p>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* Layout: Calendar on left (desktop only), Reminders list on right */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
              {/* Calendar Section - Hidden on mobile */}
              <div className="lg:col-span-2 hidden lg:block">
                <Calendar
                  reminders={reminders}
                  onAddReminder={handleAddReminder}
                />
              </div>

              {/* Reminders List Section */}
              <div className="lg:col-span-1">
                {totalReminders > 0 ? (
                  <div className="space-y-4">
                    <div>
                      <h3 className="text-2xl font-bold text-foreground mb-1">
                        {t("reminders.title")}
                      </h3>
                      <p className="text-muted-foreground text-sm mb-9">
                        {t("reminders.subtitle")}
                      </p>
                    </div>

                    <div className="space-y-3 max-h-[800px] overflow-y-auto pr-2">
                      {reminders.map((reminder) => {
                        const reminderDate = new Date(reminder.remind_at_utc);
                        const isUpcoming = reminderDate > new Date();
                        const isPaused = reminder.is_paused;

                        return (
                          <Card
                            key={reminder.id}
                            className="border-border bg-card/60 backdrop-blur hover:bg-card/80 hover:border-accent/50 transition-all cursor-pointer group"
                            onClick={() =>
                              navigate(`/reminders/${reminder.id}`)
                            }
                          >
                            <CardContent className="p-4">
                              <div className="space-y-3">
                                {/* Header: Message and Status */}
                                <div className="flex items-start justify-between gap-2">
                                  <div className="flex-1 min-w-0">
                                    <p className="text-sm font-semibold text-foreground truncate group-hover:text-accent transition-colors">
                                      {reminder.message}
                                    </p>
                                  </div>
                                  <div className="flex gap-2 flex-shrink-0">
                                    {isPaused && (
                                      <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-yellow-500/20 text-yellow-600 dark:text-yellow-400 border border-yellow-500/30">
                                        {t("reminderDetails.paused")}
                                      </span>
                                    )}
                                    {isUpcoming && !isPaused && (
                                      <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-500/20 text-green-600 dark:text-green-400 border border-green-500/30">
                                        {t("reminderDetails.active")}
                                      </span>
                                    )}
                                  </div>
                                </div>

                                {/* Date and Time */}
                                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                                  <Clock className="w-3 h-3" />
                                  <span>
                                    {reminderDate.toLocaleDateString([], {
                                      month: "short",
                                      day: "numeric",
                                    })}{" "}
                                    at{" "}
                                    {reminderDate.toLocaleTimeString([], {
                                      hour: "2-digit",
                                      minute: "2-digit",
                                    })}
                                  </span>
                                </div>

                                {/* Destinations */}
                                {reminder.destinations &&
                                  reminder.destinations.length > 0 && (
                                    <div className="flex items-center gap-2 flex-wrap">
                                      {reminder.destinations.map(
                                        (dest, idx) => (
                                          <div
                                            key={idx}
                                            className="inline-flex items-center gap-1 px-2 py-1 rounded-md bg-secondary/40 border border-border/50"
                                          >
                                            {dest.type === "discord_dm" && (
                                              <>
                                                <svg
                                                  className="w-3 h-3 text-indigo-500"
                                                  fill="currentColor"
                                                  viewBox="0 0 24 24"
                                                >
                                                  <path d="M20.317 4.37a19.791 19.791 0 00-4.885-1.515.074.074 0 00-.079.037c-.21.375-.444.864-.607 1.25a18.27 18.27 0 00-5.487 0c-.163-.386-.397-.875-.61-1.25a.077.077 0 00-.079-.037A19.736 19.736 0 003.677 4.37a.07.07 0 00-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 00.031.056 19.9 19.9 0 005.993 3.03.078.078 0 00.084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 00-.042-.106 13.107 13.107 0 01-1.872-.892.077.077 0 00-.008-.128 10.713 10.713 0 00.372-.294.075.075 0 00.03-.066c.329-.246.648-.5.954-.76a.07.07 0 00.076-.01 13.697 13.697 0 0011.086 0 .07.07 0 00.076.009c.305.26.625.514.954.759a.077.077 0 00.03.067c.12.088.246.177.371.294a.077.077 0 00-.006.127 13.227 13.227 0 01-1.873.892.076.076 0 00-.041.107c.352.699.764 1.364 1.225 1.994a.076.076 0 00.084.028 19.963 19.963 0 006.002-3.03.077.077 0 00.032-.054c.5-4.817-.838-9.033-3.55-12.765a.061.061 0 00-.031-.03zM8.02 15.33c-1.183 0-2.157-.969-2.157-2.156 0-1.193.974-2.157 2.157-2.157 1.193 0 2.156.964 2.156 2.157 0 1.187-.963 2.156-2.156 2.156zm7.975 0c-1.183 0-2.157-.969-2.157-2.156 0-1.193.974-2.157 2.157-2.157 1.193 0 2.157.964 2.157 2.157 0 1.187-.964 2.156-2.157 2.156z" />
                                                </svg>
                                                <span className="text-xs font-medium text-foreground">
                                                  DM
                                                </span>
                                              </>
                                            )}
                                            {dest.type ===
                                              "discord_channel" && (
                                              <>
                                                <svg
                                                  className="w-3 h-3 text-indigo-500"
                                                  fill="currentColor"
                                                  viewBox="0 0 24 24"
                                                >
                                                  <path d="M20.317 4.37a19.791 19.791 0 00-4.885-1.515.074.074 0 00-.079.037c-.21.375-.444.864-.607 1.25a18.27 18.27 0 00-5.487 0c-.163-.386-.397-.875-.61-1.25a.077.077 0 00-.079-.037A19.736 19.736 0 003.677 4.37a.07.07 0 00-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 00.031.056 19.9 19.9 0 005.993 3.03.078.078 0 00.084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 00-.042-.106 13.107 13.107 0 01-1.872-.892.077.077 0 00-.008-.128 10.713 10.713 0 00.372-.294.075.075 0 00.03-.066c.329-.246.648-.5.954-.76a.07.07 0 00.076-.01 13.697 13.697 0 0011.086 0 .07.07 0 00.076.009c.305.26.625.514.954.759a.077.077 0 00.03.067c.12.088.246.177.371.294a.077.077 0 00-.006.127 13.227 13.227 0 01-1.873.892.076.076 0 00-.041.107c.352.699.764 1.364 1.225 1.994a.076.076 0 00.084.028 19.963 19.963 0 006.002-3.03.077.077 0 00.032-.054c.5-4.817-.838-9.033-3.55-12.765a.061.061 0 00-.031-.03zM8.02 15.33c-1.183 0-2.157-.969-2.157-2.156 0-1.193.974-2.157 2.157-2.157 1.193 0 2.156.964 2.156 2.157 0 1.187-.963 2.156-2.156 2.156zm7.975 0c-1.183 0-2.157-.969-2.157-2.156 0-1.193.974-2.157 2.157-2.157 1.193 0 2.157.964 2.157 2.157 0 1.187-.964 2.156-2.157 2.156z" />
                                                </svg>
                                                <span className="text-xs font-medium text-foreground">
                                                  Channel
                                                </span>
                                              </>
                                            )}
                                            {dest.type === "webhook" && (
                                              <>
                                                <svg
                                                  className="w-3 h-3 text-orange-500"
                                                  fill="currentColor"
                                                  viewBox="0 0 24 24"
                                                >
                                                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" />
                                                </svg>
                                                <span className="text-xs font-medium text-foreground">
                                                  Webhook
                                                </span>
                                              </>
                                            )}
                                          </div>
                                        )
                                      )}
                                    </div>
                                  )}
                              </div>
                            </CardContent>
                          </Card>
                        );
                      })}
                    </div>
                  </div>
                ) : (
                  <>
                    {/* Empty State */}
                    <div className="space-y-4">
                      <div>
                        <h3 className="text-2xl font-bold text-foreground mb-1">
                          {t("reminders.title")}
                        </h3>
                        <p className="text-muted-foreground text-sm mb-9">
                          {t("reminders.subtitle")}
                        </p>
                      </div>
                      <Card className="border-border bg-card/60 backdrop-blur text-center py-12">
                        <Bell className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                        <h3 className="text-lg font-semibold text-foreground mb-2">
                          {t("calendar.noReminders")}
                        </h3>
                        <p className="text-muted-foreground mb-6">
                          {t("calendar.subtitle")}
                        </p>
                        <Button
                          onClick={() => navigate("/reminders/create")}
                          className="bg-accent hover:bg-accent/90 text-accent-foreground font-semibold gap-2"
                        >
                          <Plus className="w-4 h-4" />
                          {t("welcome.newReminder")}
                        </Button>
                      </Card>
                    </div>
                  </>
                )}
              </div>
            </div>
          </>
        )}
      </main>

      {/* Footer */}
      <Footer />
    </div>
  );
}
