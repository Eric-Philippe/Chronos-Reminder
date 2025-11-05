# Chronos v1.0 Roadmap

## Discord Bot

- [ ] `/help` command -> Help command with all the commands and their descriptions, and buttons to get more information on each command
- [ ] Test behavior when the bot can't send the reminder DM/CHANNEL (user blocked the bot, user left the server, bot kicked from the server, no permission to send messages in the channel...)
- [ ] Prevent snoozing a reminder after its next recurrence time
- [ ] Add embed footer with the webapp link

## Server API

- [ ] Create API endpoints for the web application to interact with the reminder engine

## Web Application

### Home

- [ ] Add more information about the service and its features
- [ ] Add screenshots or demo videos of the web application and Discord bot

### Help

- - [ ] Use Chronos Bot
- - [ ] Use Web Application
- - [ ] Use API Key
- - [ ] Self-hosting guide
- - [ ] What's New

### Layout

- [ ] Add proper footer links and information, github, version, privacy policy, terms of service...

## Reminder engine

- [ ] Logging system for sent reminders
- - [ ] Only purge a failing reminder if it has only one destination

## System

- [ ] Make a migration script from Kairos to Chronos
- [ ] Backup the database regularly

## 1.0 Launch

- [ ] Create a proper README
- [ ] Create the official Discord server for support and community
- [ ] More complex branching strategy for Git (e.g., develop, staging, production branches)
- [ ] Deploy the bot publicly
- - [ ] Make the bot verified
- - [ ] Upload updated privacy policy and terms of service
- - [ ] Update bot description and assets
- - [ ] Update Discord Discovery listing
- - [ ] Publish in bot listing websites
