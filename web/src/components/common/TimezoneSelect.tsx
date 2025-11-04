import { useEffect, useState } from "react";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { timezoneService } from "@/services";
import type { Timezone } from "@/services/types";

interface TimezoneSelectProps {
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
}

// Helper function to extract region from IANA location
function getRegion(iana: string): string {
  const parts = iana.split("/");
  if (parts.length === 2) {
    return parts[0];
  }
  return "Other";
}

export function TimezoneSelect({
  value,
  onChange,
  disabled = false,
}: TimezoneSelectProps) {
  const [timezones, setTimezones] = useState<Timezone[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Fetch timezones from API
  useEffect(() => {
    const fetchTimezones = async () => {
      try {
        setIsLoading(true);
        const data = await timezoneService.getAvailableTimezones();
        setTimezones(data);
      } catch (error) {
        console.error("Failed to fetch timezones from API", error);
        setTimezones([]);
      } finally {
        setIsLoading(false);
      }
    };

    fetchTimezones();
  }, []);

  // Group timezones by region (first part of IANA location)
  const groupedTimezones = timezones.reduce((groups, tz) => {
    const region = getRegion(tz.iana_location);
    if (!groups[region]) {
      groups[region] = [];
    }
    groups[region].push(tz);
    return groups;
  }, {} as Record<string, Timezone[]>);

  // Sort regions and timezones within each region
  const sortedRegions = Object.keys(groupedTimezones).sort();

  return (
    <Select
      value={value}
      onValueChange={onChange}
      disabled={isLoading || disabled}
    >
      <SelectTrigger className="w-full">
        <SelectValue placeholder="Select a timezone" />
      </SelectTrigger>
      <SelectContent className="max-h-[300px]">
        {sortedRegions.map((region) => (
          <SelectGroup key={region}>
            <SelectLabel>{region}</SelectLabel>
            {groupedTimezones[region]
              .sort((a, b) => a.iana_location.localeCompare(b.iana_location))
              .map((tz) => (
                <SelectItem key={tz.iana_location} value={tz.iana_location}>
                  {tz.iana_location}
                </SelectItem>
              ))}
          </SelectGroup>
        ))}
      </SelectContent>
    </Select>
  );
}
