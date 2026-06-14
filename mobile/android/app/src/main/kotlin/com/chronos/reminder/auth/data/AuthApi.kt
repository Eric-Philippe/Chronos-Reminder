package com.chronos.reminder.auth.data

import com.chronos.reminder.reminders.data.MessageResponse
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST

interface AuthApi {

    @POST("api/auth/login")
    suspend fun login(@Body body: LoginRequest): Response<AuthResponseDto>

    @POST("api/auth/register")
    suspend fun register(@Body body: RegisterRequest): Response<AuthResponseDto>

    @POST("api/auth/logout")
    suspend fun logout(): Response<MessageResponse>

    @POST("api/auth/discord/callback")
    suspend fun discordCallback(@Body body: DiscordCallbackRequest): Response<DiscordCallbackResponseDto>

    @POST("api/auth/discord/setup")
    suspend fun discordSetup(@Body body: DiscordSetupRequest): Response<DiscordCallbackResponseDto>

    @POST("api/auth/password-reset/request")
    suspend fun requestPasswordReset(@Body body: PasswordResetRequest): Response<MessageResponse>

    @POST("api/auth/verify/resend")
    suspend fun resendVerification(@Body body: ResendVerificationRequest): Response<MessageResponse>
}
