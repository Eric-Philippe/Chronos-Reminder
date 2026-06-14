package com.chronos.reminder.core.network

import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import retrofit2.Response
import java.io.IOException

sealed class ApiResult<out T> {
    data class Success<T>(val data: T) : ApiResult<T>()
    data class Error(val code: Int, val message: String) : ApiResult<Nothing>()
    data object NetworkError : ApiResult<Nothing>()

    inline fun <R> map(transform: (T) -> R): ApiResult<R> = when (this) {
        is Success -> Success(transform(data))
        is Error -> this
        is NetworkError -> this
    }

    inline fun onSuccess(block: (T) -> Unit): ApiResult<T> {
        if (this is Success) block(data)
        return this
    }
}

@Serializable
data class ApiErrorBody(val error: String? = null, val message: String? = null)

private val errorJson = Json { ignoreUnknownKeys = true; isLenient = true }

// Single bridge between Retrofit and the rest of the app: nothing throws across
// layer boundaries, IO failures collapse into NetworkError.
suspend fun <T> safeApiCall(call: suspend () -> Response<T>): ApiResult<T> = try {
    val response = call()
    val body = response.body()
    when {
        response.isSuccessful && body != null -> ApiResult.Success(body)
        response.isSuccessful -> ApiResult.Error(response.code(), "Empty response body")
        else -> {
            val raw = response.errorBody()?.string().orEmpty()
            val parsed = runCatching { errorJson.decodeFromString<ApiErrorBody>(raw) }.getOrNull()
            ApiResult.Error(response.code(), parsed?.error ?: parsed?.message ?: "Request failed (${response.code()})")
        }
    }
} catch (e: IOException) {
    ApiResult.NetworkError
} catch (e: Exception) {
    ApiResult.Error(-1, e.message ?: "Unexpected error")
}
