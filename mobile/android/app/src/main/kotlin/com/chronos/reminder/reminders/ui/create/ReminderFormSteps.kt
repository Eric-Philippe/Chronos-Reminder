package com.chronos.reminder.reminders.ui.create

import androidx.compose.foundation.horizontalScroll
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.ui.text.input.KeyboardCapitalization
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.DatePicker
import androidx.compose.material3.DatePickerDialog
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilterChip
import androidx.compose.material3.FilterChipDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TimePicker
import androidx.compose.material3.rememberDatePickerState
import androidx.compose.material3.rememberTimePickerState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.components.ChronosButton
import com.chronos.reminder.core.ui.components.ChronosButtonStyle
import com.chronos.reminder.core.ui.components.ChronosCard
import com.chronos.reminder.core.ui.components.ChronosTextField
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.BackgroundCard
import com.chronos.reminder.core.ui.theme.BackgroundMuted
import com.chronos.reminder.core.ui.theme.ForegroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMuted
import com.chronos.reminder.reminders.domain.RecurrenceType
import java.time.Instant
import java.time.LocalDate
import java.time.LocalTime
import java.time.ZoneOffset
import java.time.format.DateTimeFormatter

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun WhenStep(
    form: ReminderForm,
    onDateSelected: (LocalDate) -> Unit,
    onTimeSelected: (LocalTime) -> Unit,
    onRecurrenceSelected: (RecurrenceType) -> Unit,
    modifier: Modifier = Modifier,
) {
    var showDatePicker by rememberSaveable { mutableStateOf(false) }
    var showTimePicker by rememberSaveable { mutableStateOf(false) }

    Column(modifier = modifier) {
        Text(stringResource(R.string.step_when_title), style = MaterialTheme.typography.headlineMedium)
        Spacer(Modifier.height(24.dp))

        Row(horizontalArrangement = Arrangement.spacedBy(12.dp)) {
            ChronosCard(modifier = Modifier.weight(1f), onClick = { showDatePicker = true }) {
                Column(Modifier.padding(16.dp)) {
                    Text(
                        stringResource(R.string.pick_date),
                        style = MaterialTheme.typography.labelSmall,
                        color = ForegroundMuted,
                    )
                    Spacer(Modifier.height(4.dp))
                    Text(
                        form.date?.format(DateTimeFormatter.ofPattern("MMM d, yyyy")) ?: "—",
                        style = MaterialTheme.typography.titleMedium,
                        color = if (form.date != null) AccentOrange else ForegroundMain,
                    )
                }
            }
            ChronosCard(modifier = Modifier.weight(1f), onClick = { showTimePicker = true }) {
                Column(Modifier.padding(16.dp)) {
                    Text(
                        stringResource(R.string.pick_time),
                        style = MaterialTheme.typography.labelSmall,
                        color = ForegroundMuted,
                    )
                    Spacer(Modifier.height(4.dp))
                    Text(
                        form.time?.format(DateTimeFormatter.ofPattern("HH:mm")) ?: "—",
                        style = MaterialTheme.typography.titleMedium,
                        color = if (form.time != null) AccentOrange else ForegroundMain,
                    )
                }
            }
        }

        Spacer(Modifier.height(24.dp))
        Text(
            stringResource(R.string.recurrence),
            style = MaterialTheme.typography.labelLarge,
            color = ForegroundMuted,
        )
        Spacer(Modifier.height(8.dp))
        RecurrenceChipRow(selected = form.recurrence, onSelect = onRecurrenceSelected)
    }

    if (showDatePicker) {
        val datePickerState = rememberDatePickerState(
            initialSelectedDateMillis = form.date
                ?.atStartOfDay(ZoneOffset.UTC)?.toInstant()?.toEpochMilli()
                ?: System.currentTimeMillis(),
        )
        DatePickerDialog(
            onDismissRequest = { showDatePicker = false },
            confirmButton = {
                TextButton(onClick = {
                    datePickerState.selectedDateMillis?.let { millis ->
                        onDateSelected(Instant.ofEpochMilli(millis).atZone(ZoneOffset.UTC).toLocalDate())
                    }
                    showDatePicker = false
                }) { Text(stringResource(R.string.ok), color = AccentOrange) }
            },
            dismissButton = {
                TextButton(onClick = { showDatePicker = false }) {
                    Text(stringResource(R.string.cancel), color = ForegroundMuted)
                }
            },
        ) {
            DatePicker(state = datePickerState)
        }
    }

    if (showTimePicker) {
        val timePickerState = rememberTimePickerState(
            initialHour = form.time?.hour ?: 12,
            initialMinute = form.time?.minute ?: 0,
            is24Hour = true,
        )
        AlertDialog(
            onDismissRequest = { showTimePicker = false },
            containerColor = BackgroundCard,
            title = { Text(stringResource(R.string.pick_time), style = MaterialTheme.typography.titleLarge) },
            text = { TimePicker(state = timePickerState) },
            confirmButton = {
                TextButton(onClick = {
                    onTimeSelected(LocalTime.of(timePickerState.hour, timePickerState.minute))
                    showTimePicker = false
                }) { Text(stringResource(R.string.ok), color = AccentOrange) }
            },
            dismissButton = {
                TextButton(onClick = { showTimePicker = false }) {
                    Text(stringResource(R.string.cancel), color = ForegroundMuted)
                }
            },
        )
    }
}

@Composable
fun RecurrenceChipRow(
    selected: RecurrenceType,
    onSelect: (RecurrenceType) -> Unit,
    modifier: Modifier = Modifier,
) {
    Row(
        modifier = modifier
            .fillMaxWidth()
            .horizontalScroll(rememberScrollState()),
        horizontalArrangement = Arrangement.spacedBy(8.dp),
    ) {
        RecurrenceType.entries.sortedBy { it.apiValue }.forEach { type ->
            FilterChip(
                selected = selected == type,
                onClick = { onSelect(type) },
                label = { Text(type.displayLabel, style = MaterialTheme.typography.labelLarge) },
                colors = FilterChipDefaults.filterChipColors(
                    containerColor = BackgroundMuted,
                    labelColor = ForegroundMuted,
                    selectedContainerColor = AccentOrange,
                    selectedLabelColor = ForegroundMain,
                ),
            )
        }
    }
}

@Composable
fun WhatStep(
    form: ReminderForm,
    onMessageChange: (String) -> Unit,
    modifier: Modifier = Modifier,
) {
    Column(modifier = modifier) {
        Text(stringResource(R.string.step_what_title), style = MaterialTheme.typography.headlineMedium)
        Spacer(Modifier.height(24.dp))
        ChronosTextField(
            value = form.message,
            onValueChange = { value ->
                val capitalized = if (value.isNotEmpty()) value[0].uppercaseChar() + value.drop(1) else value
                onMessageChange(capitalized)
            },
            modifier = Modifier.fillMaxWidth(),
            placeholder = stringResource(R.string.message_placeholder),
            singleLine = false,
            minLines = 4,
            keyboardOptions = KeyboardOptions(capitalization = KeyboardCapitalization.Sentences),
        )
        Spacer(Modifier.height(8.dp))
        Text(
            text = stringResource(R.string.char_count, form.message.length, MAX_MESSAGE_LENGTH),
            style = MaterialTheme.typography.labelSmall,
            color = ForegroundMuted,
        )
    }
}

@Composable
fun StepNavButtons(
    showBack: Boolean,
    nextLabel: String,
    onBack: () -> Unit,
    onNext: () -> Unit,
    modifier: Modifier = Modifier,
    nextEnabled: Boolean = true,
    nextLoading: Boolean = false,
) {
    Row(modifier = modifier.fillMaxWidth()) {
        if (showBack) {
            ChronosButton(
                text = stringResource(R.string.back),
                onClick = onBack,
                modifier = Modifier.weight(1f),
                style = ChronosButtonStyle.Secondary,
            )
            Spacer(Modifier.width(12.dp))
        }
        ChronosButton(
            text = nextLabel,
            onClick = onNext,
            modifier = Modifier.weight(1f),
            enabled = nextEnabled,
            loading = nextLoading,
        )
    }
}
