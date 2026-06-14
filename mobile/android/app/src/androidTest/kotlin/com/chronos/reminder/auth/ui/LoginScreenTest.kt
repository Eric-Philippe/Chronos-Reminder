package com.chronos.reminder.auth.ui

import androidx.compose.ui.test.junit4.createComposeRule
import androidx.compose.ui.test.onNodeWithText
import androidx.compose.ui.test.performClick
import androidx.compose.ui.test.performTextInput
import com.chronos.reminder.core.ui.theme.ChronosTheme
import org.junit.Assert.assertEquals
import org.junit.Rule
import org.junit.Test

class LoginScreenTest {

    @get:Rule
    val composeRule = createComposeRule()

    @Test
    fun fillFormAndSubmit_invokesLoginWithEnteredValues() {
        var capturedEmail: String? = null
        var capturedPassword: String? = null

        composeRule.setContent {
            ChronosTheme {
                LoginScreen(
                    uiState = AuthUiState(),
                    onLogin = { email, password ->
                        capturedEmail = email
                        capturedPassword = password
                    },
                    onNavigateRegister = {},
                    onNavigateForgotPassword = {},
                    onClearError = {},
                    onDiscordUnconfigured = {},
                )
            }
        }

        composeRule.onNodeWithText("Email").performTextInput("user@example.com")
        composeRule.onNodeWithText("Password").performTextInput("secret123")
        composeRule.onNodeWithText("Sign In").performClick()

        assertEquals("user@example.com", capturedEmail)
        assertEquals("secret123", capturedPassword)
    }

    @Test
    fun navigationLinks_triggerCallbacks() {
        var registerClicked = false
        var forgotClicked = false

        composeRule.setContent {
            ChronosTheme {
                LoginScreen(
                    uiState = AuthUiState(),
                    onLogin = { _, _ -> },
                    onNavigateRegister = { registerClicked = true },
                    onNavigateForgotPassword = { forgotClicked = true },
                    onClearError = {},
                    onDiscordUnconfigured = {},
                )
            }
        }

        composeRule.onNodeWithText("Don't have an account? Sign up").performClick()
        composeRule.onNodeWithText("Forgot password?").performClick()

        assertEquals(true, registerClicked)
        assertEquals(true, forgotClicked)
    }

    @Test
    fun apiError_isShownInBanner() {
        composeRule.setContent {
            ChronosTheme {
                LoginScreen(
                    uiState = AuthUiState(error = "Invalid credentials"),
                    onLogin = { _, _ -> },
                    onNavigateRegister = {},
                    onNavigateForgotPassword = {},
                    onClearError = {},
                    onDiscordUnconfigured = {},
                )
            }
        }

        composeRule.onNodeWithText("Invalid credentials").assertExists()
    }
}
