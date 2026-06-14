package com.chronos.reminder.core.ui.screen

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
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Android
import androidx.compose.material.icons.filled.Build
import androidx.compose.material.icons.filled.Cloud
import androidx.compose.material.icons.filled.Code
import androidx.compose.material.icons.filled.Info
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import com.chronos.reminder.BuildConfig
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.components.ChronosCard
import com.chronos.reminder.core.ui.components.ChronosTopBar
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.BackgroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMuted

@Composable
fun AboutScreen(
    onBack: () -> Unit,
    viewModel: AboutViewModel = hiltViewModel(),
) {
    val apiVersion by viewModel.apiVersion.collectAsStateWithLifecycle()

    Scaffold(
        containerColor = BackgroundMain,
        topBar = { ChronosTopBar(title = stringResource(R.string.about_title), onBack = onBack) },
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(padding)
                .padding(horizontal = 16.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
        ) {
            Spacer(Modifier.height(32.dp))

            // App identity block
            Image(
                painter = painterResource(R.drawable.logo_chronos),
                contentDescription = stringResource(R.string.chronos_logo),
                modifier = Modifier
                    .size(88.dp)
                    .clip(RoundedCornerShape(20.dp)),
            )
            Spacer(Modifier.height(16.dp))
            Text(
                text = stringResource(R.string.app_name),
                style = MaterialTheme.typography.headlineMedium,
                fontWeight = FontWeight.Bold,
                color = ForegroundMain,
            )
            Spacer(Modifier.height(4.dp))
            Text(
                text = stringResource(R.string.about_tagline),
                style = MaterialTheme.typography.bodyMedium,
                color = ForegroundMuted,
            )

            Spacer(Modifier.height(32.dp))

            // Version info card
            ChronosCard(modifier = Modifier.fillMaxWidth()) {
                Column(modifier = Modifier.padding(4.dp)) {
                    AboutRow(
                        icon = Icons.Default.Info,
                        label = stringResource(R.string.about_version),
                        value = BuildConfig.VERSION_NAME,
                    )
                    HorizontalDivider(modifier = Modifier.padding(horizontal = 16.dp))
                    AboutRow(
                        icon = Icons.Default.Build,
                        label = stringResource(R.string.about_build),
                        value = BuildConfig.VERSION_CODE.toString(),
                    )
                    HorizontalDivider(modifier = Modifier.padding(horizontal = 16.dp))
                    AboutRow(
                        icon = Icons.Default.Android,
                        label = stringResource(R.string.about_platform),
                        value = stringResource(R.string.about_platform_value),
                        showDivider = apiVersion != null,
                    )
                    if (apiVersion != null) {
                        HorizontalDivider(modifier = Modifier.padding(horizontal = 16.dp))
                        AboutRow(
                            icon = Icons.Default.Cloud,
                            label = stringResource(R.string.about_api_version),
                            value = apiVersion!!,
                            showDivider = false,
                        )
                    }
                }
            }

            Spacer(Modifier.height(20.dp))

            // Open source card
            ChronosCard(modifier = Modifier.fillMaxWidth()) {
                Column(modifier = Modifier.padding(4.dp)) {
                    AboutRow(
                        icon = Icons.Default.Code,
                        label = stringResource(R.string.about_license),
                        value = stringResource(R.string.about_license_value),
                        showDivider = false,
                    )
                }
            }

            Spacer(Modifier.height(32.dp))

            Text(
                text = stringResource(R.string.about_copyright),
                style = MaterialTheme.typography.bodySmall,
                color = ForegroundMuted,
            )

            Spacer(Modifier.height(32.dp))
        }
    }
}

@Composable
private fun AboutRow(
    icon: ImageVector,
    label: String,
    value: String,
    showDivider: Boolean = true,
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp, vertical = 14.dp),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(12.dp),
    ) {
        Icon(
            imageVector = icon,
            contentDescription = null,
            tint = AccentOrange,
            modifier = Modifier.size(20.dp),
        )
        Text(
            text = label,
            style = MaterialTheme.typography.bodyLarge,
            color = ForegroundMain,
            modifier = Modifier.weight(1f),
        )
        Text(
            text = value,
            style = MaterialTheme.typography.bodyMedium,
            color = ForegroundMuted,
        )
    }
}
