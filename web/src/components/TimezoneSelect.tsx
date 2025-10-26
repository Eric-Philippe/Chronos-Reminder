import { useState } from "react";
import { Input } from "@/components/ui/input";

const TIMEZONES = [
  "UTC",
  "UTC+1",
  "UTC+2",
  "UTC+3",
  "UTC+3:30",
  "UTC+4",
  "UTC+4:30",
  "UTC+5",
  "UTC+5:30",
  "UTC+5:45",
  "UTC+6",
  "UTC+6:30",
  "UTC+7",
  "UTC+8",
  "UTC+8:45",
  "UTC+9",
  "UTC+9:30",
  "UTC+10",
  "UTC+10:30",
  "UTC+11",
  "UTC+12",
  "UTC+12:45",
  "UTC+13",
  "UTC+14",
  "UTC-1",
  "UTC-2",
  "UTC-3",
  "UTC-3:30",
  "UTC-4",
  "UTC-4:30",
  "UTC-5",
  "UTC-6",
  "UTC-7",
  "UTC-8",
  "UTC-8:30",
  "UTC-9",
  "UTC-9:30",
  "UTC-10",
  "UTC-11",
  "UTC-12",
  "Africa/Johannesburg",
  "Africa/Cairo",
  "Africa/Lagos",
  "Africa/Nairobi",
  "America/New_York",
  "America/Chicago",
  "America/Denver",
  "America/Los_Angeles",
  "America/Toronto",
  "America/Mexico_City",
  "America/Buenos_Aires",
  "America/Sao_Paulo",
  "Asia/Dubai",
  "Asia/Bangkok",
  "Asia/Singapore",
  "Asia/Hong_Kong",
  "Asia/Tokyo",
  "Asia/Seoul",
  "Asia/Shanghai",
  "Asia/India",
  "Asia/Jakarta",
  "Australia/Sydney",
  "Australia/Melbourne",
  "Australia/Brisbane",
  "Australia/Perth",
  "Europe/London",
  "Europe/Paris",
  "Europe/Berlin",
  "Europe/Moscow",
  "Europe/Istanbul",
  "Europe/Amsterdam",
  "Pacific/Auckland",
  "Pacific/Fiji",
  "Pacific/Honolulu",
].sort();

interface TimezoneSelectProps {
  value: string;
  onChange: (value: string) => void;
}

export function TimezoneSelect({ value, onChange }: TimezoneSelectProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

  const filteredTimezones = TIMEZONES.filter((tz) =>
    tz.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="relative">
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="w-full bg-secondary/50 border border-border text-foreground placeholder:text-muted-foreground rounded-md px-3 py-2 text-sm flex justify-between items-center hover:bg-secondary/70 transition-colors"
      >
        <span>{value || "Select timezone"}</span>
        <svg
          className={`w-4 h-4 transition-transform ${
            isOpen ? "rotate-180" : ""
          }`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 14l-7 7m0 0l-7-7m7 7V3"
          />
        </svg>
      </button>

      {isOpen && (
        <div className="absolute top-full left-0 right-0 mt-2 bg-card border border-border rounded-md shadow-lg z-50">
          {/* Search Box */}
          <div className="p-2 border-b border-border">
            <Input
              type="text"
              placeholder="Search timezone..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="bg-secondary/50 border-border text-foreground placeholder:text-muted-foreground text-sm"
              autoFocus
            />
          </div>

          {/* Timezone List */}
          <div className="max-h-64 overflow-y-auto">
            {filteredTimezones.length > 0 ? (
              filteredTimezones.map((tz) => (
                <button
                  key={tz}
                  type="button"
                  onClick={() => {
                    onChange(tz);
                    setIsOpen(false);
                    setSearchQuery("");
                  }}
                  className={`w-full text-left px-3 py-2 text-sm transition-colors ${
                    value === tz
                      ? "bg-accent text-accent-foreground"
                      : "text-foreground hover:bg-secondary/50"
                  }`}
                >
                  {tz}
                </button>
              ))
            ) : (
              <div className="px-3 py-2 text-sm text-muted-foreground text-center">
                No timezones found
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
