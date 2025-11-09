# TODOs

## Chronos v1.1 Roadmap

### Web Application

- [ ] Terms of Service and Privacy Policy pages

### Discord Bot

- [ ] Test behavior when the bot can't send the reminder DM/CHANNEL (user blocked the bot, user left the server, bot kicked from the server, no permission to send messages in the channel...)
- [ ] Prevent user from snoozing a reminder if the next occurrence is before the snooze time

### Server API

- [ ] "Zombie" account purger
- [ ] Create more API endpoints for the web application to interact with the reminder engine

### Reminder engine

- [ ] Logging system for sent reminders
- [ ] Only purge a failing reminder if it has only one destination

### Global

- [ ] Put Email, Discord Invite Link, Version in the config file instead of hardcoding them
- [ ] Add email support
- [ ] Add API Key badge to user profiles
- [ ] Badge shop
