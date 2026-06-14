package com.chronos.reminder.reminders.data

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "reminders")
data class ReminderEntity(
    @PrimaryKey val id: String,
    val message: String,
    val remindAtUtc: Long, // epoch ms
    val nextFireUtc: Long?,
    val recurrenceState: Int, // type in low 7 bits, paused flag in bit 8
    val isPaused: Boolean,
    val destinationsJson: String, // serialized JSON of destinations list
    val createdAt: Long,
    val updatedAt: Long,
)
