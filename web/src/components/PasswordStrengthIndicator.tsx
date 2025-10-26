import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

interface PasswordStrengthIndicatorProps {
  password: string;
}

interface Particle {
  id: number;
  x: number;
  y: number;
  vx: number;
  vy: number;
}

export function PasswordStrengthIndicator({
  password,
}: PasswordStrengthIndicatorProps) {
  const { t } = useTranslation();
  const [strength, setStrength] = useState(0);
  const [particles, setParticles] = useState<Particle[]>([]);
  const [nextId, setNextId] = useState(0);

  const calculatePasswordStrength = (pwd: string): number => {
    let score = 0;
    const checks = {
      length: pwd.length >= 8,
      uppercase: /[A-Z]/.test(pwd),
      lowercase: /[a-z]/.test(pwd),
      numbers: /[0-9]/.test(pwd),
      special: /[!@#$%^&*()_+=\-{};':"\\|,.<>/?]/.test(pwd),
    };

    Object.values(checks).forEach((check) => {
      if (check) score += 20;
    });

    return score;
  };

  useEffect(() => {
    const newStrength = calculatePasswordStrength(password);

    // Add particles when strength increases
    if (newStrength > strength) {
      const newParticles: Particle[] = [];
      for (let i = 0; i < 2; i++) {
        newParticles.push({
          id: nextId + i,
          x: 50 + (Math.random() - 0.5) * 15,
          y: 50,
          vx: (Math.random() - 0.5) * 1.5,
          vy: Math.random() * 1.5 + 0.5,
        });
      }
      setParticles((prev) => [...prev, ...newParticles]);
      setNextId((prev) => prev + 2);
    }

    setStrength(newStrength);
  }, [password, strength, nextId]);

  useEffect(() => {
    const interval = setInterval(() => {
      setParticles((prev) => {
        const updated = prev
          .map((p) => ({
            ...p,
            x: p.x + p.vx,
            y: p.y + p.vy,
            vy: p.vy + 0.15, // gravity
          }))
          .filter((p) => p.y < 75); // Remove particles that fall out

        return updated;
      });
    }, 30);

    return () => clearInterval(interval);
  }, []);

  const getStrengthColor = () => {
    if (strength < 20) return "#ef4444";
    if (strength < 40) return "#f97316";
    if (strength < 60) return "#eab308";
    if (strength < 80) return "#84cc16";
    return "#22c55e";
  };

  return (
    <div className="relative">
      {/* Hourglass SVG */}
      <svg
        className="absolute right-0 top-0 w-24 h-20 pointer-events-none"
        viewBox="0 0 100 100"
        xmlns="http://www.w3.org/2000/svg"
      >
        {/* Hourglass outline with rounded edges */}
        <path
          d="M 20 10 L 20 40 C 20 45, 30 50, 50 50 C 70 50, 80 45, 80 40 L 80 10 M 20 90 L 20 60 C 20 55, 30 50, 50 50 C 70 50, 80 55, 80 60 L 80 90 M 30 15 L 70 15 M 30 85 L 70 85"
          stroke="currentColor"
          strokeWidth="1.5"
          fill="none"
          className="text-border"
          strokeLinecap="round"
          strokeLinejoin="round"
        />

        {/* Top chamber - empty */}
        <path
          d="M 25 15 C 25 15, 50 40, 50 40 C 50 40, 75 15, 75 15 Z"
          fill="none"
          opacity="0.2"
        />

        {/* Bottom chamber - fill based on strength */}
        <defs>
          <clipPath id="bottomChamber">
            <path d="M 25 85 C 25 85, 50 60, 50 60 C 50 60, 75 85, 75 85 Z" />
          </clipPath>
        </defs>

        <rect
          x="25"
          y={60 + (25 - (strength / 100) * 25)}
          width="50"
          height={(strength / 100) * 25}
          fill={getStrengthColor()}
          opacity="0.8"
          clipPath="url(#bottomChamber)"
        />
      </svg>

      {/* Particles */}
      <div className="absolute right-0 top-0 w-24 h-20 pointer-events-none overflow-hidden">
        {particles.map((particle) => (
          <div
            key={particle.id}
            className="absolute w-0.5 h-0.5 rounded-full"
            style={{
              left: `${particle.x}%`,
              top: `${particle.y}%`,
              backgroundColor: getStrengthColor(),
              opacity: 0.8,
              transition: "none",
            }}
          />
        ))}
      </div>

      {/* Requirements list */}
      <div className="pr-28 space-y-1 text-[10px]">
        <div className="flex items-center gap-1.5">
          <div
            className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${
              password.length >= 8 ? "bg-green-500" : "bg-border"
            }`}
          />
          <span
            className={
              password.length >= 8 ? "text-green-500" : "text-muted-foreground"
            }
          >
            {t("passwordStrength.eightChars")}
          </span>
        </div>
        <div className="flex items-center gap-1.5">
          <div
            className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${
              /[A-Z]/.test(password) ? "bg-green-500" : "bg-border"
            }`}
          />
          <span
            className={
              /[A-Z]/.test(password)
                ? "text-green-500"
                : "text-muted-foreground"
            }
          >
            {t("passwordStrength.uppercase")}
          </span>
        </div>
        <div className="flex items-center gap-1.5">
          <div
            className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${
              /[a-z]/.test(password) ? "bg-green-500" : "bg-border"
            }`}
          />
          <span
            className={
              /[a-z]/.test(password)
                ? "text-green-500"
                : "text-muted-foreground"
            }
          >
            {t("passwordStrength.lowercase")}
          </span>
        </div>
        <div className="flex items-center gap-1.5">
          <div
            className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${
              /[0-9]/.test(password) ? "bg-green-500" : "bg-border"
            }`}
          />
          <span
            className={
              /[0-9]/.test(password)
                ? "text-green-500"
                : "text-muted-foreground"
            }
          >
            {t("passwordStrength.number")}
          </span>
        </div>
        <div className="flex items-center gap-1.5">
          <div
            className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${
              /[!@#$%^&*()_+=\-{};':"\\|,.<>/?]/.test(password)
                ? "bg-green-500"
                : "bg-border"
            }`}
          />
          <span
            className={
              /[!@#$%^&*()_+=\-{};':"\\|,.<>/?]/.test(password)
                ? "text-green-500"
                : "text-muted-foreground"
            }
          >
            {t("passwordStrength.special")}
          </span>
        </div>
      </div>
    </div>
  );
}
