// Sample types based on Go API models
export interface Account {
  id: string;
  timezone_id?: number;
  created_at: string;
  updated_at: string;
  timezone?: Timezone;
}

export interface Timezone {
  id: number;
  name: string;
  value: string;
}

export interface Reminder {
  id: string;
  account_id: string;
  remind_at_utc: string;
  snoozed_at_utc?: string;
  next_fire_utc?: string;
  message: string;
  created_at: string;
  recurrence: number;
  destinations?: ReminderDestination[];
}

export interface ReminderDestination {
  id: string;
  reminder_id: string;
  type: string;
  destination: string;
  created_at: string;
}

export interface Identity {
  id: string;
  account_id: string;
  provider: string;
  provider_id: string;
  created_at: string;
}
