import { useEffect, useRef } from "react";

interface ClockProps {
  datetime: Date;
  timezone?: string;
  label?: string;
  size?: "xs" | "sm" | "md" | "lg";
}

export function Clock({
  datetime,
  timezone = "UTC",
  label,
  size = "md",
}: ClockProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  const sizeConfig = {
    xs: { canvas: 80, radius: 33, font: "10px" },
    sm: { canvas: 100, radius: 42, font: "12px" },
    md: { canvas: 200, radius: 85, font: "16px" },
    lg: { canvas: 280, radius: 120, font: "20px" },
  };

  const config = sizeConfig[size];

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    // Get time in specific timezone
    const formatter = new Intl.DateTimeFormat("en-US", {
      timeZone: timezone,
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    });

    const parts = formatter.formatToParts(datetime);
    const timeMap = Object.fromEntries(
      parts.map((part) => [part.type, part.value])
    );

    const hours = parseInt(timeMap.hour) % 12;
    const minutes = parseInt(timeMap.minute);
    const seconds = parseInt(timeMap.second);

    const radius = config.radius;
    const centerX = canvas.width / 2;
    const centerY = canvas.height / 2;

    // Clear canvas with semi-transparent background
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Draw elegant glassmorphism background
    const gradient = ctx.createLinearGradient(
      0,
      0,
      canvas.width,
      canvas.height
    );
    gradient.addColorStop(0, "rgba(255, 255, 255, 0.08)");
    gradient.addColorStop(1, "rgba(255, 255, 255, 0.02)");
    ctx.fillStyle = gradient;
    ctx.beginPath();
    ctx.arc(centerX, centerY, radius, 0, Math.PI * 2);
    ctx.fill();

    // Draw ultra-thin primary border
    ctx.strokeStyle = "rgba(255, 255, 255, 0.25)";
    ctx.lineWidth = 1.5;
    ctx.stroke();

    // Draw delicate gold accent ring
    ctx.strokeStyle = "rgba(218, 165, 32, 0.3)";
    ctx.lineWidth = 1;
    ctx.beginPath();
    ctx.arc(centerX, centerY, radius - 3, 0, Math.PI * 2);
    ctx.stroke();

    // Draw numbers (12, 3, 6, 9) - minimalist style
    ctx.fillStyle = "rgba(255, 255, 255, 0.65)";
    ctx.font = `light ${config.font} -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif`;
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";

    const numberRadius = radius - 28;
    for (let i = 0; i < 12; i += 3) {
      const angle = (i * Math.PI) / 6 - Math.PI / 2;
      const x = centerX + Math.cos(angle) * numberRadius;
      const y = centerY + Math.sin(angle) * numberRadius;
      ctx.fillText(i === 0 ? "12" : String(i), x, y);
    }

    // Draw hour markers with refined luxury style
    for (let i = 0; i < 12; i++) {
      const angle = (i * Math.PI) / 6 - Math.PI / 2;
      const x1 = centerX + Math.cos(angle) * (radius - 12);
      const y1 = centerY + Math.sin(angle) * (radius - 12);
      const x2 = centerX + Math.cos(angle) * (radius - 20);
      const y2 = centerY + Math.sin(angle) * (radius - 20);

      // Gold accents only on cardinal points (12, 3, 6, 9)
      if (i % 3 === 0) {
        ctx.strokeStyle = "rgba(218, 165, 32, 0.7)";
        ctx.lineWidth = 2;
      } else {
        ctx.strokeStyle = "rgba(255, 255, 255, 0.35)";
        ctx.lineWidth = 1;
      }

      ctx.beginPath();
      ctx.moveTo(x1, y1);
      ctx.lineTo(x2, y2);
      ctx.stroke();
    }

    // Draw hour hand - elegant and thin
    const hourAngle =
      (hours * Math.PI) / 6 + (minutes * Math.PI) / (6 * 60) - Math.PI / 2;
    const hourLength = radius * 0.45;
    ctx.strokeStyle = "rgba(255, 255, 255, 0.85)";
    ctx.lineWidth = 4;
    ctx.lineCap = "round";
    ctx.beginPath();
    ctx.moveTo(centerX, centerY);
    ctx.lineTo(
      centerX + Math.cos(hourAngle) * hourLength,
      centerY + Math.sin(hourAngle) * hourLength
    );
    ctx.stroke();

    // Draw minute hand - refined
    const minuteAngle =
      (minutes * Math.PI) / 30 + (seconds * Math.PI) / (30 * 60) - Math.PI / 2;
    const minuteLength = radius * 0.65;
    ctx.strokeStyle = "rgba(255, 255, 255, 0.75)";
    ctx.lineWidth = 3;
    ctx.beginPath();
    ctx.moveTo(centerX, centerY);
    ctx.lineTo(
      centerX + Math.cos(minuteAngle) * minuteLength,
      centerY + Math.sin(minuteAngle) * minuteLength
    );
    ctx.stroke();

    // Draw second hand - subtle gold
    const secondAngle = (seconds * Math.PI) / 30 - Math.PI / 2;
    const secondLength = radius * 0.7;
    ctx.strokeStyle = "rgba(218, 165, 32, 0.8)";
    ctx.lineWidth = 1.5;
    ctx.beginPath();
    ctx.moveTo(centerX, centerY);
    ctx.lineTo(
      centerX + Math.cos(secondAngle) * secondLength,
      centerY + Math.sin(secondAngle) * secondLength
    );
    ctx.stroke();

    // Draw center dot with refined luxury style
    ctx.fillStyle = "rgba(255, 255, 255, 0.85)";
    ctx.beginPath();
    ctx.arc(centerX, centerY, 5, 0, Math.PI * 2);
    ctx.fill();

    // Draw subtle gold ring around center
    ctx.strokeStyle = "rgba(218, 165, 32, 0.6)";
    ctx.lineWidth = 1.5;
    ctx.stroke();
  }, [datetime, timezone, size, config]);

  return (
    <div className="flex flex-col items-center gap-3">
      <div className="backdrop-blur-2xl bg-gradient-to-br from-white/8 to-white/2 border border-white/20 rounded-3xl p-4 shadow-xl hover:shadow-2xl hover:shadow-yellow-400/10 transition-all duration-500 hover:scale-105 hover:border-white/40 group relative overflow-hidden">
        {/* Elegant animated background glow on hover */}
        <div className="absolute inset-0 bg-gradient-to-r from-yellow-400/0 via-yellow-400/4 to-yellow-400/0 opacity-0 group-hover:opacity-100 transition-opacity duration-500 rounded-3xl"></div>
        <canvas
          ref={canvasRef}
          width={config.canvas}
          height={config.canvas}
          className="block relative z-10"
        />
      </div>
      {label && (
        <div className="text-center">
          <p className="text-sm font-semibold text-white/80">{label}</p>
          <p className="text-xs text-white/50">{timezone}</p>
        </div>
      )}
    </div>
  );
}
