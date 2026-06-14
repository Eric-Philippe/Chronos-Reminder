package com.chronos.reminder.reminders.data

import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST

interface DiscordApi {

    @POST("api/discord/guilds")
    suspend fun getGuilds(@Body body: GuildsRequest): Response<GuildsResponse>

    @POST("api/discord/guilds/channels")
    suspend fun getGuildChannels(@Body body: GuildChannelsRequest): Response<GuildChannelsResponse>
}
