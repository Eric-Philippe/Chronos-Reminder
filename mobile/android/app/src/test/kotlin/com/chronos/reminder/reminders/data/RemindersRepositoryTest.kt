package com.chronos.reminder.reminders.data

import com.chronos.reminder.core.database.ReminderDao
import com.chronos.reminder.core.network.ApiResult
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.mockk
import kotlinx.coroutines.test.runTest
import kotlinx.serialization.json.Json
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.ResponseBody.Companion.toResponseBody
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test
import retrofit2.Response

class RemindersRepositoryTest {

    private val api: RemindersApi = mockk()
    private val dao: ReminderDao = mockk(relaxed = true)
    private val json = Json { ignoreUnknownKeys = true }
    private lateinit var repository: RemindersRepositoryImpl

    private val dto = ReminderDto(
        id = "r1",
        message = "Water plants",
        remindAtUtc = "2026-07-15T13:30:00Z",
        recurrenceType = "WEEKLY",
    )

    @Before
    fun setUp() {
        repository = RemindersRepositoryImpl(api, dao, json)
    }

    @Test
    fun `refresh success replaces room cache`() = runTest {
        coEvery { api.getReminders() } returns Response.success(RemindersListResponse(listOf(dto), 1))

        val result = repository.refreshReminders()

        assertTrue(result is ApiResult.Success)
        coVerify { dao.replaceAll(match { it.size == 1 && it[0].id == "r1" }) }
    }

    @Test
    fun `refresh failure surfaces api error and skips room write`() = runTest {
        coEvery { api.getReminders() } returns Response.error(
            401,
            """{"error":"Unauthorized"}""".toResponseBody("application/json".toMediaType()),
        )

        val result = repository.refreshReminders()

        assertTrue(result is ApiResult.Error)
        assertEquals(401, (result as ApiResult.Error).code)
        assertEquals("Unauthorized", result.message)
        coVerify(exactly = 0) { dao.replaceAll(any()) }
    }

    @Test
    fun `create success upserts returned reminder`() = runTest {
        coEvery { api.createReminder(any()) } returns Response.success(dto)

        val request = CreateReminderRequest(
            date = "2026-07-15",
            time = "13:30",
            message = "Water plants",
            recurrence = "WEEKLY",
            destinations = emptyList(),
        )
        val result = repository.createReminder(request)

        assertTrue(result is ApiResult.Success)
        assertEquals("Water plants", (result as ApiResult.Success).data.message)
        coVerify { dao.upsert(match { it.id == "r1" }) }
    }

    @Test
    fun `delete success removes from room`() = runTest {
        coEvery { api.deleteReminder("r1") } returns Response.success(MessageResponse("deleted"))

        val result = repository.deleteReminder("r1")

        assertTrue(result is ApiResult.Success)
        coVerify { dao.deleteById("r1") }
    }

    @Test
    fun `network failure maps to NetworkError`() = runTest {
        coEvery { api.getReminders() } throws java.io.IOException("offline")

        val result = repository.refreshReminders()

        assertTrue(result is ApiResult.NetworkError)
    }
}
