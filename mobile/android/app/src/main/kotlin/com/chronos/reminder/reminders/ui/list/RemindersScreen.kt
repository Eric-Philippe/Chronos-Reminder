package com.chronos.reminder.reminders.ui.list

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.MoreVert
import androidx.compose.material.icons.filled.Schedule
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilterChip
import androidx.compose.material3.FilterChipDefaults
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.material3.Text
import androidx.compose.material3.pulltorefresh.PullToRefreshBox
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.components.ChronosCard
import com.chronos.reminder.core.ui.components.ChronosTopBar
import com.chronos.reminder.core.ui.components.ConfirmDeleteDialog
import com.chronos.reminder.core.ui.components.DestinationChip
import com.chronos.reminder.core.ui.components.EmptyState
import com.chronos.reminder.core.ui.components.ErrorBanner
import com.chronos.reminder.core.ui.components.RecurrenceBadge
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.BackgroundCard
import com.chronos.reminder.core.ui.theme.BackgroundMain
import com.chronos.reminder.core.ui.theme.BackgroundMuted
import com.chronos.reminder.core.ui.theme.ForegroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMuted
import com.chronos.reminder.core.util.formatNextFire
import com.chronos.reminder.reminders.domain.Reminder

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun RemindersScreen(
    onCreateReminder: () -> Unit,
    onOpenReminder: (String) -> Unit,
    showCreatedBanner: Boolean = false,
    onCreatedBannerConsumed: () -> Unit = {},
    viewModel: RemindersViewModel = hiltViewModel(),
) {
    val state by viewModel.state.collectAsStateWithLifecycle()
    val reminders by viewModel.reminders.collectAsStateWithLifecycle()
    var deleteCandidate by rememberSaveable { mutableStateOf<String?>(null) }
    val snackbarHostState = remember { SnackbarHostState() }
    val reminderCreatedText = stringResource(R.string.reminder_created)

    LaunchedEffect(showCreatedBanner) {
        if (showCreatedBanner) {
            snackbarHostState.showSnackbar(reminderCreatedText)
            onCreatedBannerConsumed()
        }
    }

    Scaffold(
        containerColor = BackgroundMain,
        topBar = { ChronosTopBar(title = stringResource(R.string.reminders_title)) },
        snackbarHost = { SnackbarHost(snackbarHostState) },
        floatingActionButton = {
            FloatingActionButton(
                onClick = onCreateReminder,
                containerColor = AccentOrange,
                contentColor = ForegroundMain,
                modifier = Modifier.size(56.dp),
            ) {
                Icon(Icons.Default.Add, contentDescription = stringResource(R.string.create_reminder))
            }
        },
    ) { padding ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding),
        ) {
            Column(modifier = Modifier.fillMaxSize()) {
                FilterChipsRow(
                    selected = state.filter,
                    onSelect = viewModel::setFilter,
                    modifier = Modifier.padding(horizontal = 16.dp),
                )
                Spacer(Modifier.height(8.dp))

                PullToRefreshBox(
                    isRefreshing = state.refreshing,
                    onRefresh = viewModel::refresh,
                    modifier = Modifier.fillMaxSize(),
                ) {
                    if (reminders.isEmpty() && !state.refreshing) {
                        EmptyState(
                            icon = Icons.Default.Schedule,
                            title = stringResource(R.string.empty_reminders_title),
                            subtitle = stringResource(R.string.empty_reminders_subtitle),
                            ctaLabel = stringResource(R.string.empty_reminders_cta),
                            onCta = onCreateReminder,
                        )
                    } else {
                        LazyColumn(
                            modifier = Modifier.fillMaxSize(),
                            verticalArrangement = Arrangement.spacedBy(8.dp),
                            contentPadding = androidx.compose.foundation.layout.PaddingValues(
                                start = 16.dp,
                                end = 16.dp,
                                bottom = 96.dp,
                            ),
                        ) {
                            items(reminders, key = { it.id }) { reminder ->
                                ReminderCard(
                                    reminder = reminder,
                                    userTimezone = state.userTimezone,
                                    onClick = { onOpenReminder(reminder.id) },
                                    onEdit = { onOpenReminder(reminder.id) },
                                    onDuplicate = { viewModel.duplicate(reminder.id) },
                                    onTogglePause = {
                                        if (reminder.isPaused) {
                                            viewModel.resume(reminder.id)
                                        } else {
                                            viewModel.pause(reminder.id)
                                        }
                                    },
                                    onDelete = { deleteCandidate = reminder.id },
                                )
                            }
                        }
                    }
                }
            }

            ErrorBanner(
                message = state.error,
                modifier = Modifier.align(Alignment.BottomCenter),
                onDismiss = viewModel::clearError,
            )
        }
    }

    deleteCandidate?.let { id ->
        ConfirmDeleteDialog(
            title = stringResource(R.string.delete_reminder_title),
            text = stringResource(R.string.delete_reminder_text),
            onConfirm = {
                viewModel.delete(id)
                deleteCandidate = null
            },
            onDismiss = { deleteCandidate = null },
        )
    }
}

@Composable
private fun FilterChipsRow(
    selected: ReminderFilter,
    onSelect: (ReminderFilter) -> Unit,
    modifier: Modifier = Modifier,
) {
    val labels = mapOf(
        ReminderFilter.ALL to stringResource(R.string.filter_all),
        ReminderFilter.ACTIVE to stringResource(R.string.filter_active),
        ReminderFilter.PAUSED to stringResource(R.string.filter_paused),
    )
    Row(modifier = modifier, horizontalArrangement = Arrangement.spacedBy(8.dp)) {
        ReminderFilter.entries.forEach { filter ->
            FilterChip(
                selected = selected == filter,
                onClick = { onSelect(filter) },
                label = { Text(labels.getValue(filter), style = MaterialTheme.typography.labelLarge) },
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
fun ReminderCard(
    reminder: Reminder,
    userTimezone: String,
    onClick: () -> Unit,
    onEdit: () -> Unit,
    onDuplicate: () -> Unit,
    onTogglePause: () -> Unit,
    onDelete: () -> Unit,
) {
    var menuOpen by rememberSaveable { mutableStateOf(false) }

    ChronosCard(
        modifier = Modifier
            .fillMaxWidth()
            .alpha(if (reminder.isPaused) 0.6f else 1f),
        onClick = onClick,
    ) {
        Column(Modifier.padding(16.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                RecurrenceBadge(reminder.recurrence)
                if (reminder.isPaused) {
                    Spacer(Modifier.width(6.dp))
                    Text(
                        stringResource(R.string.paused_badge),
                        style = MaterialTheme.typography.labelSmall,
                        color = ForegroundMuted,
                    )
                }
                Spacer(Modifier.weight(1f))
                Text(
                    text = formatNextFire(reminder.nextFireUtc ?: reminder.remindAtUtc, userTimezone),
                    style = MaterialTheme.typography.labelSmall,
                    color = AccentOrange,
                )
                Box {
                    IconButton(onClick = { menuOpen = true }) {
                        Icon(
                            Icons.Default.MoreVert,
                            contentDescription = stringResource(R.string.more_options),
                            tint = ForegroundMuted,
                        )
                    }
                    DropdownMenu(
                        expanded = menuOpen,
                        onDismissRequest = { menuOpen = false },
                        containerColor = BackgroundCard,
                    ) {
                        DropdownMenuItem(
                            text = { Text(stringResource(R.string.action_edit)) },
                            onClick = {
                                menuOpen = false
                                onEdit()
                            },
                        )
                        DropdownMenuItem(
                            text = { Text(stringResource(R.string.action_duplicate)) },
                            onClick = {
                                menuOpen = false
                                onDuplicate()
                            },
                        )
                        DropdownMenuItem(
                            text = {
                                Text(
                                    stringResource(
                                        if (reminder.isPaused) R.string.action_resume else R.string.action_pause,
                                    ),
                                )
                            },
                            onClick = {
                                menuOpen = false
                                onTogglePause()
                            },
                        )
                        DropdownMenuItem(
                            text = {
                                Text(
                                    stringResource(R.string.action_delete),
                                    color = MaterialTheme.colorScheme.error,
                                )
                            },
                            onClick = {
                                menuOpen = false
                                onDelete()
                            },
                        )
                    }
                }
            }

            Text(
                text = reminder.message,
                style = MaterialTheme.typography.bodyLarge,
                maxLines = 2,
                overflow = TextOverflow.Ellipsis,
            )
            Spacer(Modifier.height(8.dp))
            LazyRow(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                items(reminder.destinations) { destination ->
                    DestinationChip(type = destination.type)
                }
            }
        }
    }
}
