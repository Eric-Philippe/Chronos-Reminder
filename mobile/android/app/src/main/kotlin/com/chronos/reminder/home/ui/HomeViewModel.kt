package com.chronos.reminder.home.ui

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.chronos.reminder.account.data.AccountDto
import com.chronos.reminder.account.data.AccountRepository
import com.chronos.reminder.dfm.data.DfmRepository
import com.chronos.reminder.reminders.data.RemindersRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

data class HomeUiState(
    val loading: Boolean = false,
    val account: AccountDto? = null,
    val totalReminders: Int = 0,
    val activeReminders: Int = 0,
    val dfmItemCount: Int = 0,
)

@HiltViewModel
class HomeViewModel @Inject constructor(
    private val accountRepository: AccountRepository,
    private val remindersRepository: RemindersRepository,
    private val dfmRepository: DfmRepository,
) : ViewModel() {

    private val _state = MutableStateFlow(HomeUiState())
    val state: StateFlow<HomeUiState> = _state.asStateFlow()

    init {
        viewModelScope.launch {
            accountRepository.account.collect { account ->
                _state.update { it.copy(account = account) }
            }
        }
        viewModelScope.launch {
            remindersRepository.getReminders().collect { reminders ->
                _state.update {
                    it.copy(
                        totalReminders = reminders.size,
                        activeReminders = reminders.count { r -> !r.isPaused },
                    )
                }
            }
        }
        viewModelScope.launch {
            dfmRepository.getItems().collect { items ->
                _state.update { it.copy(dfmItemCount = items.size) }
            }
        }
        viewModelScope.launch {
            _state.update { it.copy(loading = true) }
            accountRepository.refreshAccount()
            _state.update { it.copy(loading = false) }
        }
    }
}
