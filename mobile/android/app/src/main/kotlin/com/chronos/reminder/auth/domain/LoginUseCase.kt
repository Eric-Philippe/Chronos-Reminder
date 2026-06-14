package com.chronos.reminder.auth.domain

import com.chronos.reminder.auth.data.AuthRepository
import com.chronos.reminder.core.network.ApiResult
import javax.inject.Inject

class LoginUseCase @Inject constructor(
    private val authRepository: AuthRepository,
) {
    suspend operator fun invoke(email: String, password: String): ApiResult<Unit> {
        if (email.isBlank() || !email.contains('@')) return ApiResult.Error(-1, "Invalid email address")
        if (password.isBlank()) return ApiResult.Error(-1, "Password is required")
        return authRepository.login(email.trim(), password)
    }
}
