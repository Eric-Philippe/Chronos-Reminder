package com.chronos.reminder.home.ui

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.TaskAlt
import androidx.compose.material.icons.filled.Notifications
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Schedule
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.drawWithContent
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import com.chronos.reminder.core.ui.theme.BackgroundMain
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.components.ChronosButton
import com.chronos.reminder.core.ui.components.ChronosButtonStyle
import com.chronos.reminder.core.ui.components.ChronosCard
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.ForegroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMuted

@Composable
fun HomeScreen(
    onCreateReminder: () -> Unit,
    onOpenReminders: () -> Unit,
    onOpenDfm: () -> Unit,
    onOpenAccount: () -> Unit,
    viewModel: HomeViewModel = hiltViewModel(),
) {
    val state by viewModel.state.collectAsStateWithLifecycle()

    Column(
        modifier = Modifier
            .fillMaxSize()
            .verticalScroll(rememberScrollState())
            .padding(horizontal = 16.dp),
    ) {
        Spacer(Modifier.height(32.dp))

        // Logo + title header
        Column(
            modifier = Modifier.fillMaxWidth(),
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            Image(
                painter = painterResource(R.drawable.hourglass),
                contentDescription = null,
                modifier = Modifier
                    .size(120.dp)
                    .drawWithContent {
                        drawContent()
                        drawRect(
                            brush = Brush.radialGradient(
                                colorStops = arrayOf(
                                    0.55f to Color.Transparent,
                                    1.0f to BackgroundMain,
                                ),
                                center = Offset(size.width / 2f, size.height / 2f),
                                radius = maxOf(size.width, size.height) / 2f,
                            ),
                        )
                    },
            )
            Spacer(Modifier.height(12.dp))
            Text(
                text = stringResource(R.string.home_welcome_title),
                style = MaterialTheme.typography.displayMedium,
                color = AccentOrange,
            )
            Spacer(Modifier.height(6.dp))
            val identities = state.account?.identities ?: emptyList()
            val username = state.account?.username ?: identities.firstOrNull()?.username
            if (!username.isNullOrEmpty()) {
                Text(
                    text = stringResource(R.string.home_greeting, username),
                    style = MaterialTheme.typography.bodyMedium,
                    color = ForegroundMuted,
                )
            } else {
                Text(
                    text = stringResource(R.string.login_subtitle),
                    style = MaterialTheme.typography.bodyMedium,
                    color = ForegroundMuted,
                )
            }
        }

        Spacer(Modifier.height(32.dp))

        // Overview stats
        Text(
            text = stringResource(R.string.home_section_overview),
            style = MaterialTheme.typography.titleMedium,
            color = AccentOrange,
        )
        Spacer(Modifier.height(12.dp))
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.spacedBy(10.dp),
        ) {
            StatCard(
                modifier = Modifier.weight(1f),
                label = stringResource(R.string.home_total_reminders),
                value = state.totalReminders.toString(),
                icon = Icons.Default.Schedule,
                onClick = onOpenReminders,
            )
            StatCard(
                modifier = Modifier.weight(1f),
                label = stringResource(R.string.home_active_reminders),
                value = state.activeReminders.toString(),
                icon = Icons.Default.Notifications,
                onClick = onOpenReminders,
            )
            StatCard(
                modifier = Modifier.weight(1f),
                label = stringResource(R.string.home_dfm_items),
                value = state.dfmItemCount.toString(),
                icon = Icons.Default.TaskAlt,
                onClick = onOpenDfm,
            )
        }

        Spacer(Modifier.height(28.dp))

        // Quick access
        Text(
            text = stringResource(R.string.home_section_quick_access),
            style = MaterialTheme.typography.titleMedium,
            color = AccentOrange,
        )
        Spacer(Modifier.height(12.dp))

        ChronosButton(
            text = stringResource(R.string.home_btn_create_reminder),
            onClick = onCreateReminder,
            modifier = Modifier.fillMaxWidth(),
            leadingIcon = {
                Icon(
                    Icons.Default.Add,
                    contentDescription = null,
                    modifier = Modifier.size(18.dp),
                )
            },
        )
        Spacer(Modifier.height(8.dp))
        ChronosButton(
            text = stringResource(R.string.home_btn_open_dfm),
            onClick = onOpenDfm,
            modifier = Modifier.fillMaxWidth(),
            style = ChronosButtonStyle.Secondary,
            leadingIcon = {
                Icon(
                    Icons.Default.TaskAlt,
                    contentDescription = null,
                    modifier = Modifier.size(18.dp),
                )
            },
        )
        Spacer(Modifier.height(8.dp))
        ChronosButton(
            text = stringResource(R.string.home_btn_account),
            onClick = onOpenAccount,
            modifier = Modifier.fillMaxWidth(),
            style = ChronosButtonStyle.Secondary,
            leadingIcon = {
                Icon(
                    Icons.Default.Person,
                    contentDescription = null,
                    modifier = Modifier.size(18.dp),
                )
            },
        )

        Spacer(Modifier.height(48.dp))
    }
}

@Composable
private fun StatCard(
    modifier: Modifier = Modifier,
    label: String,
    value: String,
    icon: ImageVector,
    onClick: (() -> Unit)? = null,
) {
    ChronosCard(modifier = modifier, onClick = onClick) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 14.dp, horizontal = 8.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            Icon(
                imageVector = icon,
                contentDescription = null,
                tint = AccentOrange,
                modifier = Modifier.size(20.dp),
            )
            Spacer(Modifier.height(6.dp))
            Text(
                text = value,
                style = MaterialTheme.typography.titleLarge,
                color = ForegroundMain,
                textAlign = TextAlign.Center,
            )
            Spacer(Modifier.height(2.dp))
            Text(
                text = label,
                style = MaterialTheme.typography.labelSmall,
                color = ForegroundMuted,
                textAlign = TextAlign.Center,
            )
        }
    }
}

