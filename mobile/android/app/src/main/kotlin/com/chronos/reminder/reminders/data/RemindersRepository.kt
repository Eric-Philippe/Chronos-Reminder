package com.chronos.reminder.reminders.data

import com.chronos.reminder.core.database.ReminderDao
import com.chronos.reminder.core.network.ApiResult
import com.chronos.reminder.core.network.safeApiCall
import com.chronos.reminder.reminders.domain.Reminder
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.serialization.json.Json
import javax.inject.Inject
import javax.inject.Singleton

interface RemindersRepository {
    fun getReminders(): Flow<List<Reminder>> // from Room, always up-to-date
    fun getReminder(id: String): Flow<Reminder?>
    suspend fun refreshReminders(): ApiResult<Unit> // fetch from API → update Room
    suspend fun refreshReminder(id: String): ApiResult<Unit>
    suspend fun createReminder(request: CreateReminderRequest): ApiResult<Reminder>
    suspend fun updateReminder(id: String, request: CreateReminderRequest): ApiResult<Reminder>
    suspend fun deleteReminder(id: String): ApiResult<Unit>
    suspend fun pauseReminder(id: String): ApiResult<Unit>
    suspend fun resumeReminder(id: String): ApiResult<Unit>
    suspend fun duplicateReminder(id: String): ApiResult<Reminder>
}

@Singleton
class RemindersRepositoryImpl @Inject constructor(
    private val api: RemindersApi,
    private val dao: ReminderDao,
    private val json: Json,
) : RemindersRepository {

    override fun getReminders(): Flow<List<Reminder>> =
        dao.observeAll().map { entities -> entities.map { it.toDomain(json) } }

    override fun getReminder(id: String): Flow<Reminder?> =
        dao.observeById(id).map { it?.toDomain(json) }

    override suspend fun refreshReminders(): ApiResult<Unit> =
        safeApiCall { api.getReminders() }.map { response ->
            dao.replaceAll(response.reminders.map { it.toEntity(json) })
        }

    override suspend fun refreshReminder(id: String): ApiResult<Unit> =
        safeApiCall { api.getReminder(id) }.map { dto ->
            dao.upsert(dto.toEntity(json))
        }

    override suspend fun createReminder(request: CreateReminderRequest): ApiResult<Reminder> =
        upsertFromApi { api.createReminder(request) }

    override suspend fun updateReminder(id: String, request: CreateReminderRequest): ApiResult<Reminder> =
        upsertFromApi { api.updateReminder(id, request) }

    override suspend fun deleteReminder(id: String): ApiResult<Unit> =
        safeApiCall { api.deleteReminder(id) }.map { dao.deleteById(id) }

    override suspend fun pauseReminder(id: String): ApiResult<Unit> =
        upsertFromApi { api.pauseReminder(id) }.map { }

    override suspend fun resumeReminder(id: String): ApiResult<Unit> =
        upsertFromApi { api.resumeReminder(id) }.map { }

    override suspend fun duplicateReminder(id: String): ApiResult<Reminder> =
        upsertFromApi { api.duplicateReminder(id) }

    private suspend fun upsertFromApi(call: suspend () -> retrofit2.Response<ReminderDto>): ApiResult<Reminder> {
        val result = safeApiCall(call)
        if (result is ApiResult.Success) {
            val entity = result.data.toEntity(json)
            dao.upsert(entity)
            return ApiResult.Success(entity.toDomain(json))
        }
        @Suppress("UNCHECKED_CAST")
        return result as ApiResult<Reminder>
    }
}
