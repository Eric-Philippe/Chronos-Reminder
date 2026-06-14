package com.chronos.reminder.core.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.ForegroundMuted
import com.chronos.reminder.core.ui.theme.RecurrenceHourly
import com.chronos.reminder.core.ui.theme.RecurrenceMonthly
import com.chronos.reminder.core.ui.theme.RecurrenceWeekend
import com.chronos.reminder.core.ui.theme.RecurrenceWeekly
import com.chronos.reminder.core.ui.theme.RecurrenceWorkdays
import com.chronos.reminder.core.ui.theme.RecurrenceYearly
import com.chronos.reminder.reminders.domain.RecurrenceType

fun RecurrenceType.badgeColor(): Color = when (this) {
    RecurrenceType.ONCE -> ForegroundMuted
    RecurrenceType.HOURLY -> RecurrenceHourly
    RecurrenceType.DAILY -> AccentOrange
    RecurrenceType.WORKDAYS -> RecurrenceWorkdays
    RecurrenceType.WEEKEND -> RecurrenceWeekend
    RecurrenceType.WEEKLY -> RecurrenceWeekly
    RecurrenceType.MONTHLY -> RecurrenceMonthly
    RecurrenceType.YEARLY -> RecurrenceYearly
}

@Composable
fun RecurrenceBadge(recurrence: RecurrenceType, modifier: Modifier = Modifier) {
    Text(
        text = recurrence.apiString,
        style = MaterialTheme.typography.labelSmall,
        color = Color.White,
        modifier = modifier
            .background(recurrence.badgeColor().copy(alpha = 0.85f), RoundedCornerShape(4.dp))
            .padding(horizontal = 6.dp, vertical = 2.dp),
    )
}
