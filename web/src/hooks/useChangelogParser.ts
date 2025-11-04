/**
 * Hook to parse and structure changelog data
 */

export interface ChangelogEntry {
  section: string;
  subsection?: string;
  items: string[];
}

export interface ChangelogVersion {
  version: string;
  date: string;
  categories: {
    name: string;
    entries: ChangelogEntry[];
  }[];
}

export const useChangelogParser = () => {
  const parseChangelog = (): ChangelogVersion[] => {
    const changelogText = `# Changelog

All notable changes to this project will be documented in this file.

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
- Added \`/timezone display\` command to view current timezone.

## [Alpha 0.1.0] - 01-10-2025

### Major additions

#### Discord Bot

First release of the Discord bot with the following features:

- Added \`/timezone set/list\` commands to manage user timezones.
- Added \`/remindme\` command to set direct message reminders.
- Added \`/remindus\` command to set channel reminders.
- Implemented \`/profile\` command to view user profile and integrations.
- Introduced \`/reminders\` command with subcommands for managing reminders (delete, list, show, pause, restart).
- Added \`/calculatetime\` command for time calculations.
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
- Server is now built with a CI/CD pipeline using GitHub Actions.`;

    const versions: ChangelogVersion[] = [];
    const versionRegex = /## \[([^\]]+)\]\s*-\s*(.+)/g;

    let match;
    const versionMatches = [];

    while ((match = versionRegex.exec(changelogText)) !== null) {
      versionMatches.push({
        version: match[1],
        date: match[2],
        index: match.index,
      });
    }

    // Process each version
    for (let i = 0; i < versionMatches.length; i++) {
      const currentVersion = versionMatches[i];
      const nextVersionIndex =
        i + 1 < versionMatches.length
          ? versionMatches[i + 1].index
          : changelogText.length;

      const versionContent = changelogText.substring(
        currentVersion.index,
        nextVersionIndex
      );

      const categories = [];
      const categoryMatches = Array.from(
        versionContent.matchAll(/### (.+?)\n([\s\S]*?)(?=###|$)/g)
      );

      for (const categoryMatch of categoryMatches) {
        const categoryName = categoryMatch[1];
        const categoryBody = categoryMatch[2];

        const entries: ChangelogEntry[] = [];
        const sectionMatches = Array.from(
          categoryBody.matchAll(/#### (.+?)\n([\s\S]*?)(?=####|$)/g)
        );

        for (const sectionMatch of sectionMatches) {
          const sectionName = sectionMatch[1];
          const sectionBody = sectionMatch[2];

          const lines = sectionBody
            .split("\n")
            .filter((line) => line.trim().startsWith("-"));
          const items = lines.map((line) => {
            // Remove leading dashes and trim, handling nested items
            const item = line.replace(/^\s*-\s*/, "").trim();
            return item;
          });

          if (items.length > 0) {
            entries.push({
              section: sectionName,
              items: items,
            });
          }
        }

        // If no sections found, treat body as direct items
        if (entries.length === 0) {
          const lines = categoryBody
            .split("\n")
            .filter((line) => line.trim().startsWith("-"));
          const items = lines.map((line) =>
            line.replace(/^\s*-\s*/, "").trim()
          );

          if (items.length > 0) {
            entries.push({
              section: "",
              items: items,
            });
          }
        }

        if (entries.length > 0) {
          categories.push({
            name: categoryName,
            entries: entries,
          });
        }
      }

      versions.push({
        version: currentVersion.version,
        date: currentVersion.date,
        categories: categories,
      });
    }

    return versions;
  };

  return { parseChangelog };
};
