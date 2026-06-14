package com.chronos.reminder.core.database

import androidx.room.Database
import androidx.room.RoomDatabase
import com.chronos.reminder.dfm.data.DfmItemEntity
import com.chronos.reminder.reminders.data.ReminderEntity

@Database(
    entities = [ReminderEntity::class, DfmItemEntity::class],
    version = 1,
    exportSchema = false,
)
abstract class ChronosDatabase : RoomDatabase() {
    abstract fun reminderDao(): ReminderDao
    abstract fun dfmDao(): DfmDao
}
