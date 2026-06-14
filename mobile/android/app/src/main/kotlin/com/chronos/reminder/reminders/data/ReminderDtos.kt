package com.chronos.reminder.reminders.data

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonObject

@Serializable
data class ReminderDto(
    val id: String,
    val message: String = "",
    @SerialName("remind_at_utc") val remindAtUtc: String,
    @SerialName("next_fire_utc") val nextFireUtc: String? = null,
    // The backend sends the recurrence as a name string plus a separate paused flag.
    @SerialName("recurrence_type") val recurrenceType: String? = null,
    @SerialName("is_paused") val isPaused: Boolean = false,
    val destinations: List<DestinationDto> = emptyList(),
    @SerialName("created_at") val createdAt: String? = null,
)

@Serializable
data class DestinationDto(
    val id: String? = null,
    val type: String,
    val metadata: JsonObject = JsonObject(emptyMap()),
)

@Serializable
data class RemindersListResponse(
    val reminders: List<ReminderDto> = emptyList(),
    val count: Int = 0,
)

@Serializable
data class CreateReminderRequest(
    val date: String, // "YYYY-MM-DD" in the user's timezone
    val time: String, // "HH:mm"
    val message: String,
    val recurrence: String, // "ONCE" | "DAILY" | ...
    val destinations: List<DestinationRequest>,
)

@Serializable
data class DestinationRequest(
    val type: String,
    val metadata: JsonObject,
)

@Serializable
data class MessageResponse(
    val message: String? = null,
)

@Serializable
data class SnoozeReminderRequest(
    val minutes: Int,
)
