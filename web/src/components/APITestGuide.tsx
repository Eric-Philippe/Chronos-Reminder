import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Alert } from "@/components/ui/alert";
import { Copy, Check, Terminal, Code2, Zap, AlertCircle } from "lucide-react";
import { useToast } from "@/hooks/useToast";

type TestMethod = "curl" | "javascript" | "python";

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

export function APITestGuide() {
  const { t } = useTranslation();
  const toast = useToast();
  const [selectedMethod, setSelectedMethod] = useState<TestMethod>("curl");
  const [copiedCommand, setCopiedCommand] = useState<string | null>(null);

  // Live test state
  const [liveApiKey, setLiveApiKey] = useState("");
  const [liveResult, setLiveResult] = useState<any>(null);
  const [liveLoading, setLiveLoading] = useState(false);
  const [liveError, setLiveError] = useState<string | null>(null);

  const handleLiveTest = async () => {
    setLiveLoading(true);
    setLiveError(null);
    setLiveResult(null);
    try {
      const response = await fetch(`${API_BASE_URL}/api/reminders`, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${liveApiKey}`,
          "Content-Type": "application/json",
        },
      });
      const data = await response.json();
      if (!response.ok) {
        setLiveError(data?.error || `Error: ${response.status}`);
      } else {
        setLiveResult(data);
      }
    } catch (err: any) {
      setLiveError(err?.message || "Network error");
    } finally {
      setLiveLoading(false);
    }
  };

  const testMethods: {
    id: TestMethod;
    label: string;
    icon: React.ReactNode;
  }[] = [
    {
      id: "curl",
      label: t("apiKeys.curlExample"),
      icon: <Terminal className="w-4 h-4" />,
    },
    {
      id: "javascript",
      label: t("apiKeys.javascriptExample"),
      icon: <Code2 className="w-4 h-4" />,
    },
    {
      id: "python",
      label: t("apiKeys.pythonExample"),
      icon: <Code2 className="w-4 h-4" />,
    },
  ];

  const getCommand = (method: TestMethod): string => {
    const apiKey = "YOUR_API_KEY";

    switch (method) {
      case "curl":
        return `curl -X GET ${API_BASE_URL}/api/reminders \\
  -H "Authorization: Bearer ${apiKey}" \\
  -H "Content-Type: application/json"`;

      case "javascript":
        return `const apiKey = "${apiKey}";

fetch("${API_BASE_URL}/api/reminders", {
  method: "GET",
  headers: {
    "Authorization": \`Bearer \${apiKey}\`,
    "Content-Type": "application/json"
  }
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error("Error:", error));`;

      case "python":
        return `import requests

api_key = "${apiKey}"
headers = {
    "Authorization": f"Bearer {api_key}",
    "Content-Type": "application/json"
}

response = requests.get(
    "${API_BASE_URL}/api/reminders",
    headers=headers
)

print(response.json())`;

      default:
        return "";
    }
  };

  const handleCopyCommand = async () => {
    const command = getCommand(selectedMethod);
    try {
      await navigator.clipboard.writeText(command);
      setCopiedCommand(selectedMethod);
      toast.success(t("apiKeys.commandCopied"), {
        description: t("apiKeys.copyCommand"),
      });
      setTimeout(() => setCopiedCommand(null), 2000);
    } catch {
      toast.error(t("common.error"), {
        description: t("common.copyFailed"),
      });
    }
  };

  const sampleResponse = {
    count: 1,
    reminders: [
      {
        id: "00000000-0000-0000-0000-000000000001",
        account_id: "11111111-1111-1111-1111-111111111111",
        remind_at_utc: "2026-05-13T08:00:00+02:00",
        next_fire_utc: "2026-05-13T08:00:00+02:00",
        message: "Eric's Birthday",
        created_at: "2025-11-04T14:09:03.249548+01:00",
        recurrence_type: "DAILY",
        is_paused: false,
        destinations: [
          {
            id: "22222222-2222-2222-2222-222222222222",
            reminder_id: "00000000-0000-0000-0000-000000000001",
            type: "discord_dm",
            metadata: {
              user_id: "483939292948382384",
            },
          },
        ],
      },
    ],
  };

  return (
    <div className="space-y-6">
      {/* Introduction Card */}
      <Card className="border-blue-500/30 bg-blue-500/5">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-blue-600 dark:text-blue-400">
            <Zap className="w-5 h-5" />
            {t("apiKeys.testGuide")}
          </CardTitle>
          <p className="text-sm text-muted-foreground mt-2">
            {t("apiKeys.testGuideDesc")}
          </p>
        </CardHeader>
      </Card>

      {/* Quick Steps */}
      <div className="grid gap-4 md:grid-cols-3">
        <div className="p-6 rounded-lg bg-secondary/20 border border-border">
          <div className="flex items-center gap-3 mb-3">
            <div className="flex items-center justify-center w-8 h-8 rounded-full bg-accent text-accent-foreground text-sm font-bold">
              1
            </div>
            <p className="font-semibold text-foreground text-lg">
              {t("apiKeys.testStep1")}
            </p>
          </div>
          <p className="text-sm text-muted-foreground">
            {t("apiKeys.testStep1Desc")}
          </p>
        </div>

        <div className="p-6 rounded-lg bg-secondary/20 border border-border">
          <div className="flex items-center gap-3 mb-3">
            <div className="flex items-center justify-center w-8 h-8 rounded-full bg-accent text-accent-foreground text-sm font-bold">
              2
            </div>
            <p className="font-semibold text-foreground text-lg">
              {t("apiKeys.testStep2")}
            </p>
          </div>
          <p className="text-sm text-muted-foreground">
            {t("apiKeys.testStep2Desc")}
          </p>
        </div>

        <div className="p-6 rounded-lg bg-secondary/20 border border-border">
          <div className="flex items-center gap-3 mb-3">
            <div className="flex items-center justify-center w-8 h-8 rounded-full bg-accent text-accent-foreground text-sm font-bold">
              3
            </div>
            <p className="font-semibold text-foreground text-lg">
              {t("apiKeys.testStep3")}
            </p>
          </div>
          <p className="text-sm text-muted-foreground">
            {t("apiKeys.testStep3Desc")}
          </p>
        </div>
      </div>

      {/* Method Selection */}
      <Card className="border-border bg-card/95 backdrop-blur">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Code2 className="w-5 h-5" />
            {t("apiKeys.testStep2")}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-3 sm:grid-cols-3">
            {testMethods.map((method) => (
              <Button
                key={method.id}
                onClick={() => setSelectedMethod(method.id)}
                variant={selectedMethod === method.id ? "default" : "outline"}
                className="gap-2 h-12 text-base font-semibold"
              >
                {method.icon}
                <span>{method.label}</span>
              </Button>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Live API Key Test */}
      <Card className="border-green-500/30 bg-green-500/5">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-green-600 dark:text-green-400">
            <Zap className="w-5 h-5" />
            {t("apiKeys.liveTest") || "Live API Key Test"}
          </CardTitle>
          <p className="text-sm text-muted-foreground mt-2">
            {t("apiKeys.liveTestDesc") ||
              "Paste your API key below and test the actual API response."}
          </p>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-3">
            <input
              type="text"
              placeholder={
                t("apiKeys.keyPastePlaceholder") || "Paste your API key here"
              }
              value={liveApiKey}
              onChange={(e) => setLiveApiKey(e.target.value)}
              className="p-2 border rounded bg-background text-foreground font-mono text-sm"
              disabled={liveLoading}
              autoComplete="off"
            />
            <Button
              onClick={handleLiveTest}
              disabled={!liveApiKey.trim() || liveLoading}
              className="w-fit bg-accent text-accent-foreground"
            >
              {liveLoading ? (
                <span className="flex items-center gap-2">
                  <Zap className="w-4 h-4 animate-spin" />{" "}
                  {t("common.loading") || "Testing..."}
                </span>
              ) : (
                <span className="flex items-center gap-2">
                  <Zap className="w-4 h-4" />{" "}
                  {t("apiKeys.testNow") || "Test Now"}
                </span>
              )}
            </Button>
            {liveError && (
              <Alert className="border-red-500/50 bg-red-500/10 mt-2">
                <AlertCircle className="h-4 w-4 text-red-600 dark:text-red-400" />
                <div className="ml-3">
                  <p className="text-red-600 dark:text-red-400 font-mono text-xs">
                    {liveError}
                  </p>
                </div>
              </Alert>
            )}
            {liveResult && (
              <div className="bg-black/80 rounded-lg p-4 font-mono text-xs overflow-x-auto mt-2 border border-green-500/30">
                <pre className="text-green-400">
                  {JSON.stringify(liveResult, null, 2)}
                </pre>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Code Display */}
      <Card className="border-border bg-card/95 backdrop-blur">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center gap-2">
                <Terminal className="w-5 h-5" />
                {t("apiKeys.testEndpoint")}
              </CardTitle>
              <p className="text-sm text-muted-foreground mt-1">
                {t("apiKeys.listReminders")}
              </p>
            </div>
            <Button
              onClick={handleCopyCommand}
              size="sm"
              className="gap-2"
              variant="outline"
            >
              {copiedCommand === selectedMethod ? (
                <>
                  <Check className="w-4 h-4 text-green-600" />
                  {t("common.copied")}
                </>
              ) : (
                <>
                  <Copy className="w-4 h-4" />
                  {t("apiKeys.copyCommand")}
                </>
              )}
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="bg-black/80 rounded-lg p-6 font-mono text-base overflow-x-auto">
            <div className="text-green-400 whitespace-pre-wrap break-words">
              {getCommand(selectedMethod)
                .split("\n")
                .map((line, i) => (
                  <div key={i}>
                    <span className="text-green-400">{line}</span>
                  </div>
                ))}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Instructions Alert */}
      <Alert className="border-amber-500/50 bg-amber-500/10">
        <AlertCircle className="h-4 w-4 text-amber-600 dark:text-amber-400" />
        <div className="ml-3">
          <p className="text-sm text-amber-600 dark:text-amber-400 font-medium">
            {t("common.important")}
          </p>
          <ul className="text-xs text-amber-600/80 dark:text-amber-400/80 mt-2 space-y-1 ml-4 list-disc">
            <li>
              {t("apiKeys.importantNote1").replace("YOUR_API_KEY", "")}
              <code className="bg-amber-500/20 px-1 rounded">YOUR_API_KEY</code>
              {t("apiKeys.importantNote1").substring(
                t("apiKeys.importantNote1").indexOf("with")
              )}
            </li>
            <li>{t("apiKeys.importantNote2")}</li>
            <li>{t("apiKeys.importantNote3")}</li>
          </ul>
        </div>
      </Alert>

      {/* Sample Response */}
      <Card className="border-border bg-card/95 backdrop-blur">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Zap className="w-5 h-5 text-yellow-500" />
            {t("apiKeys.expectedResponse")}
          </CardTitle>
          <p className="text-sm text-muted-foreground mt-1">
            {t("apiKeys.successResponse")}
          </p>
        </CardHeader>
        <CardContent>
          <div className="bg-black/80 rounded-lg p-6 font-mono text-sm overflow-x-auto">
            <pre className="text-green-400">
              {JSON.stringify(sampleResponse, null, 2)}
            </pre>
          </div>
        </CardContent>
      </Card>

      {/* Security Tips */}
      <Card className="border-green-500/30 bg-green-500/5">
        <CardHeader>
          <CardTitle className="text-green-600 dark:text-green-400">
            ✓ {t("apiKeys.security")}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <div className="flex gap-2 text-sm">
            <span className="text-green-600 dark:text-green-400">✓</span>
            <p className="text-muted-foreground">{t("apiKeys.securityTip1")}</p>
          </div>
          <div className="flex gap-2 text-sm">
            <span className="text-green-600 dark:text-green-400">✓</span>
            <p className="text-muted-foreground">{t("apiKeys.securityTip2")}</p>
          </div>
          <div className="flex gap-2 text-sm">
            <span className="text-green-600 dark:text-green-400">✓</span>
            <p className="text-muted-foreground">{t("apiKeys.securityTip3")}</p>
          </div>
          <div className="flex gap-2 text-sm">
            <span className="text-green-600 dark:text-green-400">✓</span>
            <p className="text-muted-foreground">{t("apiKeys.securityTip4")}</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
