package com.chronos.reminder.reminders.ui.list

import androidx.compose.ui.test.junit4.createComposeRule
import androidx.compose.ui.test.onNodeWithContentDescription
import androidx.compose.ui.test.onNodeWithText
import androidx.compose.ui.test.performClick
import com.chronos.reminder.core.ui.theme.ChronosTheme
import com.chronos.reminder.reminders.domain.Destination
import com.chronos.reminder.reminders.domain.RecurrenceType
import com.chronos.reminder.reminders.domain.Reminder
import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.JsonPrimitive
import org.junit.Assert.assertTrue
import org.junit.Rule
import org.junit.Test
import java.time.Instant

class ReminderCardTest {

    @get:Rule
    val composeRule = createComposeRule()

    private val reminder = Reminder(
        id = "r1",
        message = "Water the plants",
        remindAtUtc = Instant.now().plusSeconds(3600),
        nextFireUtc = Instant.now().plusSeconds(3600),
        recurrence = RecurrenceType.DAILY,
        isPaused = false,
        destinations = listOf(
            Destination(
                Destination.TYPE_DISCORD_DM,
                JsonObject(mapOf("user_id" to JsonPrimitive("42"))),
            ),
        ),
        createdAt = Instant.now(),
    )

    @Test
    fun rendersMessageAndRecurrenceBadge() {
        composeRule.setContent {
            ChronosTheme {
                ReminderCard(
                    reminder = reminder,
                    userTimezone = "UTC",
                    onClick = {},
                    onEdit = {},
                    onDuplicate = {},
                    onTogglePause = {},
                    onDelete = {},
                )
            }
        }

        composeRule.onNodeWithText("Water the plants").assertExists()
        composeRule.onNodeWithText("DAILY").assertExists()
    }

    @Test
    fun pausedReminder_showsPausedBadgeAndResumeAction() {
        composeRule.setContent {
            ChronosTheme {
                ReminderCard(
                    reminder = reminder.copy(isPaused = true),
                    userTimezone = "UTC",
                    onClick = {},
                    onEdit = {},
                    onDuplicate = {},
                    onTogglePause = {},
                    onDelete = {},
                )
            }
        }

        composeRule.onNodeWithText("PAUSED").assertExists()
        composeRule.onNodeWithContentDescription("More options").performClick()
        composeRule.onNodeWithText("Resume").assertExists()
    }

    @Test
    fun deleteMenuAction_invokesCallback() {
        var deleteClicked = false
        composeRule.setContent {
            ChronosTheme {
                ReminderCard(
                    reminder = reminder,
                    userTimezone = "UTC",
                    onClick = {},
                    onEdit = {},
                    onDuplicate = {},
                    onTogglePause = {},
                    onDelete = { deleteClicked = true },
                )
            }
        }

        composeRule.onNodeWithContentDescription("More options").performClick()
        composeRule.onNodeWithText("Delete").performClick()
        assertTrue(deleteClicked)
    }
}
