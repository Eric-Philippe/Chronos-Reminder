package com.chronos.reminder.notifications

import android.util.Log
import com.chronos.reminder.core.network.safeApiCall
import com.chronos.reminder.core.storage.DeviceIdStore
import com.chronos.reminder.core.storage.TokenStore
import com.google.firebase.messaging.FirebaseMessaging
import kotlinx.coroutines.tasks.await
import javax.inject.Inject
import javax.inject.Singleton

// Registers/unregisters the FCM token with the backend. Failures are logged and
// swallowed: push registration must never block login or logout.
@Singleton
class FcmTokenManager @Inject constructor(
    private val fcmApi: FcmApi,
    private val deviceIdStore: DeviceIdStore,
    private val tokenStore: TokenStore,
) {

    suspend fun registerCurrentToken() {
        val fcmToken = currentFcmToken() ?: return
        registerToken(fcmToken)
    }

    suspend fun registerToken(fcmToken: String) {
        if (tokenStore.getToken() == null) return
        val result = safeApiCall {
            fcmApi.registerToken(FcmTokenRequest(fcmToken, deviceIdStore.getDeviceId()))
        }
        Log.d(TAG, "FCM token registration: $result")
    }

    suspend fun unregisterCurrentToken() {
        val fcmToken = currentFcmToken() ?: return
        val result = safeApiCall {
            fcmApi.unregisterToken(FcmTokenRequest(fcmToken, deviceIdStore.getDeviceId()))
        }
        Log.d(TAG, "FCM token unregistration: $result")
    }

    private suspend fun currentFcmToken(): String? = try {
        FirebaseMessaging.getInstance().token.await()
    } catch (e: Exception) {
        // Firebase is unavailable when google-services.json is missing (e.g. local dev)
        Log.w(TAG, "Could not obtain FCM token", e)
        null
    }

    private companion object {
        const val TAG = "FcmTokenManager"
    }
}
