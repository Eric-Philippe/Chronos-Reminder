import { Bell, Zap, CheckCircle2, Edit3, Trash2, Eye } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { sampleAccount, sampleReminders } from "@/data/sampleData";
import { Header } from "@/components/header";
import type { Reminder } from "@/types/models";

export function WelcomePage() {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getRecurrenceLabel = (recurrence: number) => {
    switch (recurrence) {
      case 0:
        return "Once";
      case 1:
        return "Daily";
      case 2:
        return "Weekly";
      case 3:
        return "Monthly";
      default:
        return "Custom";
    }
  };

  return (
    <div className="min-h-screen" style={{ backgroundColor: "#1a1a18" }}>
      <Header />

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Welcome Section */}
        <div className="mb-12">
          <div className="flex items-center justify-between mb-2">
            <h2 className="text-4xl font-bold text-white">
              Welcome to Chronos
            </h2>
            <Button className="bg-amber-500 hover:bg-amber-600 text-slate-950 font-semibold">
              + New Reminder
            </Button>
          </div>
          <p className="text-slate-400 text-lg">
            Manage and organize all your reminders in one place
          </p>
        </div>

        {/* Account Overview */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-12">
          <Card className="border-slate-800 bg-slate-900/50 backdrop-blur">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-slate-400 flex items-center gap-2">
                <Bell className="w-4 h-4 text-amber-500" />
                Total Reminders
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-white">
                {sampleReminders.length}
              </p>
              <p className="text-xs text-slate-500 mt-1">Active reminders</p>
            </CardContent>
          </Card>

          <Card className="border-slate-800 bg-slate-900/50 backdrop-blur">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-slate-400 flex items-center gap-2">
                <Zap className="w-4 h-4 text-amber-500" />
                Timezone
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-white truncate">
                {sampleAccount.timezone?.name}
              </p>
              <p className="text-xs text-slate-500 mt-1">
                {sampleAccount.timezone?.value}
              </p>
            </CardContent>
          </Card>

          <Card className="border-slate-800 bg-slate-900/50 backdrop-blur">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-slate-400 flex items-center gap-2">
                <CheckCircle2 className="w-4 h-4 text-amber-500" />
                Account Status
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-white">Active</p>
              <p className="text-xs text-slate-500 mt-1">
                Since {new Date(sampleAccount.created_at).toLocaleDateString()}
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Reminders List */}
        <div>
          <div className="mb-6">
            <h3 className="text-2xl font-bold text-white mb-2">
              Your Reminders
            </h3>
            <p className="text-slate-400">
              Manage and edit your scheduled reminders
            </p>
          </div>

          <div className="space-y-4">
            {sampleReminders.map((reminder: Reminder) => (
              <Card
                key={reminder.id}
                className="border-slate-800 bg-slate-900/50 backdrop-blur hover:border-amber-500/30 transition-colors"
              >
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <CardTitle className="text-white text-lg mb-2">
                        {reminder.message}
                      </CardTitle>
                      <div className="flex flex-wrap gap-2">
                        <Badge variant="outline">
                          {getRecurrenceLabel(reminder.recurrence)}
                        </Badge>
                        <Badge
                          variant="outline"
                          className="border-slate-700 text-slate-300"
                        >
                          📍{" "}
                          {reminder.destinations?.[0]?.type === "discord_dm"
                            ? "Discord DM"
                            : "Discord Channel"}
                        </Badge>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        variant="outline"
                        className="border-slate-700 text-slate-300 hover:bg-slate-800/50"
                      >
                        <Eye className="w-4 h-4 mr-1" />
                        View
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        className="border-slate-700 text-slate-300 hover:bg-slate-800/50"
                      >
                        <Edit3 className="w-4 h-4 mr-1" />
                        Edit
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        className="border-red-900/50 text-red-400 hover:bg-red-900/20 hover:border-red-700/50"
                      >
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div>
                      <p className="text-slate-500 text-xs uppercase tracking-wide">
                        Scheduled For
                      </p>
                      <p className="text-white font-medium mt-1">
                        {formatDate(reminder.remind_at_utc)}
                      </p>
                    </div>
                    <div>
                      <p className="text-slate-500 text-xs uppercase tracking-wide">
                        Created
                      </p>
                      <p className="text-white font-medium mt-1">
                        {formatDate(reminder.created_at)}
                      </p>
                    </div>
                    <div>
                      <p className="text-slate-500 text-xs uppercase tracking-wide">
                        Destination
                      </p>
                      <p className="text-amber-400 font-medium mt-1">
                        {reminder.destinations?.[0]?.destination || "N/A"}
                      </p>
                    </div>
                    <div>
                      <p className="text-slate-500 text-xs uppercase tracking-wide">
                        Status
                      </p>
                      <p className="text-green-400 font-medium mt-1">Pending</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>

        {/* Empty State Hint */}
        {sampleReminders.length === 0 && (
          <Card className="border-slate-800 bg-slate-900/50 backdrop-blur text-center py-12">
            <Bell className="w-12 h-12 text-slate-600 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-slate-300 mb-2">
              No reminders yet
            </h3>
            <p className="text-slate-500 mb-4">
              Create your first reminder to get started
            </p>
            <Button className="bg-amber-500 hover:bg-amber-600 text-slate-950 font-semibold">
              Create Reminder
            </Button>
          </Card>
        )}
      </main>

      {/* Footer */}
      <footer
        className="border-t border-stone-800 mt-16"
        style={{ backgroundColor: "#1a1a18" }}
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex flex-col md:flex-row items-center justify-between text-slate-400 text-sm">
            <p>&copy; 2025 Chronos. Never miss a moment.</p>
            <div className="flex gap-6 mt-4 md:mt-0">
              <a href="#" className="hover:text-amber-400 transition-colors">
                Documentation
              </a>
              <a href="#" className="hover:text-amber-400 transition-colors">
                Support
              </a>
              <a href="#" className="hover:text-amber-400 transition-colors">
                Discord
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
