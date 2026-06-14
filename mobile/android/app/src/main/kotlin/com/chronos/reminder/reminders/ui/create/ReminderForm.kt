package com.chronos.reminder.reminders.ui.create

import com.chronos.reminder.reminders.data.CreateReminderRequest
import com.chronos.reminder.reminders.data.DestinationRequest
import com.chronos.reminder.reminders.domain.Destination
import com.chronos.reminder.reminders.domain.RecurrenceType
import com.chronos.reminder.reminders.domain.Reminder
import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.JsonPrimitive
import kotlinx.serialization.json.jsonPrimitive
import java.time.LocalDate
import java.time.LocalTime
import java.time.ZoneId
import java.time.ZonedDateTime
import java.time.format.DateTimeFormatter

const val MAX_MESSAGE_LENGTH = 500

data class FormDestination(
    val type: String,
    val metadata: Map<String, String>,
    // Extra context shown on the chip, e.g. email address or channel name
    val detail: String? = null,
) {
    fun toRequest() = DestinationRequest(
        type = type,
        metadata = JsonObject(metadata.mapValues { JsonPrimitive(it.value) }),
    )
}

data class ReminderForm(
    val date: LocalDate? = null,
    val time: LocalTime? = null,
    val recurrence: RecurrenceType = RecurrenceType.ONCE,
    val message: String = "",
    val destinations: List<FormDestination> = emptyList(),
) {
    fun isDateTimeInFuture(userTz: String): Boolean {
        val d = date ?: return false
        val t = time ?: return false
        val zone = runCatching { ZoneId.of(userTz) }.getOrDefault(ZoneId.systemDefault())
        return ZonedDateTime.of(d, t, zone).isAfter(ZonedDateTime.now(zone))
    }

    fun toRequest(): CreateReminderRequest = CreateReminderRequest(
        date = requireNotNull(date).format(DateTimeFormatter.ISO_LOCAL_DATE),
        time = requireNotNull(time).format(DateTimeFormatter.ofPattern("HH:mm")),
        message = message.trim(),
        recurrence = recurrence.apiString,
        destinations = destinations.map { it.toRequest() },
    )

    companion object {
        fun fromReminder(reminder: Reminder, userTz: String): ReminderForm {
            val zone = runCatching { ZoneId.of(userTz) }.getOrDefault(ZoneId.systemDefault())
            val zdt = ZonedDateTime.ofInstant(reminder.nextFireUtc ?: reminder.remindAtUtc, zone)
            return ReminderForm(
                date = zdt.toLocalDate(),
                time = zdt.toLocalTime(),
                recurrence = reminder.recurrence,
                message = reminder.message,
                destinations = reminder.destinations.map { it.toFormDestination() },
            )
        }

        private fun Destination.toFormDestination(): FormDestination {
            val flat = metadata.entries.mapNotNull { (key, value) ->
                runCatching { key to value.jsonPrimitive.content }.getOrNull()
            }.toMap()
            val detail = when (type) {
                Destination.TYPE_EMAIL -> flat["email"]
                Destination.TYPE_WEBHOOK -> flat["url"]
                Destination.TYPE_DISCORD_CHANNEL -> flat["channel_name"] ?: flat["channel_id"]
                else -> null
            }
            return FormDestination(type, flat, detail)
        }
    }
}
