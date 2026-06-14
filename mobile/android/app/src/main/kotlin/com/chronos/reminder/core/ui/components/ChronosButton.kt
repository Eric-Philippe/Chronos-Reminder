package com.chronos.reminder.core.ui.components

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.BorderColor
import com.chronos.reminder.core.ui.theme.CornerRadius
import com.chronos.reminder.core.ui.theme.DestructiveRed
import com.chronos.reminder.core.ui.theme.ForegroundMain

enum class ChronosButtonStyle { Primary, Secondary, Destructive }

@Composable
fun ChronosButton(
    text: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    style: ChronosButtonStyle = ChronosButtonStyle.Primary,
    enabled: Boolean = true,
    loading: Boolean = false,
    leadingIcon: (@Composable () -> Unit)? = null,
) {
    val shape = RoundedCornerShape(CornerRadius)
    val content: @Composable () -> Unit = {
        Row(verticalAlignment = Alignment.CenterVertically) {
            if (loading) {
                CircularProgressIndicator(
                    modifier = Modifier
                        .size(18.dp)
                        .padding(end = 4.dp),
                    strokeWidth = 2.dp,
                    color = ForegroundMain,
                )
            } else if (leadingIcon != null) {
                leadingIcon()
                Spacer(Modifier.width(8.dp))
            }
            Text(text = text, style = MaterialTheme.typography.labelLarge)
        }
    }

    when (style) {
        ChronosButtonStyle.Primary -> Button(
            onClick = onClick,
            modifier = modifier.height(48.dp),
            enabled = enabled && !loading,
            shape = shape,
            colors = ButtonDefaults.buttonColors(containerColor = AccentOrange, contentColor = ForegroundMain),
        ) { content() }

        ChronosButtonStyle.Secondary -> OutlinedButton(
            onClick = onClick,
            modifier = modifier.height(48.dp),
            enabled = enabled && !loading,
            shape = shape,
            border = BorderStroke(1.dp, BorderColor),
            colors = ButtonDefaults.outlinedButtonColors(contentColor = ForegroundMain),
        ) { content() }

        ChronosButtonStyle.Destructive -> Button(
            onClick = onClick,
            modifier = modifier.height(48.dp),
            enabled = enabled && !loading,
            shape = shape,
            colors = ButtonDefaults.buttonColors(containerColor = DestructiveRed, contentColor = ForegroundMain),
        ) { content() }
    }
}
