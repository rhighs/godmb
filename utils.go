package main

import (
	dgo "github.com/bwmarrin/discordgo"
)

func InteractionTextRespond(s *dgo.Session, i *dgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: message,
		},
	})
}

func InteractionTextUpdate(s *dgo.Session, i *dgo.InteractionCreate, message string) {
	s.InteractionResponseEdit(i.Interaction, &dgo.WebhookEdit{
		Content: &message,
	})
}

func InteractionRespondDeferred(s *dgo.Session, i *dgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseDeferredChannelMessageWithSource,
		//Data: &dgo.InteractionResponseData{},
	})
}
