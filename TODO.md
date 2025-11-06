# TODOs

## Chronos v1.0 Roadmap

### Discord Bot

- [ ] Prevent snoozing a reminder after its next recurrence time
- [ ] Add embed footer with the webapp link

### Web Application

- [ ] Complete the footer with proper links and information
- [ ] Complete the Changelog page
- [ ] Bot documentation page
- [ ] Self-hosting guide page
- [ ] Terms of Service and Privacy Policy pages

### System

- [ ] Backup strategy for the database

### 1.0 Launch

- [ ] Create a proper README
- [ ] Create the official Discord server for support and community
- [ ] More complex branching strategy for Git (e.g., develop, staging, production branches)
- [ ] Deploy the bot publicly
- - [ ] Make the bot verified
- - [ ] Upload updated privacy policy and terms of service
- - [ ] Update bot description and assets
- - [ ] Update Discord Discovery listing
- - [ ] Publish in bot listing websites

## Chronos v1.1 Roadmap

### Discord Bot

- [ ] Test behavior when the bot can't send the reminder DM/CHANNEL (user blocked the bot, user left the server, bot kicked from the server, no permission to send messages in the channel...)

### Server API

- [ ] Create more API endpoints for the web application to interact with the reminder engine

### Reminder engine

- [ ] Logging system for sent reminders
- [ ] Only purge a failing reminder if it has only one destination

### Global

- [ ] Add email support
- [ ] Make a migration script from Kairos to Chronos
- [ ] Badge shop
