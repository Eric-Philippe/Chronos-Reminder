# Chronos Reminder — Android App

Native Android companion app for Chronos Reminder, built per [docs/ANDROID_SPEC.md](../../docs/ANDROID_SPEC.md).

- **Stack:** Kotlin 2.x · Jetpack Compose (Material 3) · Hilt · Retrofit + Kotlinx Serialization · Room · Coil · FCM
- **Min SDK:** 31 (Android 12) · **Target SDK:** 35

## Setup

1. **Open** this directory (`mobile/android`) in Android Studio (Ladybug or newer).
2. **Firebase:** drop your `google-services.json` into `app/`. It is gitignored. The build works
   without it (the `google-services` plugin is only applied when the file is present), but push
   notifications need it.
3. **Discord OAuth (optional):** set your client id in `app/build.gradle.kts`
   (`buildConfigField "DISCORD_CLIENT_ID"`). The backend must allow `chronos://auth/discord`
   as an OAuth redirect URI for the in-app Discord login to work.
4. **API base URL:** debug builds point to `http://10.0.2.2:8080/` (host machine `localhost`
   from the emulator). Set the release URL in `app/build.gradle.kts`.

## Build & test

```sh
./gradlew :app:assembleDebug      # build
./gradlew :app:testDebugUnitTest  # unit tests (MockK + Turbine)
./gradlew :app:connectedDebugAndroidTest  # instrumented tests (device/emulator required)
```

## Notes on spec vs. backend reality

The spec ([ANDROID_SPEC.md](../../docs/ANDROID_SPEC.md)) was written against an idealized API.
The app is built against the **actual** Go backend, which differs:

- Responses are **not** wrapped in `{ "data": ... }`; errors are `{ "error": "..." }`.
- `GET /api/reminders` returns `{ "reminders": [...], "count": n }`.
- Reminders carry `recurrence_type` (string) + `is_paused` (bool) instead of a packed int.
  The packed-int helpers from the spec are still used for the local Room cache.
- DFM: items have no server-side `position` update (no drag-reorder), and DFM reminders only
  support `discord_dm` / `email` destinations.
- API keys are created with a name only (scopes are assigned server-side).
- Discord OAuth callback takes `{ "code", "state" }` (no `redirect_uri`).

Backend features from spec §6 (FCM token endpoints, `android_push` destination type and
dispatcher) **do not exist yet**. The app already calls `POST/DELETE /api/fcm/token` and offers
the Push Notification destination; both degrade gracefully (logged, non-blocking) until the
backend lands.
