package com.chronos.reminder.reminders.data

import com.chronos.reminder.reminders.domain.Destination
import com.chronos.reminder.reminders.domain.RecurrenceType
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.JsonPrimitive
import org.junit.Assert.assertEquals
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Test
import java.time.Instant

class ReminderMappersTest {

    private val json = Json { ignoreUnknownKeys = true }

    private val dto = ReminderDto(
        id = "abc-123",
        message = "Buy groceries",
        remindAtUtc = "2026-07-15T13:30:00Z",
        nextFireUtc = "2026-07-16T13:30:00Z",
        recurrenceType = "DAILY",
        isPaused = true,
        destinations = listOf(
            DestinationDto(
                id = "d1",
                type = Destination.TYPE_DISCORD_DM,
                metadata = JsonObject(mapOf("user_id" to JsonPrimitive("42"))),
            ),
        ),
        createdAt = "2026-06-01T08:00:00Z",
    )

    @Test
    fun `dto maps to entity with packed recurrence state`() {
        val entity = dto.toEntity(json)

        assertEquals("abc-123", entity.id)
        assertEquals(Instant.parse("2026-07-15T13:30:00Z").toEpochMilli(), entity.remindAtUtc)
        assertEquals(Instant.parse("2026-07-16T13:30:00Z").toEpochMilli(), entity.nextFireUtc)
        assertTrue(entity.isPaused)
        // DAILY = 4, paused flag = bit 8
        assertEquals(RecurrenceType.DAILY.apiValue or RecurrenceType.PAUSED_FLAG, entity.recurrenceState)
    }

    @Test
    fun `entity maps to domain with destinations round-tripped`() {
        val domain = dto.toEntity(json).toDomain(json)

        assertEquals("abc-123", domain.id)
        assertEquals("Buy groceries", domain.message)
        assertEquals(RecurrenceType.DAILY, domain.recurrence)
        assertTrue(domain.isPaused)
        assertEquals(1, domain.destinations.size)
        assertEquals(Destination.TYPE_DISCORD_DM, domain.destinations[0].type)
        assertEquals("42", domain.destinations[0].metadataValue("user_id"))
    }

    @Test
    fun `unknown recurrence string falls back to ONCE`() {
        val entity = dto.copy(recurrenceType = "SOMETIMES", isPaused = false).toEntity(json)
        val domain = entity.toDomain(json)

        assertEquals(RecurrenceType.ONCE, domain.recurrence)
    }

    @Test
    fun `missing optional fields map to nulls`() {
        val minimal = ReminderDto(
            id = "x",
            message = "m",
            remindAtUtc = "2026-07-15T13:30:00Z",
        )
        val domain = minimal.toEntity(json).toDomain(json)

        assertNull(domain.nextFireUtc)
        assertNull(domain.createdAt)
        assertEquals(RecurrenceType.ONCE, domain.recurrence)
    }

    @Test
    fun `recurrence state helpers decode packed int`() {
        val state = RecurrenceType.WEEKLY.apiValue or RecurrenceType.PAUSED_FLAG
        assertTrue(RecurrenceType.isPaused(state))
        assertEquals(RecurrenceType.WEEKLY, RecurrenceType.fromState(state))
    }
}
