import type { Account, Reminder, Timezone } from "@/types/models";

export const sampleTimezone: Timezone = {
  id: 1,
  name: "Eastern Time",
  value: "America/New_York",
};

export const sampleAccount: Account = {
  id: "550e8400-e29b-41d4-a716-446655440000",
  timezone_id: 1,
  created_at: "2025-01-15T10:30:00Z",
  updated_at: "2025-10-24T14:22:00Z",
  timezone: sampleTimezone,
};

export const sampleReminders: Reminder[] = [
  {
    id: "660e8400-e29b-41d4-a716-446655440001",
    account_id: sampleAccount.id,
    remind_at_utc: "2025-10-25T14:00:00Z",
    next_fire_utc: "2025-10-25T14:00:00Z",
    message: "Team standup meeting",
    created_at: "2025-10-22T09:15:00Z",
    recurrence: 1, // Daily
    destinations: [
      {
        id: "770e8400-e29b-41d4-a716-446655440001",
        reminder_id: "660e8400-e29b-41d4-a716-446655440001",
        type: "discord_dm",
        destination: "discord_user_123",
        created_at: "2025-10-22T09:15:00Z",
      },
    ],
  },
  {
    id: "660e8400-e29b-41d4-a716-446655440002",
    account_id: sampleAccount.id,
    remind_at_utc: "2025-10-25T18:30:00Z",
    next_fire_utc: "2025-10-25T18:30:00Z",
    message: "Project deadline - Submit final deliverables",
    created_at: "2025-10-20T11:45:00Z",
    recurrence: 0, // One-time
    destinations: [
      {
        id: "770e8400-e29b-41d4-a716-446655440002",
        reminder_id: "660e8400-e29b-41d4-a716-446655440002",
        type: "discord_channel",
        destination: "project-updates",
        created_at: "2025-10-20T11:45:00Z",
      },
    ],
  },
  {
    id: "660e8400-e29b-41d4-a716-446655440003",
    account_id: sampleAccount.id,
    remind_at_utc: "2025-10-26T09:00:00Z",
    next_fire_utc: "2025-10-26T09:00:00Z",
    message: "Review pull requests",
    created_at: "2025-10-23T16:20:00Z",
    recurrence: 2, // Weekly
    destinations: [
      {
        id: "770e8400-e29b-41d4-a716-446655440003",
        reminder_id: "660e8400-e29b-41d4-a716-446655440003",
        type: "discord_dm",
        destination: "discord_user_123",
        created_at: "2025-10-23T16:20:00Z",
      },
    ],
  },
];
