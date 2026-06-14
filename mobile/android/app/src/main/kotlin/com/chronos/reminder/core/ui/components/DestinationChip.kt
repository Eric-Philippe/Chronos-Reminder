package com.chronos.reminder.core.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AlternateEmail
import androidx.compose.material.icons.filled.Chat
import androidx.compose.material.icons.filled.Notifications
import androidx.compose.material.icons.filled.Tag
import androidx.compose.material.icons.filled.Webhook
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material.icons.filled.Close
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.theme.BackgroundMuted
import com.chronos.reminder.core.ui.theme.ForegroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMuted
import com.chronos.reminder.reminders.domain.Destination

@Composable
fun destinationLabel(type: String): String = when (type) {
    Destination.TYPE_DISCORD_DM -> stringResource(R.string.dest_discord_dm)
    Destination.TYPE_DISCORD_CHANNEL -> stringResource(R.string.dest_discord_channel)
    Destination.TYPE_EMAIL -> stringResource(R.string.dest_email)
    Destination.TYPE_WEBHOOK -> stringResource(R.string.dest_webhook)
    Destination.TYPE_ANDROID_PUSH -> stringResource(R.string.dest_push)
    else -> type
}

@Composable
fun DestinationChip(
    type: String,
    modifier: Modifier = Modifier,
    detail: String? = null,
    onRemove: (() -> Unit)? = null,
) {
    Row(
        modifier = modifier
            .background(BackgroundMuted, RoundedCornerShape(16.dp))
            .padding(horizontal = 10.dp, vertical = 6.dp),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        Icon(
            imageVector = when (type) {
                Destination.TYPE_DISCORD_DM -> Icons.Default.Chat
                Destination.TYPE_DISCORD_CHANNEL -> Icons.Default.Tag
                Destination.TYPE_EMAIL -> Icons.Default.AlternateEmail
                Destination.TYPE_WEBHOOK -> Icons.Default.Webhook
                else -> Icons.Default.Notifications
            },
            contentDescription = null,
            modifier = Modifier
                .size(14.dp)
                .padding(end = 0.dp),
            tint = ForegroundMuted,
        )
        Text(
            text = detail?.let { "${destinationLabel(type)} · $it" } ?: destinationLabel(type),
            style = MaterialTheme.typography.labelSmall,
            color = ForegroundMain,
            modifier = Modifier.padding(start = 4.dp),
        )
        if (onRemove != null) {
            IconButton(onClick = onRemove, modifier = Modifier.size(20.dp)) {
                Icon(
                    imageVector = Icons.Default.Close,
                    contentDescription = stringResource(R.string.dest_remove),
                    modifier = Modifier.size(14.dp),
                    tint = ForegroundMuted,
                )
            }
        }
    }
}
