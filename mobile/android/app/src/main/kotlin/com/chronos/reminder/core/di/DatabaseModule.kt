package com.chronos.reminder.core.di

import android.content.Context
import androidx.room.Room
import com.chronos.reminder.core.database.ChronosDatabase
import com.chronos.reminder.core.database.DfmDao
import com.chronos.reminder.core.database.ReminderDao
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object DatabaseModule {

    @Provides
    @Singleton
    fun provideDatabase(@ApplicationContext context: Context): ChronosDatabase =
        Room.databaseBuilder(context, ChronosDatabase::class.java, "chronos.db")
            .fallbackToDestructiveMigration()
            .build()

    @Provides
    fun provideReminderDao(db: ChronosDatabase): ReminderDao = db.reminderDao()

    @Provides
    fun provideDfmDao(db: ChronosDatabase): DfmDao = db.dfmDao()
}
