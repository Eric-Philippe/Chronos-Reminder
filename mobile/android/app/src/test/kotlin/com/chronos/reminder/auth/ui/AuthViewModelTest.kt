package com.chronos.reminder.auth.ui

import app.cash.turbine.test
import com.chronos.reminder.MainDispatcherRule
import com.chronos.reminder.account.data.AccountRepository
import com.chronos.reminder.auth.data.AuthRepository
import com.chronos.reminder.auth.domain.LoginUseCase
import com.chronos.reminder.core.network.ApiResult
import com.chronos.reminder.core.network.AuthEvent
import com.chronos.reminder.core.network.AuthEventBus
import io.mockk.coEvery
import io.mockk.coJustRun
import io.mockk.every
import io.mockk.justRun
import io.mockk.mockk
import kotlinx.coroutines.test.runTest
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Rule
import org.junit.Test

class AuthViewModelTest {

    @get:Rule
    val mainDispatcherRule = MainDispatcherRule()

    private val authRepository: AuthRepository = mockk()
    private val accountRepository: AccountRepository = mockk(relaxed = true)
    private val loginUseCase: LoginUseCase = mockk()
    private val authEventBus = AuthEventBus()

    private lateinit var viewModel: AuthViewModel

    @Before
    fun setUp() {
        every { authRepository.isLoggedIn() } returns false
        coJustRun { authRepository.clearLocalData() }
        justRun { accountRepository.clear() }
        viewModel = AuthViewModel(authRepository, accountRepository, loginUseCase, authEventBus)
    }

    @Test
    fun `login success flips logged in state`() = runTest {
        coEvery { loginUseCase("user@example.com", "secret") } returns ApiResult.Success(Unit)

        viewModel.uiState.test {
            assertFalse(awaitItem().isLoggedIn)

            viewModel.loginWithEmail("user@example.com", "secret")

            val loggedIn = expectMostRecentItem()
            assertTrue(loggedIn.isLoggedIn)
            assertFalse(loggedIn.loading)
            assertNull(loggedIn.error)
        }
    }

    @Test
    fun `login failure exposes error message`() = runTest {
        coEvery { loginUseCase(any(), any()) } returns ApiResult.Error(401, "Invalid credentials")

        viewModel.loginWithEmail("user@example.com", "wrong")

        viewModel.uiState.test {
            val state = expectMostRecentItem()
            assertFalse(state.isLoggedIn)
            assertEquals("Invalid credentials", state.error)
        }
    }

    @Test
    fun `network failure exposes connectivity error`() = runTest {
        coEvery { loginUseCase(any(), any()) } returns ApiResult.NetworkError

        viewModel.loginWithEmail("user@example.com", "secret")

        viewModel.uiState.test {
            assertEquals("No internet connection", expectMostRecentItem().error)
        }
    }

    @Test
    fun `register success sets flag without logging in`() = runTest {
        coEvery {
            authRepository.register("new@example.com", "newbie", "secret", "Europe/Paris")
        } returns ApiResult.Success(Unit)

        viewModel.register("new@example.com", "newbie", "secret", "Europe/Paris")

        viewModel.uiState.test {
            val state = expectMostRecentItem()
            assertTrue(state.registerSuccess)
            assertFalse(state.isLoggedIn)
        }
    }

    @Test
    fun `401 event from interceptor logs the session out`() = runTest {
        every { authRepository.isLoggedIn() } returns true
        viewModel = AuthViewModel(authRepository, accountRepository, loginUseCase, authEventBus)

        viewModel.uiState.test {
            assertTrue(awaitItem().isLoggedIn)

            authEventBus.emit(AuthEvent.LoggedOut)

            assertFalse(awaitItem().isLoggedIn)
        }
    }
}
