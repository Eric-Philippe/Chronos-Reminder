package com.chronos.reminder.auth.ui

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.imePadding
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.ExperimentalMaterial3Api
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
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.components.ChronosButton
import com.chronos.reminder.core.ui.components.ChronosTextField
import com.chronos.reminder.core.ui.components.ChronosTopBar
import com.chronos.reminder.core.ui.components.ErrorBanner
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.BackgroundCard
import com.chronos.reminder.core.ui.theme.BackgroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMuted
import java.util.TimeZone

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DiscordSetupScreen(
    uiState: AuthUiState,
    setup: DiscordSetupState,
    onComplete: (email: String, username: String, password: String, timezone: String) -> Unit,
    onLoadTimezones: () -> Unit,
    onCancel: () -> Unit,
    onClearError: () -> Unit,
) {
    var username by rememberSaveable { mutableStateOf(setup.username ?: "") }
    var email by rememberSaveable { mutableStateOf(setup.email ?: "") }
    var password by rememberSaveable { mutableStateOf("") }
    var passwordConfirm by rememberSaveable { mutableStateOf("") }
    var timezone by rememberSaveable { mutableStateOf(TimeZone.getDefault().id) }
    var showTimezoneSheet by rememberSaveable { mutableStateOf(false) }
    var localError by rememberSaveable { mutableStateOf<String?>(null) }

    val mismatchError = stringResource(R.string.error_passwords_mismatch)
    val invalidEmailError = stringResource(R.string.error_invalid_email)
    val usernameRequiredError = stringResource(R.string.error_username_required)
    val passwordRequiredError = stringResource(R.string.error_password_required)

    LaunchedEffect(Unit) { onLoadTimezones() }

    Scaffold(
        containerColor = BackgroundMain,
        topBar = {
            ChronosTopBar(title = stringResource(R.string.discord_setup_title), onBack = onCancel)
        },
    ) { padding ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding),
        ) {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .verticalScroll(rememberScrollState())
                    .imePadding()
                    .padding(horizontal = 16.dp),
            ) {
                Spacer(Modifier.height(16.dp))
                Text(
                    text = stringResource(R.string.discord_setup_subtitle),
                    style = MaterialTheme.typography.bodyMedium,
                    color = ForegroundMuted,
                )
                Spacer(Modifier.height(16.dp))
                ChronosTextField(
                    value = username,
                    onValueChange = { username = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.username),
                )
                Spacer(Modifier.height(12.dp))
                ChronosTextField(
                    value = email,
                    onValueChange = { email = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.email),
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
                )
                Spacer(Modifier.height(12.dp))
                ChronosTextField(
                    value = password,
                    onValueChange = { password = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.password),
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
                    visualTransformation = PasswordVisualTransformation(),
                )
                Spacer(Modifier.height(12.dp))
                ChronosTextField(
                    value = passwordConfirm,
                    onValueChange = { passwordConfirm = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.confirm_password),
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
                    visualTransformation = PasswordVisualTransformation(),
                    isError = passwordConfirm.isNotEmpty() && passwordConfirm != password,
                )
                Spacer(Modifier.height(12.dp))

                Text(
                    text = stringResource(R.string.timezone),
                    style = MaterialTheme.typography.labelLarge,
                    color = ForegroundMuted,
                )
                Spacer(Modifier.height(4.dp))
                Text(
                    text = timezone,
                    style = MaterialTheme.typography.bodyLarge,
                    color = AccentOrange,
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { showTimezoneSheet = true }
                        .padding(vertical = 12.dp),
                )
                Spacer(Modifier.height(20.dp))

                ChronosButton(
                    text = stringResource(R.string.discord_setup_submit),
                    onClick = {
                        localError = when {
                            username.isBlank() -> usernameRequiredError
                            !email.contains('@') -> invalidEmailError
                            password.length < 8 -> passwordRequiredError
                            password != passwordConfirm -> mismatchError
                            else -> null
                        }
                        if (localError == null) {
                            onComplete(email, username, password, timezone)
                        }
                    },
                    modifier = Modifier.fillMaxWidth(),
                    loading = uiState.loading,
                )
                Spacer(Modifier.height(24.dp))
            }

            ErrorBanner(
                message = localError ?: uiState.error,
                modifier = Modifier.align(Alignment.BottomCenter),
                onDismiss = {
                    localError = null
                    onClearError()
                },
            )
        }
    }

    if (showTimezoneSheet) {
        ModalBottomSheet(
            onDismissRequest = { showTimezoneSheet = false },
            containerColor = BackgroundCard,
        ) {
            TimezonePickerSheetContent(
                timezones = uiState.timezones,
                onSelect = {
                    timezone = it
                    showTimezoneSheet = false
                },
            )
        }
    }
}
