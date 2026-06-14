package com.chronos.reminder.core.ui.screen

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.chronos.reminder.account.data.AccountApi
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class AboutViewModel @Inject constructor(
    private val accountApi: AccountApi,
) : ViewModel() {

    private val _apiVersion = MutableStateFlow<String?>(null)
    val apiVersion: StateFlow<String?> = _apiVersion.asStateFlow()

    init {
        viewModelScope.launch {
            try {
                val response = accountApi.getHealth()
                if (response.isSuccessful) {
                    _apiVersion.value = response.body()?.version
                }
            } catch (_: Exception) {
                // non-fatal: version stays null and we simply omit the row
            }
        }
    }
}
