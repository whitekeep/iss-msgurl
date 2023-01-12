package main

import (
	"fmt"
	"github.com/Goscord/goscord/goscord"
	"github.com/Goscord/goscord/goscord/discord"
	"github.com/Goscord/goscord/goscord/gateway"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Create client instance
	client := goscord.New(&gateway.Options{
		Token: os.Getenv("BOT_TOKEN"),
		Intents: gateway.IntentsGuild |
			gateway.IntentGuildMembers |
			gateway.IntentDirectMessages |
			gateway.IntentGuildMessages |
			gateway.IntentMessageContent,
	})

	// Load events
	_ = client.On("ready", OnReady(client))
	_ = client.On("interactionCreate", CommandHandler(client))

	// Login client
	if err := client.Login(); err != nil {
		panic(err)
	}

	// Wait here until term signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session
	client.Close()
}

func OnReady(client *gateway.Session) func() {
	return func() {
		fmt.Println("Logged in as ", client.Me().Tag())

		// Register slash commands
		appCmd := &discord.ApplicationCommand{
			Name:        "test",
			Type:        discord.ApplicationCommandChat,
			Description: "test command",
			Options: []*discord.ApplicationCommandOption{
				{
					Name:        "message_id",
					Type:        discord.ApplicationCommandOptionString,
					Description: "Message ID",
					Required:    true,
				},
			},
		}
		_, _ = client.Application.RegisterCommand(client.Me().Id, "", appCmd)
	}
}

func CommandHandler(client *gateway.Session) func(*discord.Interaction) {
	return func(interaction *discord.Interaction) {
		if interaction.Member == nil {
			return
		}

		// Check if yhe command is test
		if interaction.Data.(discord.ApplicationCommandData).Name != "test" {
			return
		}

		// Get message by ID
		msg, err := client.Channel.GetMessage(interaction.ChannelId, interaction.Data.(discord.ApplicationCommandData).Options[0].Value.(string))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(msg) // or set breakpoint here
	}
}
