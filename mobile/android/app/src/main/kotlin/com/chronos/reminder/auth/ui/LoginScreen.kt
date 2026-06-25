package com.chronos.reminder.auth.ui

import android.net.Uri
import androidx.browser.customtabs.CustomTabsIntent
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
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
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Visibility
import androidx.compose.material.icons.filled.VisibilityOff
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.ModalBottomSheet
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
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
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
import com.chronos.reminder.core.ui.theme.BackgroundCard
import com.chronos.reminder.core.ui.theme.ForegroundMuted
import java.util.TimeZone

const val DISCORD_REDIRECT_URI = "chronos://auth/discord"

fun discordOAuthUrl(): String =
    "https://discord.com/oauth2/authorize" +
        "?client_id=${BuildConfig.DISCORD_CLIENT_ID}" +
        "&redirect_uri=${Uri.encode(DISCORD_REDIRECT_URI)}" +
        "&response_type=code" +
        "&scope=identify+email+guilds+guilds.members.read"

enum class AuthTab { LOGIN, REGISTER }

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LoginScreen(
    uiState: AuthUiState,
    onLogin: (email: String, password: String) -> Unit,
    onRegister: (email: String, username: String, password: String, timezone: String) -> Unit = { _, _, _, _ -> },
    onLoadTimezones: () -> Unit = {},
    onNavigateForgotPassword: () -> Unit,
    onClearError: () -> Unit,
    onDiscordUnconfigured: () -> Unit,
    onResendVerification: (email: String) -> Unit = {},
) {
    val context = LocalContext.current

    var selectedTab by rememberSaveable { mutableStateOf(AuthTab.LOGIN) }

    // Login state
    var loginEmail by rememberSaveable { mutableStateOf("") }
    var loginPassword by rememberSaveable { mutableStateOf("") }
    var showPassword by rememberSaveable { mutableStateOf(false) }

    // Register state
    var regUsername by rememberSaveable { mutableStateOf("") }
    var regEmail by rememberSaveable { mutableStateOf("") }
    var regPassword by rememberSaveable { mutableStateOf("") }
    var regPasswordConfirm by rememberSaveable { mutableStateOf("") }
    var regTimezone by rememberSaveable { mutableStateOf(TimeZone.getDefault().id) }
    var showTimezoneSheet by rememberSaveable { mutableStateOf(false) }
    var localError by rememberSaveable { mutableStateOf<String?>(null) }

    val mismatchError = stringResource(R.string.error_passwords_mismatch)
    val invalidEmailError = stringResource(R.string.error_invalid_email)
    val usernameRequiredError = stringResource(R.string.error_username_required)
    val passwordRequiredError = stringResource(R.string.error_password_required)

    val isEmailNotVerified = uiState.error?.contains("verified", ignoreCase = true) == true
        || uiState.error?.contains("not verified", ignoreCase = true) == true

    LaunchedEffect(Unit) { onLoadTimezones() }

    LaunchedEffect(uiState.registerSuccess) {
        if (uiState.registerSuccess) {
            selectedTab = AuthTab.LOGIN
            regUsername = ""; regEmail = ""; regPassword = ""; regPasswordConfirm = ""
        }
    }

    Box(modifier = Modifier.fillMaxSize()) {
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
            Spacer(Modifier.height(28.dp))

            // Dual-panel tab switcher
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .clip(RoundedCornerShape(12.dp))
                    .background(BackgroundCard),
                verticalAlignment = Alignment.CenterVertically,
            ) {
                AuthTabItem(
                    label = stringResource(R.string.tab_have_account),
                    isSelected = selectedTab == AuthTab.LOGIN,
                    onClick = {
                        selectedTab = AuthTab.LOGIN
                        localError = null
                        onClearError()
                    },
                    modifier = Modifier.weight(1f),
                )
                Box(
                    modifier = Modifier
                        .width(1.dp)
                        .height(44.dp)
                        .background(ForegroundMuted.copy(alpha = 0.15f)),
                )
                AuthTabItem(
                    label = stringResource(R.string.tab_create_account),
                    isSelected = selectedTab == AuthTab.REGISTER,
                    onClick = {
                        selectedTab = AuthTab.REGISTER
                        localError = null
                        onClearError()
                    },
                    modifier = Modifier.weight(1f),
                )
            }

            Spacer(Modifier.height(24.dp))

            if (selectedTab == AuthTab.LOGIN) {
                // --- Login form ---
                ChronosTextField(
                    value = loginEmail,
                    onValueChange = { loginEmail = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.email),
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
                )
                Spacer(Modifier.height(12.dp))
                ChronosTextField(
                    value = loginPassword,
                    onValueChange = { loginPassword = it },
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
                    onClick = { onLogin(loginEmail, loginPassword) },
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
                            tint = Color.Unspecified,
                        )
                    },
                )

                Spacer(Modifier.height(16.dp))
                TextButton(onClick = onNavigateForgotPassword) {
                    Text(
                        stringResource(R.string.forgot_password_link),
                        style = MaterialTheme.typography.bodyMedium,
                        color = ForegroundMuted,
                    )
                }

                if (isEmailNotVerified && loginEmail.isNotBlank()) {
                    Spacer(Modifier.height(8.dp))
                    ChronosButton(
                        text = when (uiState.resendState) {
                            ResendState.Sending -> stringResource(R.string.resend_verification_sending)
                            ResendState.Sent -> stringResource(R.string.resend_verification_sent)
                            else -> stringResource(R.string.resend_verification)
                        },
                        onClick = { onResendVerification(loginEmail) },
                        modifier = Modifier.fillMaxWidth(),
                        style = ChronosButtonStyle.Secondary,
                        enabled = uiState.resendState == ResendState.Idle || uiState.resendState == ResendState.Error,
                        loading = uiState.resendState == ResendState.Sending,
                    )
                }
            } else {
                // --- Register form ---
                ChronosTextField(
                    value = regUsername,
                    onValueChange = { regUsername = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.username),
                )
                Spacer(Modifier.height(12.dp))
                ChronosTextField(
                    value = regEmail,
                    onValueChange = { regEmail = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.email),
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
                )
                Spacer(Modifier.height(12.dp))
                ChronosTextField(
                    value = regPassword,
                    onValueChange = { regPassword = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.password),
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
                    visualTransformation = PasswordVisualTransformation(),
                )
                Spacer(Modifier.height(12.dp))
                ChronosTextField(
                    value = regPasswordConfirm,
                    onValueChange = { regPasswordConfirm = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = stringResource(R.string.confirm_password),
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
                    visualTransformation = PasswordVisualTransformation(),
                    isError = regPasswordConfirm.isNotEmpty() && regPasswordConfirm != regPassword,
                )
                Spacer(Modifier.height(12.dp))

                Text(
                    text = stringResource(R.string.timezone),
                    style = MaterialTheme.typography.labelLarge,
                    color = ForegroundMuted,
                )
                Spacer(Modifier.height(4.dp))
                Text(
                    text = regTimezone,
                    style = MaterialTheme.typography.bodyLarge,
                    color = AccentOrange,
                    modifier = Modifier
                        .fillMaxWidth()
                        .clickable { showTimezoneSheet = true }
                        .padding(vertical = 12.dp),
                )
                Spacer(Modifier.height(20.dp))

                ChronosButton(
                    text = stringResource(R.string.create_account),
                    onClick = {
                        localError = when {
                            regUsername.isBlank() -> usernameRequiredError
                            !regEmail.contains('@') -> invalidEmailError
                            regPassword.isBlank() -> passwordRequiredError
                            regPassword != regPasswordConfirm -> mismatchError
                            else -> null
                        }
                        if (localError == null) {
                            onRegister(regEmail, regUsername, regPassword, regTimezone)
                        }
                    },
                    modifier = Modifier.fillMaxWidth(),
                    loading = uiState.loading,
                )
            }

            Spacer(Modifier.height(32.dp))
        }

        ErrorBanner(
            message = when {
                selectedTab == AuthTab.LOGIN && uiState.resendState == ResendState.Sent -> null
                selectedTab == AuthTab.LOGIN -> uiState.error
                else -> localError ?: uiState.error
            },
            modifier = Modifier.align(Alignment.BottomCenter),
            onDismiss = {
                localError = null
                onClearError()
            },
        )
    }

    if (showTimezoneSheet) {
        ModalBottomSheet(
            onDismissRequest = { showTimezoneSheet = false },
            containerColor = BackgroundCard,
        ) {
            TimezonePickerSheetContent(
                timezones = uiState.timezones,
                onSelect = {
                    regTimezone = it
                    showTimezoneSheet = false
                },
            )
        }
    }
}

@Composable
private fun AuthTabItem(
    label: String,
    isSelected: Boolean,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Column(
        modifier = modifier
            .clickable(onClick = onClick)
            .padding(vertical = 12.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        Text(
            text = label,
            style = MaterialTheme.typography.titleSmall,
            color = if (isSelected) AccentOrange else ForegroundMuted,
            fontWeight = if (isSelected) FontWeight.SemiBold else FontWeight.Normal,
        )
        Spacer(Modifier.height(4.dp))
        Box(
            modifier = Modifier
                .width(32.dp)
                .height(2.dp)
                .background(
                    color = if (isSelected) AccentOrange else Color.Transparent,
                    shape = RoundedCornerShape(1.dp),
                ),
        )
    }
}
