package com.chronos.reminder

import android.Manifest
import android.content.Context
import android.content.Intent
import android.content.pm.PackageManager
import android.content.res.Configuration
import android.net.Uri
import android.os.Build
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.result.contract.ActivityResultContracts
import androidx.activity.viewModels
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.core.content.ContextCompat
import androidx.core.splashscreen.SplashScreen.Companion.installSplashScreen
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.navigation.compose.rememberNavController
import com.chronos.reminder.auth.ui.AuthViewModel
import com.chronos.reminder.core.navigation.ChronosNavGraph
import com.chronos.reminder.core.ui.theme.ChronosTheme
import com.chronos.reminder.core.util.ConnectivityObserver
import dagger.hilt.android.AndroidEntryPoint
import javax.inject.Inject

@AndroidEntryPoint
class MainActivity : ComponentActivity() {

    @Inject
    lateinit var connectivityObserver: ConnectivityObserver

    private val authViewModel: AuthViewModel by viewModels()

    private var pendingReminderId by mutableStateOf<String?>(null)
    private var pendingDiscordLinkCode by mutableStateOf<String?>(null)

    private val notificationPermissionLauncher =
        registerForActivityResult(ActivityResultContracts.RequestPermission()) { /* best effort */ }

    override fun attachBaseContext(newBase: Context) {
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.TIRAMISU) {
            val prefs = newBase.getSharedPreferences("chronos_prefs", Context.MODE_PRIVATE)
            val tag = prefs.getString("locale_tag", null)
            if (!tag.isNullOrEmpty()) {
                val locale = java.util.Locale.forLanguageTag(tag)
                val config = Configuration(newBase.resources.configuration)
                config.setLocale(locale)
                super.attachBaseContext(newBase.createConfigurationContext(config))
                return
            }
        }
        super.attachBaseContext(newBase)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        installSplashScreen()
        super.onCreate(savedInstanceState)

        handleDeepLink(intent)

        setContent {
            ChronosTheme {
                val navController = rememberNavController()
                val uiState by authViewModel.uiState.collectAsStateWithLifecycle()
                val isOnline by connectivityObserver.isOnline.collectAsStateWithLifecycle(initialValue = true)

                LaunchedEffect(uiState.isLoggedIn) {
                    if (uiState.isLoggedIn) requestNotificationPermission()
                }

                ChronosNavGraph(
                    navController = navController,
                    authViewModel = authViewModel,
                    uiState = uiState,
                    isOnline = isOnline,
                    pendingReminderId = pendingReminderId,
                    onPendingReminderConsumed = { pendingReminderId = null },
                    pendingDiscordLinkCode = pendingDiscordLinkCode,
                    onPendingDiscordLinkConsumed = { pendingDiscordLinkCode = null },
                    onDiscordUnconfigured = {
                        authViewModel.setError(getString(R.string.error_discord_not_configured))
                    },
                )
            }
        }
    }

    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)
        handleDeepLink(intent)
    }

    private fun handleDeepLink(intent: Intent?) {
        val data: Uri = intent?.data ?: run {
            // Notification taps deliver the reminder id as an extra
            intent?.getStringExtra(EXTRA_REMINDER_ID)?.let { pendingReminderId = it }
            return
        }
        when (data.host) {
            "auth" -> data.getQueryParameter("code")?.let { code ->
                // If already authenticated, this OAuth round-trip is a "link Discord
                // to my existing account" flow; otherwise it's a login/signup.
                if (authViewModel.uiState.value.isLoggedIn) {
                    pendingDiscordLinkCode = code
                } else {
                    authViewModel.handleDiscordCode(code)
                }
            }
            "reminder" -> data.getQueryParameter("id")?.let { id ->
                pendingReminderId = id
            }
        }
    }

    private fun requestNotificationPermission() {
        val granted = ContextCompat.checkSelfPermission(
            this,
            Manifest.permission.POST_NOTIFICATIONS,
        ) == PackageManager.PERMISSION_GRANTED
        if (!granted) {
            notificationPermissionLauncher.launch(Manifest.permission.POST_NOTIFICATIONS)
        }
    }

    companion object {
        const val EXTRA_REMINDER_ID = "reminder_id"
    }
}
