package com.chronos.reminder.core.util

import java.time.Duration
import java.time.Instant
import java.time.ZoneId
import java.time.ZonedDateTime
import java.time.format.DateTimeFormatter

private val weekFormatter = DateTimeFormatter.ofPattern("EEE 'at' HH:mm")
private val dateFormatter = DateTimeFormatter.ofPattern("MMM d 'at' HH:mm")
private val fullFormatter = DateTimeFormatter.ofPattern("MMM d, yyyy 'at' HH:mm")
private val dayFormatter = DateTimeFormatter.ofPattern("MMM d, yyyy")

private fun zoneOrDefault(userTz: String): ZoneId =
    runCatching { ZoneId.of(userTz) }.getOrDefault(ZoneId.systemDefault())

// "in 3h" if <24h away, "Wed at 14:30" if <7 days, "Nov 15 at 14:30" otherwise
fun formatNextFire(instant: Instant, userTz: String): String {
    val zone = zoneOrDefault(userTz)
    val zdt = ZonedDateTime.ofInstant(instant, zone)
    val now = ZonedDateTime.now(zone)
    val untilFire = Duration.between(now, zdt)
    return when {
        untilFire.toHours() in 0 until 1 -> "in ${untilFire.toMinutes().coerceAtLeast(0)}m"
        untilFire.toHours() in 1 until 24 -> "in ${untilFire.toHours()}h"
        untilFire.toDays() in 0 until 7 -> zdt.format(weekFormatter)
        else -> zdt.format(dateFormatter)
    }
}

fun formatFullDateTime(instant: Instant, userTz: String): String =
    ZonedDateTime.ofInstant(instant, zoneOrDefault(userTz)).format(fullFormatter)

fun formatDate(instant: Instant, userTz: String): String =
    ZonedDateTime.ofInstant(instant, zoneOrDefault(userTz)).format(dayFormatter)
