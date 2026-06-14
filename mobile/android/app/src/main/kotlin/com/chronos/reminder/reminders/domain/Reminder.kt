package com.chronos.reminder.reminders.domain

import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.JsonPrimitive
import java.time.Instant

data class Destination(
    val type: String,
    val metadata: JsonObject,
) {
    // Best-effort human label for a metadata value (e.g. email address, channel name)
    fun metadataValue(key: String): String? =
        (metadata[key] as? JsonPrimitive)?.content

    companion object {
        const val TYPE_DISCORD_DM = "discord_dm"
        const val TYPE_DISCORD_CHANNEL = "discord_channel"
        const val TYPE_EMAIL = "email"
        const val TYPE_WEBHOOK = "webhook"
        const val TYPE_ANDROID_PUSH = "android_push"
    }
}

data class Reminder(
    val id: String,
    val message: String,
    val remindAtUtc: Instant,
    val nextFireUtc: Instant?,
    val recurrence: RecurrenceType,
    val isPaused: Boolean,
    val destinations: List<Destination>,
    val createdAt: Instant?,
)
