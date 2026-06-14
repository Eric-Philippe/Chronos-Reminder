package com.chronos.reminder.auth.data

import android.os.Build
import android.util.Log
import com.chronos.reminder.account.data.AccountApi
import com.chronos.reminder.account.data.MobileIdentityRequest
import com.chronos.reminder.core.database.ChronosDatabase
import com.chronos.reminder.core.network.ApiResult
import com.chronos.reminder.core.network.safeApiCall
import com.chronos.reminder.core.storage.TokenStore
import com.chronos.reminder.notifications.FcmTokenManager
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AuthRepository @Inject constructor(
    private val api: AuthApi,
    private val accountApi: AccountApi,
    private val tokenStore: TokenStore,
    private val database: ChronosDatabase,
    private val fcmTokenManager: FcmTokenManager,
) {

    // Called right after a token is stored, on any successful login path.
    // Registers push + records that this account has connected from mobile.
    private suspend fun onAuthenticated() {
        fcmTokenManager.registerCurrentToken()
        val deviceName = "${Build.MANUFACTURER} ${Build.MODEL}".trim()
        val result = safeApiCall { accountApi.ensureMobileIdentity(MobileIdentityRequest(deviceName)) }
        Log.d("AuthRepository", "Mobile identity registration: $result")
    }

    fun isLoggedIn(): Boolean = tokenStore.getToken() != null

    suspend fun login(email: String, password: String): ApiResult<Unit> {
        val result = safeApiCall { api.login(LoginRequest(email, password)) }
        return when (result) {
            is ApiResult.Success -> {
                val token = result.data.token
                if (token.isNullOrBlank()) {
                    ApiResult.Error(-1, result.data.message ?: "No token in response")
                } else {
                    tokenStore.saveToken(token)
                    onAuthenticated()
                    ApiResult.Success(Unit)
                }
            }
            is ApiResult.Error -> result
            is ApiResult.NetworkError -> result
        }
    }

    suspend fun register(
        email: String,
        username: String,
        password: String,
        timezone: String,
    ): ApiResult<Unit> =
        safeApiCall { api.register(RegisterRequest(email, username, password, timezone)) }.map { }

    suspend fun loginWithDiscordCode(code: String): ApiResult<Unit> {
        val result = safeApiCall { api.discordCallback(DiscordCallbackRequest(code)) }
        return when (result) {
            is ApiResult.Success -> {
                val data = result.data
                when {
                    !data.token.isNullOrBlank() -> {
                        tokenStore.saveToken(data.token)
                        onAuthenticated()
                        ApiResult.Success(Unit)
                    }
                    data.needsSetup -> ApiResult.Error(NEEDS_SETUP_CODE, data.message ?: "Account setup required")
                    else -> ApiResult.Error(-1, data.message ?: "Discord login failed")
                }
            }
            is ApiResult.Error -> result
            is ApiResult.NetworkError -> result
        }
    }

    suspend fun requestPasswordReset(email: String): ApiResult<Unit> =
        safeApiCall { api.requestPasswordReset(PasswordResetRequest(email)) }.map { }

    suspend fun resendVerification(email: String): ApiResult<Unit> =
        safeApiCall { api.resendVerification(ResendVerificationRequest(email)) }.map { }

    suspend fun logout() {
        fcmTokenManager.unregisterCurrentToken()
        safeApiCall { api.logout() } // best effort; local state is cleared regardless
        clearLocalData()
    }

    suspend fun clearLocalData() {
        tokenStore.clearToken()
        withContext(Dispatchers.IO) { database.clearAllTables() }
    }

    companion object {
        // Discord account exists but has no Chronos account; setup happens on the web app
        const val NEEDS_SETUP_CODE = -2
    }
}
