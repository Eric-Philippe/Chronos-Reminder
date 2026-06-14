package com.chronos.reminder.reminders.data

import com.chronos.reminder.reminders.domain.Destination
import com.chronos.reminder.reminders.domain.Reminder
import com.chronos.reminder.reminders.domain.RecurrenceType
import kotlinx.serialization.builtins.ListSerializer
import kotlinx.serialization.json.Json
import java.time.Instant

private val destinationsSerializer = ListSerializer(DestinationDto.serializer())

private fun parseInstant(iso: String?): Instant? =
    iso?.let { runCatching { Instant.parse(it) }.getOrNull() }

fun ReminderDto.toEntity(json: Json): ReminderEntity {
    val recurrence = RecurrenceType.fromApiString(recurrenceType)
    val pausedBit = if (isPaused) RecurrenceType.PAUSED_FLAG else 0
    return ReminderEntity(
        id = id,
        message = message,
        remindAtUtc = parseInstant(remindAtUtc)?.toEpochMilli() ?: 0L,
        nextFireUtc = parseInstant(nextFireUtc)?.toEpochMilli(),
        recurrenceState = recurrence.apiValue or pausedBit,
        isPaused = isPaused,
        destinationsJson = json.encodeToString(destinationsSerializer, destinations),
        createdAt = parseInstant(createdAt)?.toEpochMilli() ?: 0L,
        updatedAt = System.currentTimeMillis(),
    )
}

fun ReminderEntity.toDomain(json: Json): Reminder {
    val destinations = runCatching {
        json.decodeFromString(destinationsSerializer, destinationsJson)
    }.getOrDefault(emptyList())
    return Reminder(
        id = id,
        message = message,
        remindAtUtc = Instant.ofEpochMilli(remindAtUtc),
        nextFireUtc = nextFireUtc?.let(Instant::ofEpochMilli),
        recurrence = RecurrenceType.fromState(recurrenceState),
        isPaused = isPaused,
        destinations = destinations.map { Destination(it.type, it.metadata) },
        createdAt = createdAt.takeIf { it > 0 }?.let(Instant::ofEpochMilli),
    )
}
