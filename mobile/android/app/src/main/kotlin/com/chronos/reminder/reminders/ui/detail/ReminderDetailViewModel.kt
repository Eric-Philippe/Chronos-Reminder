package com.chronos.reminder.reminders.ui.detail

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.toRoute
import com.chronos.reminder.account.data.AccountRepository
import com.chronos.reminder.core.navigation.Screen
import com.chronos.reminder.core.network.ApiResult
import com.chronos.reminder.reminders.data.RemindersRepository
import com.chronos.reminder.reminders.domain.Reminder
import com.chronos.reminder.reminders.ui.create.ReminderForm
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

data class ReminderDetailState(
    val editing: Boolean = false,
    val form: ReminderForm = ReminderForm(),
    val saving: Boolean = false,
    val deleted: Boolean = false,
    val error: String? = null,
    val userTimezone: String = java.util.TimeZone.getDefault().id,
)

@HiltViewModel
class ReminderDetailViewModel @Inject constructor(
    savedStateHandle: SavedStateHandle,
    private val repository: RemindersRepository,
    private val accountRepository: AccountRepository,
) : ViewModel() {

    private val reminderId: String = savedStateHandle.toRoute<Screen.ReminderDetail>().id

    private val _state = MutableStateFlow(ReminderDetailState())
    val state: StateFlow<ReminderDetailState> = _state.asStateFlow()

    val reminder: StateFlow<Reminder?> = repository.getReminder(reminderId)
        .stateIn(viewModelScope, SharingStarted.WhileSubscribed(5_000), null)

    init {
        viewModelScope.launch {
            repository.refreshReminder(reminderId)
            if (accountRepository.account.value == null) {
                accountRepository.refreshAccount()
            }
            _state.update { it.copy(userTimezone = accountRepository.userTimezone) }
        }
    }

    fun startEditing() {
        val current = reminder.value ?: return
        _state.update {
            it.copy(editing = true, form = ReminderForm.fromReminder(current, it.userTimezone))
        }
    }

    fun cancelEditing() {
        _state.update { it.copy(editing = false, error = null) }
    }

    fun updateForm(form: ReminderForm) {
        _state.update { it.copy(form = form) }
    }

    fun save() {
        val form = _state.value.form
        if (form.date == null || form.time == null || form.message.isBlank()) return
        viewModelScope.launch {
            _state.update { it.copy(saving = true, error = null) }
            when (val result = repository.updateReminder(reminderId, form.toRequest())) {
                is ApiResult.Success -> _state.update { it.copy(saving = false, editing = false) }
                is ApiResult.Error -> _state.update { it.copy(saving = false, error = result.message) }
                is ApiResult.NetworkError -> _state.update {
                    it.copy(saving = false, error = "No internet connection")
                }
            }
        }
    }

    fun togglePause() {
        val current = reminder.value ?: return
        viewModelScope.launch {
            val result = if (current.isPaused) {
                repository.resumeReminder(reminderId)
            } else {
                repository.pauseReminder(reminderId)
            }
            if (result is ApiResult.Error) {
                _state.update { it.copy(error = result.message) }
            }
        }
    }

    fun delete() {
        viewModelScope.launch {
            when (val result = repository.deleteReminder(reminderId)) {
                is ApiResult.Success -> _state.update { it.copy(deleted = true) }
                is ApiResult.Error -> _state.update { it.copy(error = result.message) }
                is ApiResult.NetworkError -> _state.update { it.copy(error = "No internet connection") }
            }
        }
    }

    fun clearError() = _state.update { it.copy(error = null) }
}
