package main

import (
  "github.com/bwmarrin/discordgo"
)

type interactionProcessor func(interaction *discordgo.Interaction) (*discordgo.InteractionResponse, error)

var interactions = map[string]interactionProcessor {
  "test": handleInteractionTest,
}

func handleInteractionTest(interaction *discordgo.Interaction) (*discordgo.InteractionResponse, error) {

  data := interaction.Data.(discordgo.ApplicationCommandInteractionData)

  return &discordgo.InteractionResponse{
    Type: discordgo.InteractionResponseChannelMessageWithSource,
    Data: &discordgo.InteractionResponseData{
      Content: "<@" + interaction.Member.User.ID + "> sent the test command in <#" + interaction.ChannelID + "> with the word \"" + data.Options[0].Value.(string) + "\"!",
    },
  }, nil
}