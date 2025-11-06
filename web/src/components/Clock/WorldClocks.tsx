import { useEffect, useState } from "react";
import { Clock } from "./Clock";
import { DigitalClock } from "./DigitalClock";

interface WorldClockProps {
  title?: string;
  subtitle?: string;
}

const TIMEZONES = [
  { label: "Toronto", timezone: "America/Toronto" },
  { label: "London (UTC)", timezone: "Europe/London" },
  { label: "Paris", timezone: "Europe/Paris" },
  { label: "Tokyo", timezone: "Asia/Tokyo" },
  { label: "Adelaide", timezone: "Australia/Adelaide" },
];

export function WorldClocks({
  title = "World Clocks",
  subtitle = "Time around the globe",
}: WorldClockProps) {
  const [currentTime, setCurrentTime] = useState(new Date());
  const [isMobile, setIsMobile] = useState(window.innerWidth < 768);

  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth < 768);
    };

    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  const getDateString = (timezone: string) => {
    const formatter = new Intl.DateTimeFormat("en-US", {
      timeZone: timezone,
      month: "short",
      day: "numeric",
      weekday: "short",
    });
    return formatter.format(currentTime);
  };

  return (
    <section className="py-20 px-4 sm:px-6 lg:px-8 relative overflow-hidden">
      {/* Enhanced background glassmorphism */}
      <div className="absolute inset-0 backdrop-blur-3xl pointer-events-none -z-10"></div>

      {/* Decorative background elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none -z-10">
        <div className="absolute top-1/2 -left-32 w-64 h-64 bg-accent/5 rounded-full blur-3xl"></div>
        <div className="absolute bottom-1/4 -right-32 w-64 h-64 bg-accent/5 rounded-full blur-3xl"></div>
      </div>

      <div className="max-w-7xl mx-auto relative z-10">
        {/* Separator */}
        <div className="flex items-center gap-4 mb-12 justify-center">
          <div className="h-px flex-1 max-w-xs bg-gradient-to-r from-transparent via-accent/40 to-transparent"></div>
          <div className="w-2 h-2 rounded-full bg-accent/60"></div>
          <div className="h-px flex-1 max-w-xs bg-gradient-to-l from-transparent via-accent/40 to-transparent"></div>
        </div>

        {/* Header */}
        <div className="text-center mb-16">
          <h2 className="text-4xl sm:text-5xl font-bold text-foreground mb-4">
            {title}
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            {subtitle}
          </p>
        </div>

        {/* Clocks Grid - Responsive */}
        <div className="flex flex-wrap justify-center items-start gap-3 md:gap-4 lg:gap-6">
          {TIMEZONES.map((item) => (
            <div
              key={item.timezone}
              className="flex flex-col items-center justify-center"
            >
              <Clock
                datetime={currentTime}
                timezone={item.timezone}
                label={isMobile ? undefined : item.label}
                size={isMobile ? "xs" : "sm"}
              />
              <div className="mt-2">
                <DigitalClock
                  datetime={currentTime}
                  timezone={item.timezone}
                  format="24h"
                  size={isMobile ? "sm" : "md"}
                />
              </div>
              {isMobile && (
                <p className="text-xs text-white/50 mt-1">{item.label}</p>
              )}
              {!isMobile && (
                <p className="text-xs text-white/50 mt-1">
                  {getDateString(item.timezone)}
                </p>
              )}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
