package com.chronos.reminder.core.storage

import android.content.Context
import dagger.hilt.android.qualifiers.ApplicationContext
import java.util.UUID
import javax.inject.Inject
import javax.inject.Singleton

// Stable per-install device identifier for FCM token registration.
// Not sensitive, so plain SharedPreferences is fine.
@Singleton
class DeviceIdStore @Inject constructor(@ApplicationContext context: Context) {

    private val prefs = context.getSharedPreferences("chronos_device", Context.MODE_PRIVATE)

    fun getDeviceId(): String {
        prefs.getString(KEY_DEVICE_ID, null)?.let { return it }
        val id = UUID.randomUUID().toString()
        prefs.edit().putString(KEY_DEVICE_ID, id).apply()
        return id
    }

    private companion object {
        const val KEY_DEVICE_ID = "device_id"
    }
}
