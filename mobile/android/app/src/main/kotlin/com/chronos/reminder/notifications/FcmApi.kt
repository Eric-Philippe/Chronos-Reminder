package com.chronos.reminder.notifications

import com.chronos.reminder.reminders.data.MessageResponse
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.HTTP
import retrofit2.http.POST

interface FcmApi {

    @POST("api/fcm/token")
    suspend fun registerToken(@Body body: FcmTokenRequest): Response<MessageResponse>

    @HTTP(method = "DELETE", path = "api/fcm/token", hasBody = true)
    suspend fun unregisterToken(@Body body: FcmTokenRequest): Response<MessageResponse>
}
