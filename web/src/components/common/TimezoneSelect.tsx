import { useState } from "react";
import { Input } from "@/components/ui/input";

// IANA timezone identifiers
const TIMEZONES = [
  "Africa/Cairo",
  "Africa/Johannesburg",
  "Africa/Lagos",
  "Africa/Nairobi",
  "America/Anchorage",
  "America/Argentina/Buenos_Aires",
  "America/Chicago",
  "America/Denver",
  "America/Los_Angeles",
  "America/Mexico_City",
  "America/New_York",
  "America/Phoenix",
  "America/Toronto",
  "America/Vancouver",
  "Asia/Bangkok",
  "Asia/Dubai",
  "Asia/Hong_Kong",
  "Asia/Kolkata",
  "Asia/Jakarta",
  "Asia/Manila",
  "Asia/Shanghai",
  "Asia/Singapore",
  "Asia/Seoul",
  "Asia/Tokyo",
  "Atlantic/Azores",
  "Atlantic/Cape_Verde",
  "Australia/Adelaide",
  "Australia/Brisbane",
  "Australia/Melbourne",
  "Australia/Perth",
  "Australia/Sydney",
  "Europe/Amsterdam",
  "Europe/Athens",
  "Europe/Belgrade",
  "Europe/Berlin",
  "Europe/Brussels",
  "Europe/Bucharest",
  "Europe/Budapest",
  "Europe/Dublin",
  "Europe/Helsinki",
  "Europe/Istanbul",
  "Europe/Lisbon",
  "Europe/London",
  "Europe/Madrid",
  "Europe/Moscow",
  "Europe/Paris",
  "Europe/Prague",
  "Europe/Rome",
  "Europe/Stockholm",
  "Europe/Vienna",
  "Europe/Warsaw",
  "Europe/Zurich",
  "Pacific/Auckland",
  "Pacific/Fiji",
  "Pacific/Honolulu",
  "Pacific/Tongatapu",
  "UTC",
].sort();

interface TimezoneSelectProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  searchPlaceholder?: string;
  noResultsText?: string;
}

export function TimezoneSelect({
  value,
  onChange,
  placeholder = "Select timezone",
  searchPlaceholder = "Search timezone...",
  noResultsText = "No timezones found",
}: TimezoneSelectProps) {
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
        <span>{value || placeholder}</span>
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
              placeholder={searchPlaceholder}
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
                {noResultsText}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
