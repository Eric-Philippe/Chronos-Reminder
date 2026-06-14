package com.chronos.reminder.dfm.data

import com.chronos.reminder.core.database.DfmDao
import com.chronos.reminder.core.network.ApiResult
import com.chronos.reminder.core.network.safeApiCall
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.map
import java.time.Instant
import javax.inject.Inject
import javax.inject.Singleton

data class DfmItem(
    val id: String,
    val content: String,
    val checked: Boolean,
    val position: Int,
)

data class DfmReminderInfo(
    val nextFireUtc: Instant?,
    val recurrence: String,
    val destinations: List<String>,
)

@Singleton
class DfmRepository @Inject constructor(
    private val api: DfmApi,
    private val dao: DfmDao,
) {

    // Reminder config has no Room table; it is small and refreshed with the note.
    private val _reminderInfo = MutableStateFlow<DfmReminderInfo?>(null)
    val reminderInfo: StateFlow<DfmReminderInfo?> = _reminderInfo.asStateFlow()

    fun getItems(): Flow<List<DfmItem>> =
        dao.observeAll().map { entities ->
            entities.map { DfmItem(it.id, it.content, it.checked, it.position) }
        }

    suspend fun refresh(): ApiResult<Unit> =
        safeApiCall { api.getNote() }.map { note -> applyNote(note) }

    suspend fun addItem(content: String): ApiResult<Unit> {
        val result = safeApiCall { api.addItem(AddDfmItemRequest(content)) }
        if (result is ApiResult.Success) {
            dao.upsert(result.data.toEntity())
        }
        return result.map { }
    }

    suspend fun updateItem(id: String, content: String? = null, checked: Boolean? = null): ApiResult<Unit> {
        val result = safeApiCall { api.updateItem(id, UpdateDfmItemRequest(content, checked)) }
        if (result is ApiResult.Success) {
            dao.upsert(result.data.toEntity())
        }
        return result.map { }
    }

    // Optimistic checkbox toggle: Room first, revert on API error
    suspend fun setItemChecked(item: DfmItem, checked: Boolean): ApiResult<Unit> {
        dao.upsert(item.copy(checked = checked).toEntity())
        val result = safeApiCall { api.updateItem(item.id, UpdateDfmItemRequest(checked = checked)) }
        if (result !is ApiResult.Success) {
            dao.upsert(item.toEntity())
        }
        return result.map { }
    }

    suspend fun deleteItem(id: String): ApiResult<Unit> =
        safeApiCall { api.deleteItem(id) }.map { dao.deleteById(id) }

    suspend fun setReminder(request: DfmReminderRequest): ApiResult<Unit> =
        safeApiCall { api.setReminder(request) }.map { note -> applyNote(note) }

    suspend fun removeReminder(): ApiResult<Unit> {
        val result = safeApiCall { api.removeReminder() }
        if (result is ApiResult.Success) {
            _reminderInfo.value = null
        }
        return result.map { }
    }

    suspend fun sendNow(): ApiResult<Unit> =
        safeApiCall { api.sendNow() }.map { }

    private suspend fun applyNote(note: DfmNoteDto) {
        dao.replaceAll(note.items.map { it.toEntity() })
        _reminderInfo.value = if (note.hasReminder) {
            DfmReminderInfo(
                nextFireUtc = note.nextFireUtc?.let { runCatching { Instant.parse(it) }.getOrNull() },
                recurrence = note.recurrenceType ?: "ONCE",
                destinations = note.destinations,
            )
        } else {
            null
        }
    }

    private fun DfmItemDto.toEntity() = DfmItemEntity(
        id = id,
        content = content,
        checked = checked,
        position = position,
        createdAt = createdAt?.let { runCatching { Instant.parse(it).toEpochMilli() }.getOrNull() } ?: 0L,
    )

    private fun DfmItem.toEntity() = DfmItemEntity(
        id = id,
        content = content,
        checked = checked,
        position = position,
        createdAt = 0L,
    )
}
