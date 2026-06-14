package com.chronos.reminder.account.ui

import android.content.ClipData
import android.content.ClipboardManager
import android.content.Context
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
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.Key
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.components.ChronosButton
import com.chronos.reminder.core.ui.components.ChronosButtonStyle
import com.chronos.reminder.core.ui.components.ChronosCard
import com.chronos.reminder.core.ui.components.ChronosTextField
import com.chronos.reminder.core.ui.components.ChronosTopBar
import com.chronos.reminder.core.ui.components.ConfirmDeleteDialog
import com.chronos.reminder.core.ui.components.EmptyState
import com.chronos.reminder.core.ui.components.ErrorBanner
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.BackgroundCard
import com.chronos.reminder.core.ui.theme.BackgroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMuted

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ApiKeysScreen(
    onBack: () -> Unit,
    viewModel: AccountViewModel = hiltViewModel(),
) {
    val state by viewModel.state.collectAsStateWithLifecycle()
    val context = LocalContext.current

    var showCreateSheet by rememberSaveable { mutableStateOf(false) }
    var newKeyName by rememberSaveable { mutableStateOf("") }
    var revokeCandidate by rememberSaveable { mutableStateOf<String?>(null) }

    LaunchedEffect(Unit) { viewModel.loadApiKeys() }

    Scaffold(
        containerColor = BackgroundMain,
        topBar = { ChronosTopBar(title = stringResource(R.string.api_keys_title), onBack = onBack) },
        floatingActionButton = {
            FloatingActionButton(
                onClick = {
                    newKeyName = ""
                    showCreateSheet = true
                },
                containerColor = AccentOrange,
                contentColor = ForegroundMain,
                modifier = Modifier.size(56.dp),
            ) {
                Icon(Icons.Default.Add, contentDescription = stringResource(R.string.create_api_key))
            }
        },
    ) { padding ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding),
        ) {
            if (state.apiKeys.isEmpty()) {
                EmptyState(
                    icon = Icons.Default.Key,
                    title = stringResource(R.string.empty_api_keys_title),
                    subtitle = stringResource(R.string.empty_api_keys_subtitle),
                )
            } else {
                LazyColumn(
                    modifier = Modifier.fillMaxSize(),
                    verticalArrangement = Arrangement.spacedBy(8.dp),
                    contentPadding = androidx.compose.foundation.layout.PaddingValues(16.dp),
                ) {
                    items(state.apiKeys, key = { it.id }) { apiKey ->
                        ChronosCard(modifier = Modifier.fillMaxWidth()) {
                            Row(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .padding(16.dp),
                                verticalAlignment = Alignment.CenterVertically,
                            ) {
                                Column(Modifier.weight(1f)) {
                                    Text(apiKey.name, style = MaterialTheme.typography.titleMedium)
                                    if (apiKey.scopes.isNotBlank()) {
                                        Text(
                                            apiKey.scopes,
                                            style = MaterialTheme.typography.labelSmall,
                                            color = ForegroundMuted,
                                        )
                                    }
                                    apiKey.createdAt?.let {
                                        Text(
                                            it.take(10),
                                            style = MaterialTheme.typography.labelSmall,
                                            color = ForegroundMuted,
                                        )
                                    }
                                }
                                ChronosButton(
                                    text = stringResource(R.string.revoke),
                                    onClick = { revokeCandidate = apiKey.id },
                                    style = ChronosButtonStyle.Destructive,
                                )
                            }
                        }
                    }
                }
            }

            ErrorBanner(
                message = state.error,
                modifier = Modifier.align(Alignment.BottomCenter),
                onDismiss = viewModel::clearMessages,
            )
        }
    }

    if (showCreateSheet) {
        ModalBottomSheet(onDismissRequest = { showCreateSheet = false }, containerColor = BackgroundCard) {
            Column(Modifier.padding(16.dp)) {
                Text(stringResource(R.string.create_api_key), style = MaterialTheme.typography.titleMedium)
                Spacer(Modifier.height(12.dp))
                ChronosTextField(
                    value = newKeyName,
                    onValueChange = { newKeyName = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.api_key_name),
                )
                Spacer(Modifier.height(16.dp))
                ChronosButton(
                    text = stringResource(R.string.create),
                    onClick = {
                        viewModel.createApiKey(newKeyName)
                        showCreateSheet = false
                    },
                    modifier = Modifier.fillMaxWidth(),
                    enabled = newKeyName.isNotBlank(),
                )
                Spacer(Modifier.height(32.dp))
            }
        }
    }

    state.createdKey?.let { created ->
        AlertDialog(
            onDismissRequest = viewModel::dismissCreatedKey,
            containerColor = BackgroundCard,
            title = { Text(stringResource(R.string.api_key_created_title), style = MaterialTheme.typography.titleLarge) },
            text = {
                Column {
                    Text(
                        created.key.orEmpty(),
                        style = MaterialTheme.typography.bodyMedium,
                        color = AccentOrange,
                    )
                    Spacer(Modifier.height(8.dp))
                    Text(
                        stringResource(R.string.api_key_created_warning),
                        style = MaterialTheme.typography.labelSmall,
                        color = ForegroundMuted,
                    )
                }
            },
            confirmButton = {
                TextButton(onClick = {
                    val clipboard = context.getSystemService(Context.CLIPBOARD_SERVICE) as ClipboardManager
                    clipboard.setPrimaryClip(ClipData.newPlainText("Chronos API key", created.key.orEmpty()))
                }) {
                    Text(stringResource(R.string.copy_to_clipboard), color = AccentOrange)
                }
            },
            dismissButton = {
                TextButton(onClick = viewModel::dismissCreatedKey) {
                    Text(stringResource(R.string.ok), color = ForegroundMuted)
                }
            },
        )
    }

    revokeCandidate?.let { id ->
        ConfirmDeleteDialog(
            title = stringResource(R.string.revoke_key_title),
            text = stringResource(R.string.revoke_key_text),
            confirmLabel = stringResource(R.string.revoke),
            onConfirm = {
                viewModel.revokeApiKey(id)
                revokeCandidate = null
            },
            onDismiss = { revokeCandidate = null },
        )
    }
}
