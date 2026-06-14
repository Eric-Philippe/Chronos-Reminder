package com.chronos.reminder.auth.ui

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.chronos.reminder.account.data.AccountRepository
import com.chronos.reminder.account.data.TimezoneDto
import com.chronos.reminder.auth.data.AuthRepository
import com.chronos.reminder.auth.domain.LoginUseCase
import com.chronos.reminder.core.network.ApiResult
import com.chronos.reminder.core.network.AuthEvent
import com.chronos.reminder.core.network.AuthEventBus
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

data class AuthUiState(
    val isLoggedIn: Boolean = false,
    val loading: Boolean = false,
    val error: String? = null,
    val registerSuccess: Boolean = false,
    val resetEmailSent: Boolean = false,
    val resendState: ResendState = ResendState.Idle,
    val timezones: List<TimezoneDto> = emptyList(),
)

enum class ResendState { Idle, Sending, Sent, Error }

@HiltViewModel
class AuthViewModel @Inject constructor(
    private val authRepository: AuthRepository,
    private val accountRepository: AccountRepository,
    private val loginUseCase: LoginUseCase,
    authEventBus: AuthEventBus,
) : ViewModel() {

    private val _uiState = MutableStateFlow(AuthUiState(isLoggedIn = authRepository.isLoggedIn()))
    val uiState: StateFlow<AuthUiState> = _uiState.asStateFlow()

    init {
        viewModelScope.launch {
            authEventBus.events.collect { event ->
                if (event is AuthEvent.LoggedOut) {
                    authRepository.clearLocalData()
                    accountRepository.clear()
                    _uiState.update { it.copy(isLoggedIn = false) }
                }
            }
        }
    }

    fun loginWithEmail(email: String, password: String) {
        runAuthOp(
            op = { loginUseCase(email, password) },
            onSuccess = { state -> state.copy(isLoggedIn = true) },
        )
    }

    fun register(email: String, username: String, password: String, timezone: String) {
        runAuthOp(
            op = { authRepository.register(email.trim(), username.trim(), password, timezone) },
            onSuccess = { state -> state.copy(registerSuccess = true) },
        )
    }

    fun handleDiscordCode(code: String) {
        runAuthOp(
            op = { authRepository.loginWithDiscordCode(code) },
            onSuccess = { state -> state.copy(isLoggedIn = true) },
        )
    }

    fun resendVerification(email: String) {
        viewModelScope.launch {
            _uiState.update { it.copy(resendState = ResendState.Sending) }
            when (authRepository.resendVerification(email)) {
                is ApiResult.Success -> _uiState.update { it.copy(resendState = ResendState.Sent) }
                else -> _uiState.update { it.copy(resendState = ResendState.Error) }
            }
        }
    }

    fun forgotPassword(email: String) {
        runAuthOp(
            op = { authRepository.requestPasswordReset(email.trim()) },
            onSuccess = { state -> state.copy(resetEmailSent = true) },
        )
    }

    fun logout() {
        viewModelScope.launch {
            _uiState.update { it.copy(loading = true) }
            authRepository.logout()
            accountRepository.clear()
            _uiState.update { AuthUiState(isLoggedIn = false) }
        }
    }

    fun loadTimezones() {
        if (_uiState.value.timezones.isNotEmpty()) return
        viewModelScope.launch {
            val result = accountRepository.getTimezones()
            if (result is ApiResult.Success) {
                _uiState.update { it.copy(timezones = result.data) }
            }
        }
    }

    fun clearError() {
        _uiState.update { it.copy(error = null) }
    }

    fun setError(message: String) {
        _uiState.update { it.copy(error = message) }
    }

    fun clearFlags() {
        _uiState.update { it.copy(registerSuccess = false, resetEmailSent = false, error = null, resendState = ResendState.Idle) }
    }

    private fun runAuthOp(
        op: suspend () -> ApiResult<Unit>,
        onSuccess: (AuthUiState) -> AuthUiState,
    ) {
        viewModelScope.launch {
            _uiState.update { it.copy(loading = true, error = null) }
            when (val result = op()) {
                is ApiResult.Success -> _uiState.update { onSuccess(it.copy(loading = false)) }
                is ApiResult.Error -> _uiState.update { it.copy(loading = false, error = result.message) }
                is ApiResult.NetworkError -> _uiState.update {
                    it.copy(loading = false, error = "No internet connection")
                }
            }
        }
    }
}
