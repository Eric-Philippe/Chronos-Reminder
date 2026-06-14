package com.chronos.reminder.reminders.ui.detail

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ExperimentalLayoutApi
import androidx.compose.foundation.layout.FlowRow
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.components.ChronosButton
import com.chronos.reminder.core.ui.components.ChronosButtonStyle
import com.chronos.reminder.core.ui.components.ChronosCard
import com.chronos.reminder.core.ui.components.ChronosTopBar
import com.chronos.reminder.core.ui.components.ConfirmDeleteDialog
import com.chronos.reminder.core.ui.components.DestinationChip
import com.chronos.reminder.core.ui.components.ErrorBanner
import com.chronos.reminder.core.ui.components.RecurrenceBadge
import com.chronos.reminder.core.ui.theme.BackgroundMain
import com.chronos.reminder.core.ui.theme.DestructiveRed
import com.chronos.reminder.core.ui.theme.ForegroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMuted
import com.chronos.reminder.core.util.formatDate
import com.chronos.reminder.core.util.formatFullDateTime
import com.chronos.reminder.reminders.ui.create.WhatStep
import com.chronos.reminder.reminders.ui.create.WhenStep

@OptIn(ExperimentalLayoutApi::class)
@Composable
fun ReminderDetailScreen(
    onBack: () -> Unit,
    viewModel: ReminderDetailViewModel = hiltViewModel(),
) {
    val state by viewModel.state.collectAsStateWithLifecycle()
    val reminder by viewModel.reminder.collectAsStateWithLifecycle()
    var showDeleteDialog by rememberSaveable { mutableStateOf(false) }

    LaunchedEffect(state.deleted) {
        if (state.deleted) onBack()
    }

    Scaffold(
        containerColor = BackgroundMain,
        topBar = {
            ChronosTopBar(
                title = stringResource(R.string.reminder_detail_title),
                onBack = onBack,
                actions = {
                    if (!state.editing) {
                        IconButton(onClick = viewModel::startEditing) {
                            Icon(
                                Icons.Default.Edit,
                                contentDescription = stringResource(R.string.edit_reminder),
                                tint = ForegroundMain,
                            )
                        }
                        IconButton(onClick = { showDeleteDialog = true }) {
                            Icon(
                                Icons.Default.Delete,
                                contentDescription = stringResource(R.string.delete),
                                tint = DestructiveRed,
                            )
                        }
                    }
                },
            )
        },
    ) { padding ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding),
        ) {
            val current = reminder
            if (current != null) {
                Column(
                    modifier = Modifier
                        .fillMaxSize()
                        .verticalScroll(rememberScrollState())
                        .padding(horizontal = 16.dp),
                ) {
                    Spacer(Modifier.height(8.dp))
                    if (state.editing) {
                        WhenStep(
                            form = state.form,
                            onDateSelected = { viewModel.updateForm(state.form.copy(date = it)) },
                            onTimeSelected = { viewModel.updateForm(state.form.copy(time = it)) },
                            onRecurrenceSelected = { viewModel.updateForm(state.form.copy(recurrence = it)) },
                        )
                        Spacer(Modifier.height(24.dp))
                        WhatStep(
                            form = state.form,
                            onMessageChange = { viewModel.updateForm(state.form.copy(message = it)) },
                        )
                        Spacer(Modifier.height(32.dp))
                        Row(Modifier.fillMaxWidth()) {
                            ChronosButton(
                                text = stringResource(R.string.cancel),
                                onClick = viewModel::cancelEditing,
                                modifier = Modifier.weight(1f),
                                style = ChronosButtonStyle.Secondary,
                            )
                            Spacer(Modifier.width(12.dp))
                            ChronosButton(
                                text = stringResource(R.string.save),
                                onClick = viewModel::save,
                                modifier = Modifier.weight(1f),
                                loading = state.saving,
                                enabled = state.form.message.isNotBlank(),
                            )
                        }
                    } else {
                        ChronosCard(modifier = Modifier.fillMaxWidth()) {
                            Column(Modifier.padding(16.dp)) {
                                Row(verticalAlignment = Alignment.CenterVertically) {
                                    RecurrenceBadge(current.recurrence)
                                    if (current.isPaused) {
                                        Spacer(Modifier.width(8.dp))
                                        Text(
                                            stringResource(R.string.paused_badge),
                                            style = MaterialTheme.typography.labelSmall,
                                            color = ForegroundMuted,
                                        )
                                    }
                                }
                                Spacer(Modifier.height(12.dp))
                                Text(current.message, style = MaterialTheme.typography.bodyLarge)
                            }
                        }
                        Spacer(Modifier.height(16.dp))

                        current.nextFireUtc?.let { nextFire ->
                            Text(
                                text = stringResource(
                                    R.string.next_fire_label,
                                    formatFullDateTime(nextFire, state.userTimezone),
                                ),
                                style = MaterialTheme.typography.bodyMedium,
                            )
                            Spacer(Modifier.height(4.dp))
                        }
                        current.createdAt?.let { createdAt ->
                            Text(
                                text = stringResource(
                                    R.string.created_label,
                                    formatDate(createdAt, state.userTimezone),
                                ),
                                style = MaterialTheme.typography.bodyMedium,
                                color = ForegroundMuted,
                            )
                        }
                        Spacer(Modifier.height(16.dp))

                        Text(
                            stringResource(R.string.destinations_label),
                            style = MaterialTheme.typography.labelLarge,
                            color = ForegroundMuted,
                        )
                        Spacer(Modifier.height(8.dp))
                        FlowRow(
                            horizontalArrangement = Arrangement.spacedBy(8.dp),
                            verticalArrangement = Arrangement.spacedBy(8.dp),
                        ) {
                            current.destinations.forEach { destination ->
                                DestinationChip(type = destination.type)
                            }
                        }
                        Spacer(Modifier.height(32.dp))

                        ChronosButton(
                            text = stringResource(
                                if (current.isPaused) R.string.action_resume else R.string.action_pause,
                            ),
                            onClick = viewModel::togglePause,
                            modifier = Modifier.fillMaxWidth(),
                            style = ChronosButtonStyle.Secondary,
                        )
                    }
                    Spacer(Modifier.height(24.dp))
                }
            }

            ErrorBanner(
                message = state.error,
                modifier = Modifier.align(Alignment.BottomCenter),
                onDismiss = viewModel::clearError,
            )
        }
    }

    if (showDeleteDialog) {
        ConfirmDeleteDialog(
            title = stringResource(R.string.delete_reminder_title),
            text = stringResource(R.string.delete_reminder_text),
            onConfirm = {
                showDeleteDialog = false
                viewModel.delete()
            },
            onDismiss = { showDeleteDialog = false },
        )
    }
}
