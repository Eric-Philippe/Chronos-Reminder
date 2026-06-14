package com.chronos.reminder.reminders.data

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class GuildsRequest(@SerialName("account_id") val accountId: String)

@Serializable
data class GuildsResponse(val guilds: List<GuildDto> = emptyList())

@Serializable
data class GuildDto(
    val id: String,
    val name: String,
    val icon: String? = null,
)

@Serializable
data class GuildChannelsRequest(
    @SerialName("account_id") val accountId: String,
    @SerialName("guild_id") val guildId: String,
)

@Serializable
data class GuildChannelsResponse(val channels: List<ChannelDto> = emptyList())

@Serializable
data class ChannelDto(
    val id: String,
    val name: String,
    val type: Int = 0,
    val position: Int = 0,
)
