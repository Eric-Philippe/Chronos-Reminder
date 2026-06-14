package com.chronos.reminder.core.ui.theme

import androidx.compose.ui.graphics.Color

// --- Dark theme (primary) ---
val BackgroundMain = Color(0xFF1E1A17) // oklch(0.15 0.01 49.25) — main background
val BackgroundCard = Color(0xFF2C2723) // oklch(0.22 0.01 49.25) — card / surface
val BackgroundMuted = Color(0xFF3A342E) // slightly lighter surface, e.g. input fields
val ForegroundMain = Color(0xFFF2F1EE) // oklch(0.95 0.001 106) — primary text
val ForegroundMuted = Color(0xFF9C9489) // secondary / hint text
val AccentOrange = Color(0xFFC47A3A) // oklch(0.65 0.15 40) — CTA buttons, highlights
val AccentOrangeDark = Color(0xFF7A4E2A) // oklch(0.40 0.10 40) — pressed state
val DestructiveRed = Color(0xFFC0392B) // oklch(0.577 0.245 27) — delete / error
val BorderColor = Color(0xFF4A433C) // subtle border / divider

// --- Light theme (optional, lower priority) ---
val LightBackground = Color(0xFFF9F7F4)
val LightCard = Color(0xFFFFFFFF)
val LightForeground = Color(0xFF2C2420)
val LightAccent = Color(0xFF7A4E2A)

// --- Recurrence badge colors ---
val RecurrenceHourly = Color(0xFF4A90D9)
val RecurrenceWorkdays = Color(0xFF27AE60)
val RecurrenceWeekend = Color(0xFF8E44AD)
val RecurrenceWeekly = Color(0xFFD4AC0D)
val RecurrenceMonthly = Color(0xFF2980B9)
val RecurrenceYearly = Color(0xFFE67E22)
val LinkedGreen = Color(0xFF27AE60)
