package com.chronos.reminder.core.ui.screen

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import com.chronos.reminder.R
import com.chronos.reminder.core.ui.components.ChronosCard
import com.chronos.reminder.core.ui.components.ChronosTopBar
import com.chronos.reminder.core.ui.theme.AccentOrange
import com.chronos.reminder.core.ui.theme.BackgroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMain
import com.chronos.reminder.core.ui.theme.ForegroundMuted

@Composable
fun TermsScreen(onBack: () -> Unit) {
    Scaffold(
        containerColor = BackgroundMain,
        topBar = { ChronosTopBar(title = stringResource(R.string.terms_title), onBack = onBack) },
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(padding)
                .padding(horizontal = 16.dp),
        ) {
            Spacer(Modifier.height(16.dp))

            Text(
                text = stringResource(R.string.terms_subtitle),
                style = MaterialTheme.typography.bodyMedium,
                color = ForegroundMuted,
            )
            Spacer(Modifier.height(4.dp))
            Text(
                text = stringResource(R.string.terms_updated),
                style = MaterialTheme.typography.labelSmall,
                color = ForegroundMuted,
            )

            Spacer(Modifier.height(20.dp))

            val sections = listOf(
                R.string.terms_s1_title to R.string.terms_s1_body,
                R.string.terms_s2_title to R.string.terms_s2_body,
                R.string.terms_s3_title to R.string.terms_s3_body,
                R.string.terms_s4_title to R.string.terms_s4_body,
                R.string.terms_s5_title to R.string.terms_s5_body,
                R.string.terms_s6_title to R.string.terms_s6_body,
                R.string.terms_s7_title to R.string.terms_s7_body,
                R.string.terms_s8_title to R.string.terms_s8_body,
                R.string.terms_s9_title to R.string.terms_s9_body,
            )

            sections.forEach { (titleRes, bodyRes) ->
                ChronosCard(modifier = Modifier.fillMaxWidth()) {
                    Column(Modifier.padding(16.dp)) {
                        Text(
                            text = stringResource(titleRes),
                            style = MaterialTheme.typography.titleMedium,
                            color = AccentOrange,
                        )
                        Spacer(Modifier.height(8.dp))
                        Text(
                            text = stringResource(bodyRes),
                            style = MaterialTheme.typography.bodyMedium,
                            color = ForegroundMain,
                        )
                    }
                }
                Spacer(Modifier.height(10.dp))
            }

            Spacer(Modifier.height(32.dp))
        }
    }
}
