# Chronos v1.0 Roadmap

## Database

- [ ] Add a `last_notified_at` column to the `reminders` table to store the last time a reminder was sent
- [ ] Add a logging table to track reminder deliveries and failures
- [ ] Use REDIS ?
- - [ ] Cache paused reminders
- - [x] Cache user account for Discord.EnsureAccount method
- - [x] Cache user timezone preferences

## Discord Bot

- [x] `/remindme` command -> Quickly set {DISCORD_DM}Reminders type
- [x] `/remindus` command -> Set {DISCORD_CHANNEL}Reminders type, (MUST BE ENTERED IN A SERVER, NOT IN A DM), the channel field is an autocomplete that gives the user channels select where the user has the manage channel permission
- [x] Draw the discord reminder with a prettier display (image...)
- [x] `/profile` command -> View user profile information and their integrations and buttons to manage them
- [x] `/reminders delete|list|show|pause|unpause` commands -> Manage reminders
- - [x] `list` subcommand -> Give a first embed with all the list, and a second embed with a paginated list of one reminders that the user can scroll through, show the remaining time before the next reminder
- - [x] `delete` subcommand -> Delete a reminder by its content/ID, or if not given, launch a select menu to first select between all its different type of reminders (DM, Channel, Webhook...) or skip this step if only one type exists, then display a second select menu with all the reminders of this type to select which one to delete, and then remove the message and put a confirmation message. Only when targeting a {DISCORD_CHANNEL}Reminder, If the user is admin of the current guild, even if a reminder was created by someone else in this guild, he can delete it.
- - [x] `show` subcommand -> Get a reminder by its content/ID and display it in a pretty embed or image, and little icons to show its types/destination(s) (footer with: create website account for update/delete reminders), show the remaining time before the next reminder
- - [x] `pause` subcommand -> Pause a reminder by its content/ID, return an error message if the reminder is not a recurring one.
- - [x] `restart` subcommand -> Unpause a reminder by its content/ID, return an error message if the reminder is not a recurring one.
- [x] `/calcultime` command -> Time calculator
- [ ] `/help` command -> Help command with all the commands and their descriptions, and buttons to get more information on each command
- [ ] Test behavior when the bot can't send the reminder DM/CHANNEL (user blocked the bot, user left the server, bot kicked from the server, no permission to send messages in the channel...)
- [ ] Being able to snooze a reminder when received (only for DM reminders)
- - [x] Add the necessary fields in the database
- - [x] Add buttons to the reminder message
- - [x] Add the logic to handle a incoming reminder from snoozing
- - [x] Add the delete reminder queue for not snoozed one time reminders
- - [x] Add snooze operation handler and plug it to the buttons
- - [ ] Prevent snoozing a reminder after its next recurrence time
- [x] Add testings for the timeparser

## API

- [ ] Create API endpoints for the web application to interact with the reminder engine
- [ ] Implement authentication for API endpoints
- [ ] API Key for third-party integrations

## Web Application

- [ ] https://chronosreminder.app/.fr?
- [ ] Welcome unsigned page
- [x] Create account/login => (Possibility to merge from an already existing Discord account => OAuth2 Discord => Create account prefiled with Discord info)
- [ ] CRUD on reminders
- [ ] View profile and integrations
- [ ] Link Discord account
- [ ] View history/logs
- [x] Same slick style as SnapFileThing

## Reminder engine

- [x] Implement main scheduler queue
- [x] Implement DM_DISCORD dispatcher
- [x] Implement CHANNEL_DISCORD dispatcher
- - [x] Add role mention support
- [ ] Implement WEBHOOK dispatcher
- [x] Implement recurrence handling
- [x] Recalculate next occurrence on restart to avoid reminders trying to catch up spam
- [x] Add a ReminderError model to log errors when sending reminders and prevent retrying to send reminders that have failed multiple times
- [ ] Implement a purge system for deleting discord account not linked without activity for more than 3 months with no reminders
- [ ] Logging system for sent reminders and errors
- - [ ] Only purge a failing reminder if it has only one destination
- [ ] When sending a past reminder that has recurrences, calculate the next occurrence from now instead of sending all the missed occurrences
- [ ] Impelemnt emote support in reminders content

## Global

- [x] Dockerize the application
- [x] Build container during CI/CD
- [ ] Make a migration script from Kairos to Chronos
- [x] Create an alpha version for testing
- [ ] Create a beta version when the bot is ready
- [ ] Create a proper README
- [ ] Deploy the bot publicly

## Todo

- [ ] Forbid the user to create reminders in the past
- [ ] Clean components directory
- [ ] Being able to display the next month for the next reminder
