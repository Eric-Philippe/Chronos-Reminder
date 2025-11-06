import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  Mail,
  MessageSquare,
  Bug,
  Lightbulb,
  Send,
  CheckCircle,
} from "lucide-react";
import { Header } from "../components/common/header";
import { Footer } from "@/components/common/footer";
import { toast } from "sonner";

type MessageType = "general" | "feedback" | "bug" | "feature";

interface ContactForm {
  name: string;
  email: string;
  type: MessageType;
  subject: string;
  message: string;
}

export function ContactPage() {
  const { t } = useTranslation();
  const [formData, setFormData] = useState<ContactForm>({
    name: "",
    email: "",
    type: "general",
    subject: "",
    message: "",
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);

  const messageTypes = [
    {
      id: "general" as const,
      label: t("contact.general"),
      icon: Mail,
      description: t("contact.generalDesc"),
      color: "text-blue-500",
    },
    {
      id: "feedback" as const,
      label: t("contact.feedback"),
      icon: MessageSquare,
      description: t("contact.feedbackDesc"),
      color: "text-purple-500",
    },
    {
      id: "bug" as const,
      label: t("contact.bug"),
      icon: Bug,
      description: t("contact.bugDesc"),
      color: "text-red-500",
    },
    {
      id: "feature" as const,
      label: t("contact.feature"),
      icon: Lightbulb,
      description: t("contact.featureDesc"),
      color: "text-yellow-500",
    },
  ];

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleTypeSelect = (type: MessageType) => {
    setFormData((prev) => ({
      ...prev,
      type,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validation
    if (!formData.name.trim()) {
      toast.error(t("contact.nameRequired"));
      return;
    }
    if (!formData.email.trim()) {
      toast.error(t("contact.emailRequired"));
      return;
    }
    if (!formData.subject.trim()) {
      toast.error(t("contact.subjectRequired"));
      return;
    }
    if (!formData.message.trim()) {
      toast.error(t("contact.messageRequired"));
      return;
    }

    setIsSubmitting(true);

    try {
      // Get the API URL from environment or use localhost as default
      const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";
      const response = await fetch(`${apiUrl}/api/contact`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        throw new Error("Failed to send message");
      }

      setIsSuccess(true);
      toast.success(t("contact.sentSuccessfully"));

      // Reset form after 2 seconds
      setTimeout(() => {
        setFormData({
          name: "",
          email: "",
          type: "general",
          subject: "",
          message: "",
        });
        setIsSuccess(false);
      }, 2000);
    } catch (error) {
      console.error("Error sending message:", error);
      toast.error(t("contact.sendFailed"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const getTypeInfo = (type: MessageType) => {
    return messageTypes.find((t) => t.id === type);
  };

  const currentType = getTypeInfo(formData.type);

  return (
    <>
      <Header />
      <main className="min-h-screen bg-gradient-to-br from-background to-background-secondary py-12 px-4 sm:px-6 lg:px-8 pt-24">
        <div className="max-w-4xl mx-auto">
          {/* Header Section */}
          <div className="text-center mb-16">
            <h1 className="text-4xl sm:text-5xl font-bold mb-4 text-foreground">
              {t("contact.title")}
            </h1>
            <p className="text-lg text-foreground/70 max-w-2xl mx-auto">
              {t("contact.description")}
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Contact Form */}
            <div className="lg:col-span-2">
              <div className="rounded-2xl border border-border/50 dark:border-white/10 backdrop-blur-sm p-8 bg-white/30 dark:bg-black/20">
                {isSuccess ? (
                  <div className="flex flex-col items-center justify-center py-12 gap-4">
                    <div className="w-16 h-16 rounded-full bg-green-500/20 flex items-center justify-center animate-pulse">
                      <CheckCircle className="w-8 h-8 text-green-500" />
                    </div>
                    <h2 className="text-2xl font-bold text-foreground">
                      {t("contact.successTitle")}
                    </h2>
                    <p className="text-foreground/70 text-center">
                      {t("contact.successMessage")}
                    </p>
                  </div>
                ) : (
                  <form onSubmit={handleSubmit} className="space-y-6">
                    {/* Name and Email */}
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-foreground mb-2">
                          {t("contact.name")} *
                        </label>
                        <input
                          type="text"
                          name="name"
                          value={formData.name}
                          onChange={handleChange}
                          placeholder={t("contact.namePlaceholder")}
                          className="w-full px-4 py-2 rounded-lg bg-background/50 border border-border/30 text-foreground placeholder-foreground/50 focus:outline-none focus:border-amber-500/50 focus:ring-2 focus:ring-amber-500/20 transition-all"
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-foreground mb-2">
                          {t("contact.email")} *
                        </label>
                        <input
                          type="email"
                          name="email"
                          value={formData.email}
                          onChange={handleChange}
                          placeholder={t("contact.emailPlaceholder")}
                          className="w-full px-4 py-2 rounded-lg bg-background/50 border border-border/30 text-foreground placeholder-foreground/50 focus:outline-none focus:border-amber-500/50 focus:ring-2 focus:ring-amber-500/20 transition-all"
                        />
                      </div>
                    </div>

                    {/* Message Type Selection */}
                    <div>
                      <label className="block text-sm font-medium text-foreground mb-3">
                        {t("contact.messageType")} *
                      </label>
                      <div className="grid grid-cols-2 sm:grid-cols-4 gap-2">
                        {messageTypes.map((type) => (
                          <button
                            key={type.id}
                            type="button"
                            onClick={() => handleTypeSelect(type.id)}
                            className={`p-3 rounded-lg border-2 transition-all duration-200 ${
                              formData.type === type.id
                                ? "border-amber-500 bg-amber-500/20"
                                : "border-border/30 dark:border-white/5 hover:border-border/50 dark:hover:border-white/20"
                            }`}
                          >
                            <type.icon
                              className={`w-5 h-5 mx-auto mb-1 ${type.color}`}
                            />
                            <div className="text-xs font-medium text-foreground">
                              {type.label}
                            </div>
                          </button>
                        ))}
                      </div>
                    </div>

                    {/* Subject */}
                    <div>
                      <label className="block text-sm font-medium text-foreground mb-2">
                        {t("contact.subject")} *
                      </label>
                      <input
                        type="text"
                        name="subject"
                        value={formData.subject}
                        onChange={handleChange}
                        placeholder={`${currentType?.label}`}
                        className="w-full px-4 py-2 rounded-lg bg-background/50 border border-border/30 text-foreground placeholder-foreground/50 focus:outline-none focus:border-amber-500/50 focus:ring-2 focus:ring-amber-500/20 transition-all"
                      />
                    </div>

                    {/* Message */}
                    <div>
                      <label className="block text-sm font-medium text-foreground mb-2">
                        {t("contact.message")} *
                      </label>
                      <textarea
                        name="message"
                        value={formData.message}
                        onChange={handleChange}
                        placeholder={t("contact.messagePlaceholder")}
                        rows={6}
                        className="w-full px-4 py-2 rounded-lg bg-background/50 border border-border/30 text-foreground placeholder-foreground/50 focus:outline-none focus:border-amber-500/50 focus:ring-2 focus:ring-amber-500/20 transition-all resize-none"
                      />
                    </div>

                    {/* Submit Button */}
                    <button
                      type="submit"
                      disabled={isSubmitting}
                      className="w-full px-6 py-3 bg-amber-600 hover:bg-amber-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors flex items-center justify-center gap-2"
                    >
                      <Send className="w-4 h-4" />
                      {isSubmitting
                        ? t("contact.sending")
                        : t("contact.sendMessage")}
                    </button>
                  </form>
                )}
              </div>
            </div>

            {/* Contact Info Sidebar */}
            <div className="lg:col-span-1 space-y-4">
              {/* Message Type Info */}
              {currentType && (
                <div className="rounded-2xl border border-border/50 dark:border-white/10 backdrop-blur-sm p-6 bg-white/30 dark:bg-black/20">
                  <div className="flex items-start gap-3 mb-4">
                    <currentType.icon
                      className={`w-6 h-6 ${currentType.color}`}
                    />
                    <div>
                      <h3 className="font-semibold text-foreground">
                        {currentType.label}
                      </h3>
                      <p className="text-sm text-foreground/70">
                        {currentType.description}
                      </p>
                    </div>
                  </div>
                </div>
              )}

              {/* Quick Links */}
              <div className="rounded-2xl border border-border/50 dark:border-white/10 backdrop-blur-sm p-6 bg-white/30 dark:bg-black/20 space-y-3">
                <h3 className="font-semibold text-foreground mb-4">
                  {t("contact.otherWays")}
                </h3>
                <a
                  href="mailto:contact@chronos-reminder.com"
                  className="block p-3 rounded-lg bg-background/50 hover:bg-background/70 border border-border/30 text-foreground text-sm transition-colors"
                >
                  <div className="font-medium">{t("contact.directEmail")}</div>
                  <div className="text-foreground/70 text-xs mt-1">
                    contact@chronos-reminder.com
                  </div>
                </a>

                <div className="h-px bg-gradient-to-r from-transparent via-border/50 to-transparent my-4" />

                <div className="text-xs text-foreground/60 space-y-2">
                  <p>
                    <strong>{t("contact.responseTime")}</strong>
                  </p>
                  <p>
                    <strong>{t("contact.tip")}</strong>
                  </p>
                </div>
              </div>

              {/* All Message Types */}
              <div className="rounded-2xl border border-border/50 dark:border-white/10 backdrop-blur-sm p-6 bg-white/30 dark:bg-black/20 space-y-3">
                <h3 className="font-semibold text-foreground mb-4 text-sm">
                  {t("contact.allMessageTypes")}
                </h3>
                {messageTypes.map((type) => (
                  <button
                    key={type.id}
                    onClick={() => handleTypeSelect(type.id)}
                    className={`w-full text-left p-3 rounded-lg border transition-all ${
                      formData.type === type.id
                        ? "border-amber-500 bg-amber-500/20"
                        : "border-border/30 dark:border-white/5 hover:border-border/50"
                    }`}
                  >
                    <div className="flex items-center gap-2">
                      <type.icon className={`w-4 h-4 ${type.color}`} />
                      <span className="font-medium text-sm text-foreground">
                        {type.label}
                      </span>
                    </div>
                  </button>
                ))}
              </div>
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </>
  );
}
