package com.chronos.reminder.reminders.data

import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Path

interface RemindersApi {

    @GET("api/reminders")
    suspend fun getReminders(): Response<RemindersListResponse>

    @GET("api/reminders/{id}")
    suspend fun getReminder(@Path("id") id: String): Response<ReminderDto>

    @POST("api/reminders")
    suspend fun createReminder(@Body body: CreateReminderRequest): Response<ReminderDto>

    @PUT("api/reminders/{id}")
    suspend fun updateReminder(@Path("id") id: String, @Body body: CreateReminderRequest): Response<ReminderDto>

    @DELETE("api/reminders/{id}")
    suspend fun deleteReminder(@Path("id") id: String): Response<MessageResponse>

    @POST("api/reminders/{id}/pause")
    suspend fun pauseReminder(@Path("id") id: String): Response<ReminderDto>

    @POST("api/reminders/{id}/resume")
    suspend fun resumeReminder(@Path("id") id: String): Response<ReminderDto>

    @POST("api/reminders/{id}/duplicate")
    suspend fun duplicateReminder(@Path("id") id: String): Response<ReminderDto>

    @POST("api/reminders/{id}/snooze")
    suspend fun snoozeReminder(@Path("id") id: String, @Body body: SnoozeReminderRequest): Response<MessageResponse>
}
