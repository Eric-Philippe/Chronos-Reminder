package com.chronos.reminder.core.ui.components

import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.OutlinedTextFieldDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.VisualTransformation
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.BackgroundMuted
import com.chronos.reminder.core.ui.theme.BorderColor
import com.chronos.reminder.core.ui.theme.CornerRadius
import com.chronos.reminder.core.ui.theme.DestructiveRed
import com.chronos.reminder.core.ui.theme.ForegroundMuted

@Composable
fun ChronosTextField(
    value: String,
    onValueChange: (String) -> Unit,
    modifier: Modifier = Modifier,
    placeholder: String? = null,
    singleLine: Boolean = true,
    minLines: Int = 1,
    enabled: Boolean = true,
    isError: Boolean = false,
    errorText: String? = null,
    visualTransformation: VisualTransformation = VisualTransformation.None,
    keyboardOptions: KeyboardOptions = KeyboardOptions.Default,
    keyboardActions: KeyboardActions = KeyboardActions.Default,
    trailingIcon: (@Composable () -> Unit)? = null,
) {
    OutlinedTextField(
        value = value,
        onValueChange = onValueChange,
        modifier = modifier,
        enabled = enabled,
        singleLine = singleLine,
        minLines = minLines,
        isError = isError,
        shape = RoundedCornerShape(CornerRadius),
        placeholder = placeholder?.let {
            { Text(it, style = MaterialTheme.typography.bodyLarge, color = ForegroundMuted) }
        },
        supportingText = errorText?.let {
            { Text(it, style = MaterialTheme.typography.labelSmall, color = DestructiveRed) }
        },
        visualTransformation = visualTransformation,
        keyboardOptions = keyboardOptions,
        keyboardActions = keyboardActions,
        trailingIcon = trailingIcon,
        textStyle = MaterialTheme.typography.bodyLarge,
        colors = OutlinedTextFieldDefaults.colors(
            focusedContainerColor = BackgroundMuted,
            unfocusedContainerColor = BackgroundMuted,
            disabledContainerColor = BackgroundMuted,
            errorContainerColor = BackgroundMuted,
            focusedBorderColor = AccentOrange,
            unfocusedBorderColor = BorderColor,
            errorBorderColor = DestructiveRed,
            cursorColor = AccentOrange,
        ),
    )
}
