package com.chronos.reminder.core.di

import com.chronos.reminder.reminders.data.RemindersRepository
import com.chronos.reminder.reminders.data.RemindersRepositoryImpl
import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class RepositoryModule {

    @Binds
    @Singleton
    abstract fun bindRemindersRepository(impl: RemindersRepositoryImpl): RemindersRepository
}
