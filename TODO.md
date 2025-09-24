# Chronos v1.0 Roadmap

## Database

- [ ] Add a `last_notified_at` column to the `reminders` table to store the last time a reminder was sent
- [ ] Add a logging table to track reminder deliveries and failures
- [ ] Use REDIS ?

## Discord Bot

- [x] `/remindme` command -> Quickly set {DISCORD_DM}Reminders type
- [x] `/remindus` command -> Set {DISCORD_CHANNEL}Reminders type, (MUST BE ENTERED IN A SERVER, NOT IN A DM), the channel field is an autocomplete that gives the user channels select where the user has the manage channel permission
- [ ] Draw the discord reminder with a prettier display (image...)
- [ ] `/profile` command -> View user profile information and their integrations and buttons to manage them
- [ ] `/reminders delete|list|show|pause|unpause` commands -> Manage reminders
- - [ ] `list` subcommand -> Give a first embed with all the list, and a second embed with a paginated list of one reminders that the user can scroll through
- - [ ] `delete` subcommand -> Delete a reminder by its content/ID, or if not given, launch a select menu to first select between all its different type of reminders (DM, Channel, Webhook...) or skip this step if only one type exists, then display a second select menu with all the reminders of this type to select which one to delete, and then remove the message and put a confirmation message. Only when targeting a {DISCORD_CHANNEL}Reminder, If the user is admin of the current guild, even if a reminder was created by someone else in this guild, he can delete it.
- - [ ] `show` subcommand -> Get a reminder by its content/ID and display it in a pretty embed or image, and little icons to show its types/destination(s) (footer with: create website account for update/delete reminders)
- - [ ] `pause` subcommand -> Pause a reminder by its content/ID, return an error message if the reminder is not a recurring one.
- - [ ] `unpause` subcommand -> Unpause a reminder by its content/ID, return an error message if the reminder is not a recurring one.
- [ ] `/calcultime` command -> Time calculator
- [ ] `/help` command -> Help command with all the commands and their descriptions, and buttons to get more information on each command
- [ ] Test behavior when the bot can't send the reminder DM/CHANNEL (user blocked the bot, user left the server, bot kicked from the server, no permission to send messages in the channel...)

## API

- [ ] Create API endpoints for the web application to interact with the reminder engine
- [ ] Implement authentication for API endpoints
- [ ] API Key for third-party integrations

## Web Application

- [ ] https://chronosreminder.app/.fr?
- [ ] Welcome unsigned page
- [ ] Create account/login => (Possibility to merge from an already existing Discord account => OAuth2 Discord => Create account prefiled with Discord info)
- [ ] CRUD on reminders
- [ ] View profile and integrations
- [ ] Link Discord account
- [ ] View history/logs
- [ ] Same slick style as SnapFileThing

## Reminder engine

- [x] Implement main scheduler queue
- [x] Implement DM_DISCORD dispatcher
- [ ] Implement CHANNEL_DISCORD dispatcher
- [ ] Implement WEBHOOK dispatcher
- [ ] Implement a purge system for deleting discord account not linked without activity for more than 3 months with no reminders
- [ ] Logging system for sent reminders and errors

## Global

- [ ] Dockerize the application
- [ ] Build container during CI/CD
- [ ] Make a migration script from Kairos to Chronos
- [ ] Create a beta version when the bot is ready
- [ ] Create a proper README
- [ ] Deploy the bot publicly
