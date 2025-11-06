import { useEffect, useState } from "react";

interface DigitalClockProps {
  datetime: Date;
  timezone?: string;
  label?: string;
  format?: "12h" | "24h";
  size?: "sm" | "md";
}

export function DigitalClock({
  datetime,
  timezone = "UTC",
  label,
  format = "24h",
  size = "md",
}: DigitalClockProps) {
  const [time, setTime] = useState("");
  const [seconds, setSeconds] = useState("");

  const sizeClasses = {
    sm: {
      container: "px-3 py-1.5 rounded",
      time: "text-xs sm:text-sm",
      seconds: "text-xs",
      gap: "gap-0.5",
    },
    md: {
      container: "px-4 py-2 rounded-md",
      time: "text-sm sm:text-base",
      seconds: "text-xs",
      gap: "gap-0.5",
    },
  };

  useEffect(() => {
    const formatter = new Intl.DateTimeFormat("en-US", {
      timeZone: timezone,
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: format === "12h",
    });

    const parts = formatter.formatToParts(datetime);
    const timeMap = Object.fromEntries(
      parts.map((part) => [part.type, part.value])
    );

    const hour = timeMap.hour || "00";
    const minute = timeMap.minute || "00";
    const second = timeMap.second || "00";
    const period = timeMap.dayPeriod || "";

    const timeStr = `${hour}:${minute}`;
    const secondsStr = second;
    const fullTime = period ? `${timeStr} ${period}` : timeStr;

    setTime(fullTime);
    setSeconds(secondsStr);
  }, [datetime, timezone, format]);

  return (
    <div className="flex flex-col items-center gap-2">
      {/* Digital Clock Display */}
      <div
        className={`backdrop-blur-2xl bg-gradient-to-br from-white/10 to-white/3 border border-white/20 ${sizeClasses[size].container} shadow-lg hover:shadow-xl hover:shadow-yellow-400/5 transition-all duration-500 hover:scale-110 hover:border-yellow-400/30 group relative overflow-hidden`}
      >
        {/* Subtle animated background glow */}
        <div className="absolute inset-0 bg-gradient-to-r from-yellow-400/0 via-yellow-400/3 to-yellow-400/0 opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>

        <div className="relative z-10">
          {/* Main time display - luxury minimalist */}
          <div className={`flex items-baseline ${sizeClasses[size].gap}`}>
            <div
              className={`font-light tracking-wide text-white font-mono ${sizeClasses[size].time}`}
            >
              {time}
            </div>
            <div
              className={`font-light text-white/50 font-mono ${sizeClasses[size].seconds}`}
            >
              {seconds}
            </div>
          </div>

          {/* Decorative accent line - subtle */}
          <div className="absolute -bottom-1.5 left-0 right-0 h-px bg-gradient-to-r from-transparent via-yellow-400/40 to-transparent rounded-full blur-sm"></div>
        </div>
      </div>

      {/* Label and Timezone */}
      {label && (
        <div className="text-center">
          <p className="text-sm font-semibold text-white/80">{label}</p>
          <p className="text-xs text-white/50">{timezone}</p>
        </div>
      )}
    </div>
  );
}
