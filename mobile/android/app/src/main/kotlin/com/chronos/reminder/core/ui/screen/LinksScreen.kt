package com.chronos.reminder.core.ui.screen

import android.content.Intent
import android.net.Uri
import androidx.compose.foundation.clickable
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
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.Article
import androidx.compose.material.icons.automirrored.filled.HelpOutline
import androidx.compose.material.icons.automirrored.filled.KeyboardArrowRight
import androidx.compose.material.icons.filled.Code
import androidx.compose.material.icons.filled.Forum
import androidx.compose.material.icons.filled.Language
import androidx.compose.material.icons.filled.Policy
import androidx.compose.material.icons.filled.Speed
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.components.ChronosCard
import com.chronos.reminder.core.ui.components.ChronosTopBar
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.BackgroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMuted

private const val WEB_APP_URL = com.chronos.reminder.core.AppConstants.Urls.WEB_APP
private const val DOCS_URL = com.chronos.reminder.core.AppConstants.Urls.DOCS
private const val DISCORD_URL = com.chronos.reminder.core.AppConstants.Urls.DISCORD_INVITE
private const val GITHUB_URL = com.chronos.reminder.core.AppConstants.Urls.GITHUB
private const val STATUS_URL = com.chronos.reminder.core.AppConstants.Urls.STATUS
private const val CONTACT_URL = com.chronos.reminder.core.AppConstants.Urls.CONTACT

@Composable
fun LinksScreen(
    onBack: () -> Unit,
    onOpenTerms: () -> Unit,
    onOpenPrivacy: () -> Unit,
) {
    val context = LocalContext.current

    fun openUrl(url: String) {
        context.startActivity(Intent(Intent.ACTION_VIEW, Uri.parse(url)))
    }

    Scaffold(
        containerColor = BackgroundMain,
        topBar = { ChronosTopBar(title = stringResource(R.string.links_screen_title), onBack = onBack) },
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(padding)
                .padding(horizontal = 16.dp),
        ) {
            SectionLabel(stringResource(R.string.section_links))

            ChronosCard(modifier = Modifier.fillMaxWidth()) {
                Column(Modifier.padding(vertical = 4.dp)) {
                    LinkRow(
                        icon = Icons.Default.Language,
                        label = stringResource(R.string.link_web_app),
                        onClick = { openUrl(WEB_APP_URL) },
                    )
                    LinkRow(
                        icon = Icons.AutoMirrored.Filled.Article,
                        label = stringResource(R.string.link_documentation),
                        onClick = { openUrl(DOCS_URL) },
                    )
                    LinkRow(
                        icon = Icons.Default.Forum,
                        label = stringResource(R.string.link_discord),
                        onClick = { openUrl(DISCORD_URL) },
                    )
                    LinkRow(
                        icon = Icons.Default.Code,
                        label = stringResource(R.string.link_github),
                        onClick = { openUrl(GITHUB_URL) },
                    )
                    LinkRow(
                        icon = Icons.AutoMirrored.Filled.HelpOutline,
                        label = stringResource(R.string.link_support),
                        onClick = { openUrl(CONTACT_URL) },
                    )
                    LinkRow(
                        icon = Icons.Default.Speed,
                        label = stringResource(R.string.link_status),
                        onClick = { openUrl(STATUS_URL) },
                        showDivider = false,
                    )
                }
            }

            SectionLabel(stringResource(R.string.section_legal))

            ChronosCard(modifier = Modifier.fillMaxWidth()) {
                Column(Modifier.padding(vertical = 4.dp)) {
                    LinkRow(
                        icon = Icons.AutoMirrored.Filled.Article,
                        label = stringResource(R.string.terms_of_use),
                        onClick = onOpenTerms,
                    )
                    LinkRow(
                        icon = Icons.Default.Policy,
                        label = stringResource(R.string.privacy_policy),
                        onClick = onOpenPrivacy,
                        showDivider = false,
                    )
                }
            }

            Spacer(Modifier.height(32.dp))
        }
    }
}

@Composable
private fun SectionLabel(text: String) {
    Spacer(Modifier.height(20.dp))
    Text(
        text = text,
        style = MaterialTheme.typography.titleMedium,
        color = AccentOrange,
    )
    Spacer(Modifier.height(8.dp))
}

@Composable
private fun LinkRow(
    icon: ImageVector,
    label: String,
    onClick: () -> Unit,
    showDivider: Boolean = true,
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 14.dp),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        Icon(
            imageVector = icon,
            contentDescription = null,
            tint = AccentOrange,
            modifier = Modifier.size(20.dp),
        )
        Spacer(Modifier.width(12.dp))
        Text(
            text = label,
            style = MaterialTheme.typography.bodyLarge,
            color = ForegroundMain,
            modifier = Modifier.weight(1f),
        )
        Icon(
            Icons.AutoMirrored.Filled.KeyboardArrowRight,
            contentDescription = null,
            tint = ForegroundMuted,
            modifier = Modifier.size(18.dp),
        )
    }
}
