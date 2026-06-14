package com.chronos.reminder.notifications

import android.app.NotificationManager
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import com.chronos.reminder.reminders.data.RemindersApi
import com.chronos.reminder.reminders.data.SnoozeReminderRequest
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class SnoozeReceiver : BroadcastReceiver() {

    @Inject
    lateinit var remindersApi: RemindersApi

    @Inject
    lateinit var notificationHelper: NotificationHelper

    override fun onReceive(context: Context, intent: Intent) {
        if (intent.action != ACTION_SNOOZE) return

        val reminderId = intent.getStringExtra(EXTRA_REMINDER_ID) ?: return
        val minutes = intent.getIntExtra(EXTRA_MINUTES, 10)
        val notifId = intent.getIntExtra(EXTRA_NOTIFICATION_ID, 0)

        // Dismiss the original notification immediately
        context.getSystemService(NotificationManager::class.java)?.cancel(notifId)

        val result = goAsync()
        CoroutineScope(Dispatchers.IO).launch {
            try {
                remindersApi.snoozeReminder(reminderId, SnoozeReminderRequest(minutes))
                notificationHelper.showSnoozeConfirmation(minutes, notifId)
            } catch (_: Exception) {
                // silent — snooze is best-effort from the lock screen
            } finally {
                result.finish()
            }
        }
    }

    companion object {
        const val ACTION_SNOOZE = "com.chronos.reminder.ACTION_SNOOZE"
        const val EXTRA_REMINDER_ID = "reminder_id"
        const val EXTRA_MINUTES = "snooze_minutes"
        const val EXTRA_NOTIFICATION_ID = "notification_id"
    }
}
