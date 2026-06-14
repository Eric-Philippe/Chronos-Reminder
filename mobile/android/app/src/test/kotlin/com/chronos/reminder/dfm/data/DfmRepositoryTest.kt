package com.chronos.reminder.dfm.data

import com.chronos.reminder.core.database.DfmDao
import com.chronos.reminder.core.network.ApiResult
import io.mockk.coEvery
import io.mockk.coVerifyOrder
import io.mockk.mockk
import kotlinx.coroutines.test.runTest
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.ResponseBody.Companion.toResponseBody
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test
import retrofit2.Response

class DfmRepositoryTest {

    private val api: DfmApi = mockk()
    private val dao: DfmDao = mockk(relaxed = true)
    private lateinit var repository: DfmRepository

    private val item = DfmItem(id = "i1", content = "Milk", checked = false, position = 0)

    @Before
    fun setUp() {
        repository = DfmRepository(api, dao)
    }

    @Test
    fun `checkbox toggle is optimistic`() = runTest {
        coEvery { api.updateItem("i1", any()) } returns Response.success(
            DfmItemDto(id = "i1", content = "Milk", checked = true, position = 0),
        )

        val result = repository.setItemChecked(item, checked = true)

        assertTrue(result is ApiResult.Success)
        coVerifyOrder {
            dao.upsert(match { it.id == "i1" && it.checked }) // optimistic write
        }
    }

    @Test
    fun `failed toggle reverts the optimistic write`() = runTest {
        coEvery { api.updateItem("i1", any()) } returns Response.error(
            500,
            """{"error":"boom"}""".toResponseBody("application/json".toMediaType()),
        )

        val result = repository.setItemChecked(item, checked = true)

        assertTrue(result is ApiResult.Error)
        coVerifyOrder {
            dao.upsert(match { it.id == "i1" && it.checked }) // optimistic write
            dao.upsert(match { it.id == "i1" && !it.checked }) // revert
        }
    }

    @Test
    fun `refresh replaces items and exposes reminder info`() = runTest {
        coEvery { api.getNote() } returns Response.success(
            DfmNoteDto(
                id = "note-1",
                hasReminder = true,
                nextFireUtc = "2026-07-15T09:00:00Z",
                recurrenceType = "DAILY",
                destinations = listOf("discord_dm"),
                items = listOf(DfmItemDto(id = "i1", content = "Milk")),
            ),
        )

        val result = repository.refresh()

        assertTrue(result is ApiResult.Success)
        assertTrue(repository.reminderInfo.value != null)
        coVerifyOrder { dao.replaceAll(match { it.size == 1 && it[0].id == "i1" }) }
    }
}
