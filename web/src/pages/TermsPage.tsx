import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import { Header } from "@/components/common/header";
import { Footer } from "@/components/common/footer";
import {
  Shield,
  Lock,
  Database,
  Heart,
  Server,
  Trash2,
  HelpCircle,
  ExternalLink,
  Package,
} from "lucide-react";
import { ROUTES } from "@/config/routes";

interface SectionProps {
  icon: React.ReactNode;
  title: string;
  children: React.ReactNode;
}

function Section({ icon, title, children }: SectionProps) {
  return (
    <div className="rounded-2xl border border-border/50 dark:border-white/10 backdrop-blur-sm p-8 bg-white/30 dark:bg-black/20">
      <div className="flex items-center gap-3 mb-4">
        <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-amber-500/15 flex-shrink-0">
          {icon}
        </div>
        <h2 className="text-xl font-bold text-foreground">{title}</h2>
      </div>
      <div className="space-y-3 text-foreground/80 leading-relaxed">
        {children}
      </div>
    </div>
  );
}

export function TermsPage() {
  const { t } = useTranslation();

  const tableRows = [
    {
      data: t("terms.dataTableRow1Data"),
      why: t("terms.dataTableRow1Why"),
      access: t("terms.dataTableRow1Access"),
    },
    {
      data: t("terms.dataTableRow2Data"),
      why: t("terms.dataTableRow2Why"),
      access: t("terms.dataTableRow2Access"),
    },
    {
      data: t("terms.dataTableRow3Data"),
      why: t("terms.dataTableRow3Why"),
      access: t("terms.dataTableRow3Access"),
      highlight: true,
    },
    {
      data: t("terms.dataTableRow4Data"),
      why: t("terms.dataTableRow4Why"),
      access: t("terms.dataTableRow4Access"),
    },
    {
      data: t("terms.dataTableRow5Data"),
      why: t("terms.dataTableRow5Why"),
      access: t("terms.dataTableRow5Access"),
    },
    {
      data: t("terms.dataTableRow6Data"),
      why: t("terms.dataTableRow6Why"),
      access: t("terms.dataTableRow6Access"),
    },
  ];

  return (
    <>
      <Header />
      <main className="min-h-screen bg-gradient-to-br from-background to-background-secondary py-12 px-4 sm:px-6 lg:px-8 pt-24">
        <div className="max-w-3xl mx-auto">
          {/* Header */}
          <div className="text-center mb-14">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-amber-500/15 mb-6">
              <Shield className="w-8 h-8 text-amber-500" />
            </div>
            <h1 className="text-4xl sm:text-5xl font-bold mb-4 text-foreground">
              {t("terms.title")}
            </h1>
            <p className="text-lg text-foreground/70 max-w-xl mx-auto">
              {t("terms.subtitle")}
            </p>
            <p className="text-sm text-foreground/40 mt-3">
              {t("terms.lastUpdated")}
            </p>
          </div>

          {/* Data at a glance table */}
          <div className="rounded-2xl border border-amber-500/20 bg-amber-500/5 backdrop-blur-sm p-6 mb-6">
            <h2 className="text-base font-semibold text-amber-600 dark:text-amber-400 mb-4 uppercase tracking-wide">
              {t("terms.dataTableTitle")}
            </h2>
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/30">
                    <th className="text-left py-2 pr-4 text-foreground/60 font-medium">
                      {t("terms.dataTableColData")}
                    </th>
                    <th className="text-left py-2 pr-4 text-foreground/60 font-medium">
                      {t("terms.dataTableColWhy")}
                    </th>
                    <th className="text-left py-2 text-foreground/60 font-medium">
                      {t("terms.dataTableColAccess")}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {tableRows.map((row, i) => (
                    <tr
                      key={i}
                      className={`border-b border-border/10 last:border-0 ${
                        row.highlight ? "bg-green-500/5" : ""
                      }`}
                    >
                      <td className="py-2 pr-4 font-medium text-foreground">
                        {row.data}
                      </td>
                      <td className="py-2 pr-4 text-foreground/70">
                        {row.why}
                      </td>
                      <td className="py-2">
                        <span
                          className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
                            row.access.startsWith("Nobody") ||
                            row.access.startsWith("Personne") ||
                            row.access.startsWith("Nadie")
                              ? "bg-green-500/15 text-green-600 dark:text-green-400"
                              : "bg-secondary/50 text-foreground/70"
                          }`}
                        >
                          {row.access}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Sections */}
          <div className="space-y-4">
            <Section
              icon={<Heart className="w-5 h-5 text-amber-500" />}
              title={t("terms.section1Title")}
            >
              <p>{t("terms.section1Body")}</p>
            </Section>

            <Section
              icon={<Database className="w-5 h-5 text-amber-500" />}
              title={t("terms.section2Title")}
            >
              <p>{t("terms.section2Intro")}</p>
              <ul className="space-y-1.5 ml-2">
                {[
                  t("terms.section2Item1"),
                  t("terms.section2Item2"),
                  t("terms.section2Item3"),
                  t("terms.section2Item4"),
                  t("terms.section2Item5"),
                  t("terms.section2Item6"),
                ].map((item, i) => (
                  <li key={i} className="flex items-start gap-2">
                    <span className="mt-1.5 w-1.5 h-1.5 rounded-full bg-amber-500 flex-shrink-0" />
                    <span>{item}</span>
                  </li>
                ))}
              </ul>
              <p className="text-sm text-foreground/60 italic">
                {t("terms.section2Outro")}
              </p>
            </Section>

            <Section
              icon={<Lock className="w-5 h-5 text-green-500" />}
              title={t("terms.section3Title")}
            >
              <div className="p-4 rounded-xl bg-green-500/10 border border-green-500/20">
                <p>{t("terms.section3Body")}</p>
              </div>
            </Section>

            <Section
              icon={<Shield className="w-5 h-5 text-amber-500" />}
              title={t("terms.section4Title")}
            >
              <p>{t("terms.section4Body")}</p>
            </Section>

            <Section
              icon={<Heart className="w-5 h-5 text-pink-500" />}
              title={t("terms.section5Title")}
            >
              <p>{t("terms.section5Body")}</p>
              <p className="text-sm text-foreground/60 italic border-l-2 border-amber-500/30 pl-3">
                {t("terms.section5Aside")}
              </p>
            </Section>

            <Section
              icon={<Server className="w-5 h-5 text-amber-500" />}
              title={t("terms.section6Title")}
            >
              <p>{t("terms.section6Body")}</p>
              <a
                href="https://status.chronosrmd.com/status/chronos"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1.5 text-sm text-amber-600 dark:text-amber-400 hover:underline"
              >
                {t("terms.section6StatusLink")}
                <ExternalLink className="w-3.5 h-3.5" />
              </a>
            </Section>

            <Section
              icon={<Trash2 className="w-5 h-5 text-amber-500" />}
              title={t("terms.section7Title")}
            >
              <p>{t("terms.section7Body")}</p>
              <p className="text-sm text-foreground/60 italic border-l-2 border-border/50 pl-3">
                {t("terms.section7GdprNote")}
              </p>
            </Section>

            <Section
              icon={<Package className="w-5 h-5 text-amber-500" />}
              title={t("terms.section8Title")}
            >
              <p>{t("terms.section8Body")}</p>
              <Link
                to={ROUTES.SELFHOST.path}
                className="inline-flex items-center gap-1.5 text-sm text-amber-600 dark:text-amber-400 hover:underline"
              >
                {t("terms.section8Link")}
              </Link>
            </Section>

            <Section
              icon={<HelpCircle className="w-5 h-5 text-amber-500" />}
              title={t("terms.section9Title")}
            >
              <p>{t("terms.section9Body")}</p>
              <div className="flex flex-wrap gap-3">
                <Link
                  to={ROUTES.CONTACT.path}
                  className="inline-flex items-center gap-1.5 text-sm text-amber-600 dark:text-amber-400 hover:underline"
                >
                  {t("terms.section9ContactLink")}
                </Link>
                <span className="text-foreground/30">·</span>
                <a
                  href="https://discord.gg/m3MsM922QD"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1.5 text-sm text-amber-600 dark:text-amber-400 hover:underline"
                >
                  {t("terms.section9DiscordLink")}
                  <ExternalLink className="w-3.5 h-3.5" />
                </a>
              </div>
            </Section>
          </div>
        </div>
      </main>
      <Footer />
    </>
  );
}
