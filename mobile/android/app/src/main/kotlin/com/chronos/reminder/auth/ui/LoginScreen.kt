package com.chronos.reminder.auth.ui

import android.net.Uri
import androidx.browser.customtabs.CustomTabsIntent
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
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Visibility
import androidx.compose.material.icons.filled.VisibilityOff
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.unit.dp
import com.chronos.reminder.BuildConfig
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.components.ChronosButton
import com.chronos.reminder.core.ui.components.ChronosButtonStyle
import com.chronos.reminder.core.ui.components.ChronosTextField
import com.chronos.reminder.core.ui.components.ErrorBanner
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.ForegroundMuted

const val DISCORD_REDIRECT_URI = "chronos://auth/discord"

fun discordOAuthUrl(): String =
    "https://discord.com/oauth2/authorize" +
        "?client_id=${BuildConfig.DISCORD_CLIENT_ID}" +
        "&redirect_uri=${Uri.encode(DISCORD_REDIRECT_URI)}" +
        "&response_type=code" +
        "&scope=identify+email"

@Composable
fun LoginScreen(
    uiState: AuthUiState,
    onLogin: (email: String, password: String) -> Unit,
    onNavigateRegister: () -> Unit,
    onNavigateForgotPassword: () -> Unit,
    onClearError: () -> Unit,
    onDiscordUnconfigured: () -> Unit,
    onResendVerification: (email: String) -> Unit = {},
) {
    val context = LocalContext.current
    var email by rememberSaveable { mutableStateOf("") }
    var password by rememberSaveable { mutableStateOf("") }
    var showPassword by rememberSaveable { mutableStateOf(false) }

    val isEmailNotVerified = uiState.error?.contains("verified", ignoreCase = true) == true
        || uiState.error?.contains("not verified", ignoreCase = true) == true

    androidx.compose.foundation.layout.Box(modifier = Modifier.fillMaxSize()) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .verticalScroll(rememberScrollState())
            .padding(horizontal = 16.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center,
    ) {
        Image(
            painter = painterResource(R.drawable.logo_chronos),
            contentDescription = stringResource(R.string.chronos_logo),
            modifier = Modifier
                .size(96.dp)
                .clip(RoundedCornerShape(22.dp)),
        )
        Spacer(Modifier.height(8.dp))
        Text(
            text = stringResource(R.string.login_title),
            style = MaterialTheme.typography.displayMedium,
            color = AccentOrange,
        )
        Text(
            text = stringResource(R.string.login_subtitle),
            style = MaterialTheme.typography.bodyMedium,
            color = ForegroundMuted,
        )
        Spacer(Modifier.height(32.dp))

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
            visualTransformation = if (showPassword) VisualTransformation.None else PasswordVisualTransformation(),
            trailingIcon = {
                IconButton(onClick = { showPassword = !showPassword }) {
                    Icon(
                        imageVector = if (showPassword) Icons.Default.VisibilityOff else Icons.Default.Visibility,
                        contentDescription = stringResource(
                            if (showPassword) R.string.hide_password else R.string.show_password,
                        ),
                        tint = ForegroundMuted,
                    )
                }
            },
        )
        Spacer(Modifier.height(20.dp))

        ChronosButton(
            text = stringResource(R.string.sign_in),
            onClick = { onLogin(email, password) },
            modifier = Modifier.fillMaxWidth(),
            loading = uiState.loading,
        )

        Spacer(Modifier.height(16.dp))
        Row(verticalAlignment = Alignment.CenterVertically) {
            HorizontalDivider(modifier = Modifier.weight(1f))
            Text(
                text = stringResource(R.string.or_divider),
                style = MaterialTheme.typography.labelSmall,
                color = ForegroundMuted,
                modifier = Modifier.padding(horizontal = 12.dp),
            )
            HorizontalDivider(modifier = Modifier.weight(1f))
        }
        Spacer(Modifier.height(16.dp))

        ChronosButton(
            text = stringResource(R.string.continue_with_discord),
            onClick = {
                if (BuildConfig.DISCORD_CLIENT_ID.isBlank()) {
                    onDiscordUnconfigured()
                } else {
                    CustomTabsIntent.Builder().build()
                        .launchUrl(context, Uri.parse(discordOAuthUrl()))
                }
            },
            modifier = Modifier.fillMaxWidth(),
            style = ChronosButtonStyle.Secondary,
            leadingIcon = {
                Icon(
                    painter = painterResource(R.drawable.ic_discord),
                    contentDescription = stringResource(R.string.discord_logo),
                    modifier = Modifier
                        .size(20.dp)
                        .padding(end = 2.dp),
                    tint = androidx.compose.ui.graphics.Color.Unspecified,
                )
            },
        )

        Spacer(Modifier.height(24.dp))
        TextButton(onClick = onNavigateForgotPassword) {
            Text(
                stringResource(R.string.forgot_password_link),
                style = MaterialTheme.typography.bodyMedium,
                color = ForegroundMuted,
            )
        }
        TextButton(onClick = onNavigateRegister) {
            Text(
                stringResource(R.string.no_account_sign_up),
                style = MaterialTheme.typography.bodyMedium,
                color = AccentOrange,
            )
        }

        // Show resend verification button when the error is about unverified email
        if (isEmailNotVerified && email.isNotBlank()) {
            Spacer(Modifier.height(8.dp))
            ChronosButton(
                text = when (uiState.resendState) {
                    ResendState.Sending -> stringResource(R.string.resend_verification_sending)
                    ResendState.Sent -> stringResource(R.string.resend_verification_sent)
                    else -> stringResource(R.string.resend_verification)
                },
                onClick = { onResendVerification(email) },
                modifier = Modifier.fillMaxWidth(),
                style = ChronosButtonStyle.Secondary,
                enabled = uiState.resendState == ResendState.Idle || uiState.resendState == ResendState.Error,
                loading = uiState.resendState == ResendState.Sending,
            )
        }
    }

    ErrorBanner(
        message = if (uiState.resendState == ResendState.Sent) null else uiState.error,
        modifier = Modifier.align(Alignment.BottomCenter),
        onDismiss = onClearError,
    )
    }
}
