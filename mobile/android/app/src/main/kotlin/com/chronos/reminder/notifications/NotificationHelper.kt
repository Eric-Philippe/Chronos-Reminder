package com.chronos.reminder.notifications

import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Context
import android.content.Intent
import android.graphics.BitmapFactory
import androidx.core.app.NotificationCompat
import com.chronos.reminder.ChronosApp
import com.chronos.reminder.MainActivity
import com.chronos.reminder.R
import dagger.hilt.android.qualifiers.ApplicationContext
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class NotificationHelper @Inject constructor(
    @ApplicationContext private val context: Context,
) {

    fun showReminderNotification(message: String, reminderId: String?) {
        val notifId = reminderId?.hashCode() ?: System.currentTimeMillis().toInt()

        val tapIntent = Intent(context, MainActivity::class.java).apply {
            flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TOP
            reminderId?.let { putExtra(MainActivity.EXTRA_REMINDER_ID, it) }
        }
        val tapPending = PendingIntent.getActivity(
            context,
            notifId,
            tapIntent,
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE,
        )

        fun snoozePending(minutes: Int): PendingIntent {
            val intent = Intent(context, SnoozeReceiver::class.java).apply {
                action = SnoozeReceiver.ACTION_SNOOZE
                putExtra(SnoozeReceiver.EXTRA_REMINDER_ID, reminderId)
                putExtra(SnoozeReceiver.EXTRA_MINUTES, minutes)
                putExtra(SnoozeReceiver.EXTRA_NOTIFICATION_ID, notifId)
            }
            return PendingIntent.getBroadcast(
                context,
                notifId xor minutes,
                intent,
                PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE,
            )
        }

        val notification = NotificationCompat.Builder(context, ChronosApp.REMINDERS_CHANNEL_ID)
            .setSmallIcon(R.drawable.ic_notification)
            .setLargeIcon(BitmapFactory.decodeResource(context.resources, R.mipmap.ic_launcher))
            .setContentTitle(context.getString(R.string.notification_title))
            .setContentText(message)
            .setStyle(NotificationCompat.BigTextStyle().bigText(message))
            .setContentIntent(tapPending)
            .setAutoCancel(true)
            .setPriority(NotificationCompat.PRIORITY_HIGH)
            .apply {
                if (reminderId != null) {
                    addAction(R.drawable.ic_notification, context.getString(R.string.notification_snooze_10m), snoozePending(10))
                    addAction(R.drawable.ic_notification, context.getString(R.string.notification_snooze_1h), snoozePending(60))
                    addAction(R.drawable.ic_notification, context.getString(R.string.notification_snooze_1d), snoozePending(1440))
                }
            }
            .build()

        val manager = context.getSystemService(NotificationManager::class.java)
        manager.notify(notifId, notification)
    }

    fun showSnoozeConfirmation(minutes: Int, originalNotifId: Int) {
        val text = when (minutes) {
            10 -> context.getString(R.string.notification_snoozed_10m)
            60 -> context.getString(R.string.notification_snoozed_1h)
            1440 -> context.getString(R.string.notification_snoozed_1d)
            else -> context.getString(R.string.notification_snoozed_10m)
        }

        val notification = NotificationCompat.Builder(context, ChronosApp.REMINDERS_CHANNEL_ID)
            .setSmallIcon(R.drawable.ic_notification)
            .setContentTitle(context.getString(R.string.notification_snoozed_title))
            .setContentText(text)
            .setAutoCancel(true)
            .setPriority(NotificationCompat.PRIORITY_LOW)
            .build()

        val manager = context.getSystemService(NotificationManager::class.java)
        manager.notify(originalNotifId xor 0x5A5A, notification)
    }
}
