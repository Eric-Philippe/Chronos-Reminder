package com.chronos.reminder.dfm.data

import com.chronos.reminder.reminders.data.MessageResponse
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Path

interface DfmApi {

    @GET("api/dfm")
    suspend fun getNote(): Response<DfmNoteDto>

    @POST("api/dfm/items")
    suspend fun addItem(@Body body: AddDfmItemRequest): Response<DfmItemDto>

    @PUT("api/dfm/items/{id}")
    suspend fun updateItem(@Path("id") id: String, @Body body: UpdateDfmItemRequest): Response<DfmItemDto>

    @DELETE("api/dfm/items/{id}")
    suspend fun deleteItem(@Path("id") id: String): Response<MessageResponse>

    @PUT("api/dfm/reminder")
    suspend fun setReminder(@Body body: DfmReminderRequest): Response<DfmNoteDto>

    @DELETE("api/dfm/reminder")
    suspend fun removeReminder(): Response<DfmNoteDto>

    @POST("api/dfm/send")
    suspend fun sendNow(): Response<MessageResponse>
}
