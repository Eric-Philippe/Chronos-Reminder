package com.chronos.reminder.reminders.domain

enum class RecurrenceType(val apiValue: Int, val apiString: String, val displayLabel: String) {
    ONCE(0, "ONCE", "Once"),
    YEARLY(1, "YEARLY", "Yearly"),
    MONTHLY(2, "MONTHLY", "Monthly"),
    WEEKLY(3, "WEEKLY", "Weekly"),
    DAILY(4, "DAILY", "Daily"),
    HOURLY(5, "HOURLY", "Hourly"),
    WORKDAYS(6, "WORKDAYS", "Workdays"),
    WEEKEND(7, "WEEKEND", "Weekend"),
    ;

    companion object {
        const val PAUSED_FLAG = 128

        fun isPaused(state: Int): Boolean = (state and PAUSED_FLAG) != 0

        fun fromState(state: Int): RecurrenceType =
            entries.firstOrNull { it.apiValue == (state and 0x7F) } ?: ONCE

        fun fromApiString(value: String?): RecurrenceType =
            entries.firstOrNull { it.apiString.equals(value, ignoreCase = true) } ?: ONCE
    }
}
