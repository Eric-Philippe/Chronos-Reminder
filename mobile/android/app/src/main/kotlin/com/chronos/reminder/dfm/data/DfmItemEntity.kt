package com.chronos.reminder.dfm.data

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "dfm_items")
data class DfmItemEntity(
    @PrimaryKey val id: String,
    val content: String,
    val checked: Boolean,
    val position: Int,
    val createdAt: Long,
)
