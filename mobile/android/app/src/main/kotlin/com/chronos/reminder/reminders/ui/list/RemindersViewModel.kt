package com.chronos.reminder.reminders.ui.list

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.chronos.reminder.account.data.AccountRepository
import com.chronos.reminder.core.network.ApiResult
import com.chronos.reminder.reminders.data.RemindersRepository
import com.chronos.reminder.reminders.domain.Reminder
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.combine
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

enum class ReminderFilter { ALL, ACTIVE, PAUSED }

data class RemindersListState(
    val refreshing: Boolean = false,
    val error: String? = null,
    val filter: ReminderFilter = ReminderFilter.ALL,
    val userTimezone: String = java.util.TimeZone.getDefault().id,
)

@HiltViewModel
class RemindersViewModel @Inject constructor(
    private val repository: RemindersRepository,
    private val accountRepository: AccountRepository,
) : ViewModel() {

    private val _state = MutableStateFlow(RemindersListState())
    val state: StateFlow<RemindersListState> = _state.asStateFlow()

    val reminders: StateFlow<List<Reminder>> =
        combine(repository.getReminders(), _state) { reminders, state ->
            when (state.filter) {
                ReminderFilter.ALL -> reminders
                ReminderFilter.ACTIVE -> reminders.filter { !it.isPaused }
                ReminderFilter.PAUSED -> reminders.filter { it.isPaused }
            }
        }.stateIn(viewModelScope, SharingStarted.WhileSubscribed(5_000), emptyList())

    init {
        refresh()
        viewModelScope.launch {
            accountRepository.refreshAccount()
            _state.update { it.copy(userTimezone = accountRepository.userTimezone) }
        }
    }

    fun refresh() {
        viewModelScope.launch {
            _state.update { it.copy(refreshing = true) }
            val result = repository.refreshReminders()
            _state.update { it.copy(refreshing = false, error = result.errorMessage()) }
        }
    }

    fun setFilter(filter: ReminderFilter) {
        _state.update { it.copy(filter = filter) }
    }

    fun delete(id: String) = runOp { repository.deleteReminder(id) }
    fun pause(id: String) = runOp { repository.pauseReminder(id) }
    fun resume(id: String) = runOp { repository.resumeReminder(id) }
    fun duplicate(id: String) = runOp { repository.duplicateReminder(id).map { } }

    fun clearError() = _state.update { it.copy(error = null) }

    private fun runOp(op: suspend () -> ApiResult<Unit>) {
        viewModelScope.launch {
            val result = op()
            _state.update { it.copy(error = result.errorMessage()) }
        }
    }

    private fun <T> ApiResult<T>.errorMessage(): String? = when (this) {
        is ApiResult.Success -> null
        is ApiResult.Error -> message
        is ApiResult.NetworkError -> "No internet connection"
    }
}
