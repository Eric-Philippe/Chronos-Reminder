package com.chronos.reminder.reminders.ui.create

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.layout.FlowRow
import androidx.compose.foundation.layout.ExperimentalLayoutApi
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AlternateEmail
import androidx.compose.material.icons.filled.Chat
import androidx.compose.material.icons.filled.Notifications
import androidx.compose.material.icons.filled.Tag
import androidx.compose.material.icons.filled.Webhook
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilterChip
import androidx.compose.material3.FilterChipDefaults
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.components.ChronosButton
import com.chronos.reminder.core.ui.components.ChronosCard
import com.chronos.reminder.core.ui.components.ChronosTextField
import com.chronos.reminder.core.ui.components.DestinationChip
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.BackgroundCard
import com.chronos.reminder.core.ui.theme.BackgroundMuted
import com.chronos.reminder.core.ui.theme.ForegroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMuted
import com.chronos.reminder.reminders.data.ChannelDto
import com.chronos.reminder.reminders.data.GuildDto

private const val DISCORD_TEXT_CHANNEL = 0

@OptIn(ExperimentalLayoutApi::class, ExperimentalMaterial3Api::class)
@Composable
fun WhereStep(
    uiState: CreateReminderUiState,
    onAddDiscordDm: () -> Unit,
    onAddEmail: (String) -> Unit,
    onAddPush: () -> Unit,
    onAddChannel: (GuildDto, ChannelDto) -> Unit,
    onAddWebhook: (url: String, platform: String) -> Unit,
    onRemove: (FormDestination) -> Unit,
    onLoadGuilds: () -> Unit,
    onLoadChannels: (guildId: String) -> Unit,
    modifier: Modifier = Modifier,
) {
    var showEmailInput by rememberSaveable { mutableStateOf(false) }
    var showWebhookInput by rememberSaveable { mutableStateOf(false) }
    var showChannelSheet by rememberSaveable { mutableStateOf(false) }

    Column(modifier = modifier) {
        Text(stringResource(R.string.step_where_title), style = MaterialTheme.typography.headlineMedium)
        Spacer(Modifier.height(16.dp))

        if (uiState.form.destinations.isNotEmpty()) {
            FlowRow(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(8.dp),
                verticalArrangement = Arrangement.spacedBy(8.dp),
            ) {
                uiState.form.destinations.forEach { destination ->
                    DestinationChip(
                        type = destination.type,
                        detail = destination.detail,
                        onRemove = { onRemove(destination) },
                    )
                }
            }
            Spacer(Modifier.height(16.dp))
        }

        if (uiState.hasDiscordIdentity) {
            DestinationOptionRow(
                icon = Icons.Default.Chat,
                label = stringResource(R.string.dest_discord_dm),
                onClick = onAddDiscordDm,
            )
        }
        if (uiState.hasEmailIdentity) {
            DestinationOptionRow(
                icon = Icons.Default.AlternateEmail,
                label = stringResource(R.string.dest_email),
                onClick = { showEmailInput = true },
            )
        }
        DestinationOptionRow(
            icon = Icons.Default.Notifications,
            label = stringResource(R.string.dest_push),
            onClick = onAddPush,
        )
        if (uiState.hasDiscordIdentity) {
            DestinationOptionRow(
                icon = Icons.Default.Tag,
                label = stringResource(R.string.dest_discord_channel),
                onClick = {
                    onLoadGuilds()
                    showChannelSheet = true
                },
            )
        }
        DestinationOptionRow(
            icon = Icons.Default.Webhook,
            label = stringResource(R.string.dest_webhook),
            onClick = { showWebhookInput = true },
        )
    }

    if (showEmailInput) {
        ModalBottomSheet(onDismissRequest = { showEmailInput = false }, containerColor = BackgroundCard) {
            var email by rememberSaveable { mutableStateOf(uiState.accountEmail.orEmpty()) }
            Column(Modifier.padding(16.dp)) {
                Text(stringResource(R.string.email_address), style = MaterialTheme.typography.titleMedium)
                Spacer(Modifier.height(12.dp))
                ChronosTextField(
                    value = email,
                    onValueChange = { email = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.email),
                )
                Spacer(Modifier.height(16.dp))
                ChronosButton(
                    text = stringResource(R.string.dest_add),
                    onClick = {
                        onAddEmail(email.trim())
                        showEmailInput = false
                    },
                    modifier = Modifier.fillMaxWidth(),
                    enabled = email.contains('@'),
                )
                Spacer(Modifier.height(24.dp))
            }
        }
    }

    if (showWebhookInput) {
        ModalBottomSheet(onDismissRequest = { showWebhookInput = false }, containerColor = BackgroundCard) {
            var url by rememberSaveable { mutableStateOf("") }
            var platform by rememberSaveable { mutableStateOf("generic") }
            Column(Modifier.padding(16.dp)) {
                Text(stringResource(R.string.dest_webhook), style = MaterialTheme.typography.titleMedium)
                Spacer(Modifier.height(12.dp))
                ChronosTextField(
                    value = url,
                    onValueChange = { url = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.webhook_url),
                )
                Spacer(Modifier.height(12.dp))
                Text(
                    stringResource(R.string.webhook_platform),
                    style = MaterialTheme.typography.labelLarge,
                    color = ForegroundMuted,
                )
                Spacer(Modifier.height(8.dp))
                Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                    listOf("generic", "discord", "slack").forEach { option ->
                        FilterChip(
                            selected = platform == option,
                            onClick = { platform = option },
                            label = { Text(option, style = MaterialTheme.typography.labelLarge) },
                            colors = FilterChipDefaults.filterChipColors(
                                containerColor = BackgroundMuted,
                                labelColor = ForegroundMuted,
                                selectedContainerColor = AccentOrange,
                                selectedLabelColor = ForegroundMain,
                            ),
                        )
                    }
                }
                Spacer(Modifier.height(16.dp))
                ChronosButton(
                    text = stringResource(R.string.dest_add),
                    onClick = {
                        onAddWebhook(url.trim(), platform)
                        showWebhookInput = false
                    },
                    modifier = Modifier.fillMaxWidth(),
                    enabled = url.startsWith("http"),
                )
                Spacer(Modifier.height(24.dp))
            }
        }
    }

    if (showChannelSheet) {
        ModalBottomSheet(onDismissRequest = { showChannelSheet = false }, containerColor = BackgroundCard) {
            var selectedGuild by rememberSaveable { mutableStateOf<String?>(null) }
            val guild = uiState.guilds.firstOrNull { it.id == selectedGuild }
            Column(Modifier.padding(16.dp)) {
                Text(
                    text = if (guild == null) {
                        stringResource(R.string.select_guild)
                    } else {
                        stringResource(R.string.select_channel)
                    },
                    style = MaterialTheme.typography.titleMedium,
                )
                Spacer(Modifier.height(12.dp))
                if (uiState.loadingGuilds || uiState.loadingChannels) {
                    CircularProgressIndicator(
                        color = AccentOrange,
                        modifier = Modifier
                            .align(Alignment.CenterHorizontally)
                            .padding(24.dp),
                    )
                } else if (guild == null) {
                    LazyColumn {
                        items(uiState.guilds, key = { it.id }) { g ->
                            Text(
                                text = g.name,
                                style = MaterialTheme.typography.bodyLarge,
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .clickable {
                                        selectedGuild = g.id
                                        onLoadChannels(g.id)
                                    }
                                    .padding(vertical = 14.dp),
                            )
                        }
                    }
                } else {
                    LazyColumn {
                        items(
                            uiState.channels.filter { it.type == DISCORD_TEXT_CHANNEL },
                            key = { it.id },
                        ) { channel ->
                            Text(
                                text = "#${channel.name}",
                                style = MaterialTheme.typography.bodyLarge,
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .clickable {
                                        onAddChannel(guild, channel)
                                        showChannelSheet = false
                                    }
                                    .padding(vertical = 14.dp),
                            )
                        }
                    }
                }
                Spacer(Modifier.height(24.dp))
            }
        }
    }
}

@Composable
private fun DestinationOptionRow(
    icon: ImageVector,
    label: String,
    onClick: () -> Unit,
) {
    ChronosCard(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 4.dp),
        onClick = onClick,
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Icon(imageVector = icon, contentDescription = null, tint = AccentOrange, modifier = Modifier.size(22.dp))
            Spacer(Modifier.padding(start = 12.dp))
            Text(label, style = MaterialTheme.typography.bodyLarge, modifier = Modifier.weight(1f))
            Text(
                stringResource(R.string.dest_add),
                style = MaterialTheme.typography.labelLarge,
                color = ForegroundMuted,
            )
        }
    }
}
