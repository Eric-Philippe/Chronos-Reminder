package com.chronos.reminder.core.database

import androidx.room.Dao
import androidx.room.Query
import androidx.room.Transaction
import androidx.room.Upsert
import com.chronos.reminder.dfm.data.DfmItemEntity
import kotlinx.coroutines.flow.Flow

@Dao
interface DfmDao {

    @Query("SELECT * FROM dfm_items ORDER BY position ASC, createdAt ASC")
    fun observeAll(): Flow<List<DfmItemEntity>>

    @Upsert
    suspend fun upsert(item: DfmItemEntity)

    @Upsert
    suspend fun upsertAll(items: List<DfmItemEntity>)

    @Query("DELETE FROM dfm_items WHERE id = :id")
    suspend fun deleteById(id: String)

    @Query("DELETE FROM dfm_items")
    suspend fun clear()

    @Transaction
    suspend fun replaceAll(items: List<DfmItemEntity>) {
        clear()
        upsertAll(items)
    }
}
