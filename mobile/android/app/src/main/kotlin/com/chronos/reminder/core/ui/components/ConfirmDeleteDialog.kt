package com.chronos.reminder.core.ui.components

import androidx.compose.material3.AlertDialog
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.ui.res.stringResource
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.theme.BackgroundCard
import com.chronos.reminder.core.ui.theme.DestructiveRed
import com.chronos.reminder.core.ui.theme.ForegroundMuted

@Composable
fun ConfirmDeleteDialog(
    title: String,
    text: String,
    onConfirm: () -> Unit,
    onDismiss: () -> Unit,
    confirmLabel: String = stringResource(R.string.delete),
    confirmEnabled: Boolean = true,
    extraContent: (@Composable () -> Unit)? = null,
) {
    AlertDialog(
        onDismissRequest = onDismiss,
        containerColor = BackgroundCard,
        title = { Text(title, style = MaterialTheme.typography.titleLarge) },
        text = {
            if (extraContent != null) {
                androidx.compose.foundation.layout.Column {
                    Text(text, style = MaterialTheme.typography.bodyMedium)
                    extraContent()
                }
            } else {
                Text(text, style = MaterialTheme.typography.bodyMedium)
            }
        },
        confirmButton = {
            TextButton(onClick = onConfirm, enabled = confirmEnabled) {
                Text(confirmLabel, color = DestructiveRed, style = MaterialTheme.typography.labelLarge)
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) {
                Text(
                    stringResource(R.string.cancel),
                    color = ForegroundMuted,
                    style = MaterialTheme.typography.labelLarge,
                )
            }
        },
    )
}
