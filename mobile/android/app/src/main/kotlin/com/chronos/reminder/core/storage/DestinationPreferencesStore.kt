package com.chronos.reminder.core.storage

import android.content.Context
import android.content.SharedPreferences
import dagger.hilt.android.qualifiers.ApplicationContext
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Persists the last two destination types the user picked when creating a reminder.
 * Discord Channel is excluded — it requires a specific guild/channel that may no longer exist.
 */
@Singleton
class DestinationPreferencesStore @Inject constructor(@ApplicationContext context: Context) {

    private val prefs: SharedPreferences =
        context.getSharedPreferences("chronos_destination_prefs", Context.MODE_PRIVATE)

    fun save(destinationTypes: List<String>) {
        val eligible = destinationTypes
            .filter { it != "discord_channel" && it != "webhook" }
            .distinct()
            .takeLast(2)
        prefs.edit().putString(KEY_LAST_DESTINATIONS, eligible.joinToString(",")).apply()
    }

    fun getLast(): List<String> {
        val raw = prefs.getString(KEY_LAST_DESTINATIONS, null) ?: return emptyList()
        return raw.split(",").filter { it.isNotBlank() }
    }

    private companion object {
        const val KEY_LAST_DESTINATIONS = "last_destinations"
    }
}
