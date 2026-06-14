package com.chronos.reminder.core.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.unit.dp

val CornerRadius = 8.dp

// Always dark, no dynamic color
@Composable
fun ChronosTheme(content: @Composable () -> Unit) {
    MaterialTheme(
        colorScheme = darkColorScheme(
            background = BackgroundMain,
            surface = BackgroundCard,
            surfaceContainer = BackgroundCard,
            surfaceContainerHigh = BackgroundCard,
            surfaceContainerHighest = BackgroundMuted,
            primary = AccentOrange,
            onPrimary = ForegroundMain,
            secondary = AccentOrangeDark,
            onSecondary = ForegroundMain,
            onBackground = ForegroundMain,
            onSurface = ForegroundMain,
            onSurfaceVariant = ForegroundMuted,
            error = DestructiveRed,
            onError = ForegroundMain,
            outline = BorderColor,
            outlineVariant = BorderColor,
            surfaceVariant = BackgroundMuted,
        ),
        typography = ChronosTypography,
        content = content,
    )
}
