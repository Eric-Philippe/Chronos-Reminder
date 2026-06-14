package com.chronos.reminder.dfm.data

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class DfmNoteDto(
    val id: String,
    @SerialName("remind_at_utc") val remindAtUtc: String? = null,
    @SerialName("next_fire_utc") val nextFireUtc: String? = null,
    @SerialName("recurrence_type") val recurrenceType: String? = null,
    @SerialName("has_reminder") val hasReminder: Boolean = false,
    val destinations: List<String> = emptyList(), // "discord_dm" | "email"
    val items: List<DfmItemDto> = emptyList(),
)

@Serializable
data class DfmItemDto(
    val id: String,
    val content: String = "",
    val checked: Boolean = false,
    val position: Int = 0,
    @SerialName("created_at") val createdAt: String? = null,
)

@Serializable
data class AddDfmItemRequest(val content: String)

@Serializable
data class UpdateDfmItemRequest(
    val content: String? = null,
    val checked: Boolean? = null,
)

@Serializable
data class DfmReminderRequest(
    val date: String = "",
    val time: String,
    val recurrence: String,
    val destinations: List<String>, // "discord_dm" and/or "email"
)
