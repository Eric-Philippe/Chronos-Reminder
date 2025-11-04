import { useState } from "react";
import { ChevronDown } from "lucide-react";

interface IntegrationCardProps {
  name: string;
  icon: React.ReactNode;
  description: string;
  features: string[];
  apiKeySupport?: boolean;
  isApiKeyCard?: boolean;
}

export function IntegrationCard({
  name,
  icon,
  description,
  features,
}: IntegrationCardProps) {
  const [isExpanded, setIsExpanded] = useState(false);

  return (
    <div
      className="relative h-full"
      onMouseEnter={() => setIsExpanded(true)}
      onMouseLeave={() => setIsExpanded(false)}
    >
      {/* Glowing border animation - inline style */}
      {isExpanded && (
        <style>{`
          @keyframes glowingBorder {
            0% {
              box-shadow: 
                0 0 20px rgba(var(--accent-rgb), 0),
                inset 0 0 20px rgba(var(--accent-rgb), 0);
            }
            25% {
              box-shadow: 
                20px 0 20px rgba(var(--accent-rgb), 0.5),
                inset -20px 0 20px rgba(var(--accent-rgb), 0.1);
            }
            50% {
              box-shadow: 
                0 20px 20px rgba(var(--accent-rgb), 0.5),
                inset 0 -20px 20px rgba(var(--accent-rgb), 0.1);
            }
            75% {
              box-shadow: 
                -20px 0 20px rgba(var(--accent-rgb), 0.5),
                inset 20px 0 20px rgba(var(--accent-rgb), 0.1);
            }
            100% {
              box-shadow: 
                0 0 20px rgba(var(--accent-rgb), 0),
                inset 0 0 20px rgba(var(--accent-rgb), 0);
            }
          }
          .glow-animation {
            animation: glowingBorder 3s ease-in-out infinite !important;
          }
        `}</style>
      )}

      {/* Collapsed State */}
      <div
        className={`
          relative p-6 rounded-xl bg-background/50 border border-white/5 
          transition-all duration-300 cursor-pointer h-full
          ${
            isExpanded
              ? "glow-animation lg:ring-2 lg:ring-accent/50"
              : "hover:border-accent/30 hover:shadow-lg hover:shadow-accent/10"
          }
        `}
      >
        <div className="flex flex-col items-center gap-4">
          <div className="w-16 h-16 rounded-lg bg-accent/20 flex items-center justify-center transition-colors">
            {icon}
          </div>
          <h3 className="text-lg font-semibold text-foreground text-center">
            {name}
          </h3>
          <p className="text-sm text-muted-foreground text-center leading-relaxed">
            {description}
          </p>

          {/* Expand indicator */}
          <div
            className={`
              transition-transform duration-300 text-muted-foreground
              ${isExpanded ? "rotate-180" : ""}
            `}
          >
            <ChevronDown className="w-5 h-5" />
          </div>
        </div>

        {/* Expanded State - Overlay */}
        {isExpanded && (
          <div
            className={`
              absolute inset-0 p-6 rounded-xl bg-background/95 backdrop-blur-sm
              border border-accent/50 flex flex-col gap-4
              transition-opacity duration-300
              opacity-100
            `}
          >
            <div className="space-y-3">
              <h4 className="text-sm font-semibold text-accent uppercase tracking-wide">
                Features
              </h4>
              <ul className="space-y-2">
                {features.map((feature, index) => (
                  <li key={index} className="flex items-start gap-2">
                    <span className="text-accent mt-1">✓</span>
                    <span className="text-sm text-muted-foreground">
                      {feature}
                    </span>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
