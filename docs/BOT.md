# Chronos Bot Commands Documentation

Welcome to the Chronos Reminder Bot! This document provides a complete guide to all available commands. The bot is designed to help you manage reminders across Discord with timezone support and flexible scheduling options.

---

## Table of Contents

1. [Reminder Commands](#reminder-commands)
2. [User Management Commands](#user-management-commands)
3. [Utility Commands](#utility-commands)
4. [General Commands](#general-commands)

---

## Reminder Commands

### 📝 `/reminders` - Manage Your Reminders

**Category:** Reminders

**Short Description:** Manage your reminders

**Full Description:** List, show, pause, restart, or delete your existing reminders. This is the main command for managing all your created reminders with various subcommands.

**Usage:**

```
/reminders <subcommand> [options]
```

**Subcommands:**

#### `list`

Lists all your reminders with their details.

- **Usage:** `/reminders list`
- **What it does:** Displays all active reminders you've created

#### `show`

Shows detailed information about a specific reminder.

- **Usage:** `/reminders show reminder:<reminder>`
- **Parameters:**
  - `reminder` (Required): Select the reminder from autocomplete suggestions
- **Example:** `/reminders show reminder:[14:30] Take medication`

#### `pause`

Pauses a recurring reminder (does not pause one-time reminders).

- **Usage:** `/reminders pause reminder:<reminder>`
- **Parameters:**
  - `reminder` (Required): Select a recurring reminder from autocomplete
- **Note:** Only works with recurring reminders
- **Example:** `/reminders pause reminder:[Daily 09:00] Daily standup`

#### `unpause`

Restarts a paused reminder.

- **Usage:** `/reminders unpause reminder:<reminder>`
- **Parameters:**
  - `reminder` (Required): Select a paused reminder from autocomplete
- **Note:** Only works with currently paused reminders
- **Example:** `/reminders unpause reminder:[Daily 09:00] Daily standup`

#### `delete`

Permanently deletes a reminder.

- **Usage:** `/reminders delete reminder:<reminder>`
- **Parameters:**
  - `reminder` (Required): Select the reminder to delete from autocomplete
- **⚠️ Warning:** This action is permanent and cannot be undone
- **Example:** `/reminders delete reminder:[Tomorrow 15:00] Team meeting`

---

### ⏰ `/remindme` - Create a Personal Reminder

**Category:** Reminders

**Short Description:** Create a new reminder

**Full Description:** Create a new reminder that will be sent to you via direct message at the specified date and time. Your personal reminders are sent as private messages.

**Usage:**

```
/remindme message:<text> date:<date> time:<time> [recurrence:<type>]
```

**Parameters:**

- `message` (Required, String): The reminder message content
- `date` (Required, String): The date for the reminder
  - Supported formats: `today`, `tomorrow`, `next week`, `next month`, `25/12/2024`, `2024-12-25`, etc.
  - Autocomplete suggestions available
- `time` (Required, String): The time for the reminder
  - Supported formats: `15:30`, `3pm`, `9:30am`, `15.5` (hours)
- `recurrence` (Optional, String, Default: `ONCE`): How often to repeat the reminder
  - **Options:**
    - `ONCE` - One-time reminder (default)
    - `HOURLY` - Every hour
    - `DAILY` - Every day
    - `WEEKLY` - Every week (same day)
    - `MONTHLY` - Every month (same date)
    - `YEARLY` - Every year (same date and month)
    - `WORKDAYS` - Monday through Friday
    - `WEEKEND` - Saturday and Sunday

**Examples:**

```
/remindme message:"Take medicine" date:"today" time:"15:30"
/remindme message:"Team meeting" date:"25/12/2024" time:"10:00" recurrence:daily
/remindme message:"Birthday" date:"tomorrow" time:"9am"
/remindme message:"Weekly review" date:"next week" time:"2pm" recurrence:weekly
```

**Notes:**

- Reminders are sent via direct message
- Times are in your configured timezone
- You can create as many reminders as needed
- Times must be in the future

---

### 📢 `/remindus` - Create a Channel Reminder

**Category:** Reminders

**Short Description:** Create a new reminder in a channel

**Full Description:** Create a new reminder that will be sent in a specified channel at the specified date and time. Requires 'Manage Channel', 'Administrator' permission, or server ownership. You can optionally mention a role to ping when the reminder is sent.

**Usage:**

```
/remindus message:<text> date:<date> time:<time> channel:<channel> [role:<role>] [recurrence:<type>]
```

**Parameters:**

- `message` (Required, String): The reminder message content
- `date` (Required, String): The date for the reminder
  - Same formats as `/remindme`
- `time` (Required, String): The time for the reminder
  - Same formats as `/remindme`
- `channel` (Required, Channel): The Discord channel where the reminder will be posted
- `role` (Optional, Role): Role to mention in the reminder message
  - Requires 'Manage Roles' permission to mention roles
  - The bot must have a higher role than the target role
- `recurrence` (Optional, String, Default: `ONCE`): How often to repeat
  - Same options as `/remindme`

**Permissions Required:**

- You need one of: `Manage Channel`, `Administrator`, or server ownership
- If mentioning a role, you also need `Manage Roles` permission
- The bot needs permission to mention the role (higher role hierarchy)

**Examples:**

```
/remindus message:"Team meeting" date:"25/12/2024" time:"10:00" channel:#general
/remindus message:"Daily standup" date:"tomorrow" time:"9am" channel:#dev-team recurrence:workdays
/remindus message:"Important announcement" date:"next week" time:"2pm" channel:#announcements role:@everyone recurrence:weekly
```

**Notes:**

- Reminders are posted as channel messages
- Only users with proper permissions can create channel reminders
- Role mentions require additional permissions
- Times are in the creator's configured timezone

---

## User Management Commands

### 🌏 `/timezone` - Manage Timezones

**Category:** User

**Short Description:** Manage timezones

**Full Description:** List available timezones or change your current timezone. All reminders use your configured timezone for scheduling.

**Usage:**

```
/timezone <subcommand>
```

**Subcommands:**

#### `list`

Displays all available timezones.

- **Usage:** `/timezone list`
- **What it does:** Shows a paginated list of all supported IANA timezones

#### `change`

Opens an interactive menu to change your timezone.

- **Usage:** `/timezone change`
- **What it does:** Presents timezone options organized by region for easy selection
- **Note:** Changing your timezone affects all future reminders

#### `display`

Shows your currently configured timezone.

- **Usage:** `/timezone display`
- **What it does:** Displays your current timezone setting

**Examples:**

```
/timezone list
/timezone change
/timezone display
```

**Notes:**

- Timezones follow the IANA timezone database (e.g., `America/New_York`, `Europe/London`, `Asia/Tokyo`)
- Your timezone is used to interpret all date and time inputs
- Reminders are stored in UTC but displayed in your timezone

---

### 👤 `/profile` - View User Profile

**Category:** User

**Short Description:** View user profile

**Full Description:** Display a user's profile with their avatar, account creation date, reminder count, and platform badges. Use without parameters to view your own profile, or specify a user to view theirs.

**Usage:**

```
/profile [user:@user]
```

**Parameters:**

- `user` (Optional, User): The user whose profile to view
  - Leave empty to view your own profile
  - Specify a user to view their profile

**Profile Information Displayed:**

- User avatar
- Account creation date with Chronos
- Total reminder count
- Platform badges (if applicable)
- User mention link

**Examples:**

```
/profile
/profile user:@john_doe
/profile user:@Alice#1234
```

**Notes:**

- You can view any user's profile
- Only basic public information is displayed
- Profile data is synced when you first use the bot

---

## Utility Commands

### 🧮 `/calcultime` - Calculate Time Operations

**Category:** Tools

**Short Description:** Calculate time operations

**Full Description:** Perform calculations between times or with factors. Supports addition, subtraction, multiplication, and division of time values. Useful for quick time math.

**Usage:**

```
/calcultime time1:<time> operation:<operation> time2:<time/factor>
```

**Parameters:**

- `time1` (Required, String): First time value
  - Supported formats: `2h 30m`, `14:30`, `2.5h`, `150m`, `9000s`, etc.
- `operation` (Required, String): Operation to perform
  - **Options:**
    - `ADD` (+) - Add two time values
    - `SUBTRACT` (-) - Subtract one time from another
    - `MULTIPLY` (×) - Multiply a time by a factor
    - `DIVIDE` (÷) - Divide a time by a factor
- `time2` (Required, String): Second time value or factor
  - For ADD/SUBTRACT: time value (e.g., `1h 15m`)
  - For MULTIPLY/DIVIDE: numeric factor (e.g., `2.5`, `0.75`)

**Examples:**

```
/calcultime time1:"2h 30m" operation:add time2:"1h 15m"
(Result: 3h 45m)

/calcultime time1:"5h" operation:subtract time2:"1h 30m"
(Result: 3h 30m)

/calcultime time1:"1h 30m" operation:multiply time2:"2"
(Result: 3h)

/calcultime time1:"2h" operation:divide time2:"4"
(Result: 30m)
```

**Supported Time Formats:**

- Hours and minutes: `2h 30m`, `1h 15m`
- Clock time: `14:30`, `09:45`
- Decimal hours: `2.5h`, `1.75h`
- Minutes only: `150m`, `90m`
- Seconds: `9000s`, `5400s`
- Combined: `1h 30m 45s`

**Notes:**

- All calculations are performed precisely
- Results are formatted in human-readable format
- Useful for planning and time management

---

## General Commands

### ⏰ `/tic` - Ping the Bot

**Category:** General

**Short Description:** Ping the bot

**Full Description:** Ping the bot and get a response to verify that it's online and responsive. A quick way to check if the bot is working.

**Usage:**

```
/tic
```

**Response:** The bot responds with "The bot is alive! ⏰ Tac!"

**Example:**

```
/tic
→ The bot is alive! ⏰ Tac !
```

**Notes:**

- No parameters needed
- Useful for troubleshooting connectivity
- Doesn't require an account

---

### ⏳ `/hourglass` - Start a Short Timer

**Category:** General

**Short Description:** Start a short timer (max 30 minutes)!

**Full Description:** Start a quick in-memory timer that will notify you when it ends. For longer or persistent reminders, use `/remindme` or `/remindus`. Perfect for quick focus sessions, cooking timers, or short-term countdowns.

**Usage:**

```
/hourglass duration:<duration> message:<message>
```

**Parameters:**

- `duration` (Required, String): Duration in seconds or minutes (maximum 30 minutes)
  - Supported formats: `10s`, `5m`, `2m30s`, etc.
  - Maximum: `30m` (30 minutes)
- `message` (Required, String): The message to send when the timer ends
  - Can include any text

**Examples:**

```
/hourglass duration:25m message:Focus time!
/hourglass duration:10m message:Take a break
/hourglass duration:2m30s message:Check oven
/hourglass duration:1m message:Presentation starts
```

**Limitations:**

- ⏳ Maximum duration: 30 minutes
- ⏳ Only supports seconds and minutes (not hours)
- ⏳ Timers are in-memory (not persistent)
- ⏳ For longer reminders, use `/remindme` or `/remindus`

**Notes:**

- Timers are sent as immediate acknowledgement with end time
- Notification is sent to the same channel when timer expires
- You are mentioned in the notification
- Perfect for Pomodoro technique (25 minute focus sessions)

---

### 📕 `/help` - Get Help with Commands

**Category:** General

**Short Description:** Get help with bot commands

**Full Description:** Provides information about available commands and how to use them. Use this command to learn more about the bot's features and functionalities. Can provide general help or specific information about a particular command.

**Usage:**

```
/help [command:<name>]
```

**Parameters:**

- `command` (Optional, String): The specific command to get help with
  - Leave empty for general help with all commands
  - Examples: `remindme`, `remindus`, `profile`, `timezone`, etc.

**Examples:**

```
/help
(Shows all available commands organized by category)

/help command:remindme
(Shows detailed information about /remindme)

/help command:timezone
(Shows detailed information about /timezone)
```

**Features:**

- Organized by command categories
- Shows usage examples
- Displays available options and parameters
- Explains permissions and requirements

**Notes:**

- Help is always available and accessible
- No account needed
- Great for learning new commands

---

### 💡 `/support` - Get Help and Support Resources

**Category:** General

**Short Description:** Get help and support resources

**Full Description:** Access documentation, FAQs, contact information, and our official Discord server for support, feedback, and community connection.

**Usage:**

```
/support
```

**Resources Provided:**

- **📖 Documentation & Contact**: https://chronosrmd.com/
  - Full documentation and FAQ
  - Contact information
- **💬 Official Discord Server**: https://discord.gg/m3MsM922QD
  - Community support
  - Share feedback
  - Connect with other users

**Example:**

```
/support
→ Shows support resources and links
```

**Notes:**

- No parameters needed
- Always accessible
- Great for reporting issues or suggesting features

---

## Quick Reference

### Command Categories

**Reminders:**

- `/remindme` - Personal DM reminders
- `/remindus` - Channel reminders
- `/reminders` - Manage existing reminders

**User Management:**

- `/timezone` - Manage your timezone
- `/profile` - View user profiles

**Utilities:**

- `/calcultime` - Time calculations

**General:**

- `/tic` - Ping the bot
- `/hourglass` - Quick timer
- `/help` - Get command help
- `/support` - Support resources

### Recurrence Options

All reminder commands support these recurrence types:

- `ONCE` - One-time reminder (default)
- `HOURLY` - Every hour
- `DAILY` - Every day
- `WEEKLY` - Every week (same day)
- `MONTHLY` - Every month (same date)
- `YEARLY` - Every year
- `WORKDAYS` - Monday through Friday
- `WEEKEND` - Saturday and Sunday

### Supported Time Formats

- **Date:** `today`, `tomorrow`, `next week`, `25/12/2024`, `2024-12-25`
- **Time:** `15:30`, `3pm`, `9:30am`, `15.5`
- **Duration:** `10s`, `5m`, `2h`, `2h 30m`

---

## Tips & Tricks

1. **Use Autocomplete**: When creating reminders, use the date autocomplete to quickly select common dates
2. **Timezone Management**: Set your timezone correctly to ensure reminders trigger at the right time
3. **Channel Reminders**: Use `/remindus` for team reminders with optional role mentions
4. **Pausing Reminders**: Pause recurring reminders instead of deleting them if you might want them back
5. **Quick Timers**: Use `/hourglass` for Pomodoro sessions (25 minutes is recommended)
6. **Permission Levels**: Admin users can manage server reminders, regular users can only manage their own

---

## Troubleshooting

**Issue: "The specified date and time is in the past"**

- Solution: Make sure you're setting a future date and time

**Issue: "Insufficient Permissions" on `/remindus`**

- Solution: Ensure you have Manage Channel or Administrator permissions, or are the server owner

**Issue: "Invalid Timezone"**

- Solution: Use `/timezone list` to see available timezones and change with `/timezone change`

**Issue: "The bot needs 'Mention Everyone' permission"**

- Solution: Grant the bot proper Discord permissions in channel/role settings if you want role mentions in reminders

---

## Getting Help

- Use `/help` to learn about any command
- Use `/support` to access documentation and community
- Visit https://chronosrmd.com/ for full documentation
- Join the official Discord: https://discord.gg/m3MsM922QD

---

_Last Updated: 2025_
_Chronos Reminder Bot - Making time management easy_
