# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - 9/11/2025

Final release of Chronos Reminder with the following features:

### Major additions

- Complete home page with 3D models, descriptions and app details.
- Header UI improvements
- Footer UI improvements
- Reminders dashboard with pagination, search and filters
- Changelog page
- Contact page with form submission
- Account settings page
- API Key management
- Rate limiting for API calls
- Status page with uptime kuma integration
- Documentation site with guides and API references with GitBook
- Reset password flow
- Email verification flow
- Official Discord server for support and community
- Reminders getting paused then later resumed don't try to resend missed reminders, instead just continue from the next occurrence.
- Complete containerization and deployment scripts for production use.
- New Bot commands:
  - `/hourglass` launch a quick timer
  - `/support` get the link to the official support server
  - `/help` complete help command listing all available commands and their usage.
  - `/timezone display` show the user's current timezone

### Minor additions

- Improved responsiveness for mobile devices.
- Performance optimizations on main engine
- User can now delete their account from settings.
- Minor translation fixes and improvements.

### Fixes

- Reminders editing now properly updates the reminder without putting back the old values.
- Fixed various UI bugs and layout issues.
- Fixed timezone handling issues in the reminders creation flow.

## [Alpha 0.2.0] - 20-10-2025

### Major additions

##### Web App

First release of the web application with the following features:

- User authentication
- - Login with email and password
- - OAuth2 login with Discord
- API authentication using JWT tokens and cookie sessions
- Dark/Light mode toggle
- French, English and Spanish language support
- First version of the Header and Footer components
- Main reminders dashboard
- Discord DM Reminders creation
- CRUD operations for reminders

### Minor additions

#### Discord Bot

- Now display pretty /profile embeds with onfly drawn graphs.
- Added `/timezone display` command to view current timezone.

## [Alpha 0.1.0] - 01-10-2025

### Major additions

#### Discord Bot

First release of the Discord bot with the following features:

- Added `/timezone set/list` commands to manage user timezones.
- Added `/remindme` command to set direct message reminders.
- Added `/remindus` command to set channel reminders.
- Implemented `/profile` command to view user profile and integrations.
- Introduced `/reminders` command with subcommands for managing reminders (delete, list, show, pause, restart).
- Added `/calculatetime` command for time calculations.
- Enabled snooze functionality for direct message reminders.

#### Engine

- Implemented the main scheduler for handling reminders.
- Implemented garbage collection for old reminders.
- Implemented CHANNEL_DISCORD dispatcher with role mention support.
- Implemented DM_DISCORD dispatcher.
- Implemented recurrence handling for reminders.
- Added recalculation of next occurrence on restart to avoid spam.
- Added a full timezone handling system.
- Implemented testing for the recalculation of next occurrences and the datetime input parser.
- Added Swagger documentation for the REST API.
- Database table storing reminder errors to log failed reminder deliveries.

#### System

- Cache Sytem: Added caching for user account data to reduce database load.
- Whole server is dockerized for easier deployment and development.
- Server is now built with a CI/CD pipeline using GitHub Actions.
