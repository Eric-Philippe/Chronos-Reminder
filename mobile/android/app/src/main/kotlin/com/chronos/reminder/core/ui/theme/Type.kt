package com.chronos.reminder.core.ui.theme

import androidx.compose.material3.Typography
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.Font
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp
import com.chronos.reminder.R

val Nasa21 = FontFamily(Font(R.font.nasa21))

val ChronosTypography = Typography(
    displayLarge = TextStyle(fontFamily = Nasa21, fontSize = 57.sp),
    displayMedium = TextStyle(fontFamily = Nasa21, fontSize = 45.sp),
    headlineLarge = TextStyle(fontFamily = Nasa21, fontSize = 32.sp),
    headlineMedium = TextStyle(fontFamily = Nasa21, fontSize = 28.sp),
    titleLarge = TextStyle(fontFamily = Nasa21, fontSize = 22.sp),
    titleMedium = TextStyle(fontFamily = Nasa21, fontSize = 16.sp, fontWeight = FontWeight.Medium),
    bodyLarge = TextStyle(fontFamily = Nasa21, fontSize = 16.sp),
    bodyMedium = TextStyle(fontFamily = Nasa21, fontSize = 14.sp),
    labelLarge = TextStyle(fontFamily = Nasa21, fontSize = 14.sp, fontWeight = FontWeight.Medium),
    labelSmall = TextStyle(fontFamily = Nasa21, fontSize = 11.sp),
)
