package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/logic"
)

// Register the calcultime command
func init() {
	RegisterCommand(&Command{
		Description: Description{
			Name:             "calcultime",
			Emoji:            "🧮",
			CategoryName:     "Tools",
			ShortDescription: "Calculate time operations",
			FullDescription:  "Perform calculations between times or with factors. Supports addition, subtraction, multiplication, and division of time values.",
			Usage:            "/calcultime time1:<time> operation:<operation> time2:<time/factor>",
			Example:          "/calcultime time1:\"2h 30m\" operation:add time2:\"1h 15m\"",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "calcultime",
			Description: "Calculate time operations (add, subtract, multiply, divide)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time1",
					Description: "First time value (e.g., '2h 30m', '14:30', '2.5h')",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "operation",
					Description: "Operation to perform",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Add (+)",
							Value: "ADD",
						},
						{
							Name:  "Subtract (-)",
							Value: "SUBTRACT",
						},
						{
							Name:  "Multiply (×)",
							Value: "MULTIPLY",
						},
						{
							Name:  "Divide (÷)",
							Value: "DIVIDE",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time2",
					Description: "Second time value or factor (e.g., '1h 15m', '2.5', '0.75')",
					Required:    true,
				},
			},
		},
		NeedsAccount: true,
		Run:          logic.CalculTimeHandler,
	})
}
