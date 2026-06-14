package com.chronos.reminder.notifications

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class FcmTokenRequest(
    val token: String,
    @SerialName("device_id") val deviceId: String,
)
