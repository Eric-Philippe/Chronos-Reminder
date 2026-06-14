package com.chronos.reminder.auth.data

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class LoginRequest(
    val email: String,
    val password: String,
    @SerialName("remember_me") val rememberMe: Boolean = true,
)

@Serializable
data class AuthResponseDto(
    val id: String? = null,
    val email: String? = null,
    val username: String? = null,
    val token: String? = null,
    @SerialName("expires_at") val expiresAt: String? = null,
    val message: String? = null,
)

@Serializable
data class RegisterRequest(
    val email: String,
    val username: String,
    val password: String,
    val timezone: String,
)

@Serializable
data class DiscordCallbackRequest(
    val code: String,
    val state: String = "",
)

// Outcome of a Discord OAuth callback on mobile.
sealed interface DiscordLoginResult {
    // Existing account with credentials — the user is now logged in.
    data object LoggedIn : DiscordLoginResult

    // Brand-new Discord-only account that needs an email/password to finish.
    data class NeedsSetup(
        val accountId: String,
        val email: String?,
        val username: String?,
    ) : DiscordLoginResult
}

@Serializable
data class DiscordSetupRequest(
    @SerialName("account_id") val accountId: String,
    val email: String,
    val username: String,
    val password: String,
    val timezone: String,
)

@Serializable
data class DiscordCallbackResponseDto(
    val id: String? = null,
    val email: String? = null,
    val username: String? = null,
    val token: String? = null,
    @SerialName("expires_at") val expiresAt: String? = null,
    val message: String? = null,
    // Present when the Discord account has no Chronos account yet
    @SerialName("needs_setup") val needsSetup: Boolean = false,
    @SerialName("account_id") val accountId: String? = null,
    @SerialName("discord_email") val discordEmail: String? = null,
    @SerialName("discord_username") val discordUsername: String? = null,
)

@Serializable
data class PasswordResetRequest(
    val email: String,
)

@Serializable
data class ResendVerificationRequest(
    val email: String,
)
