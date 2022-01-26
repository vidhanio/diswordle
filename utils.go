package bot

import "github.com/bwmarrin/discordgo"

const ephemeralFlag = 1 << 6

func ephemeralify(r *discordgo.InteractionResponse) *discordgo.InteractionResponse {
	r.Data.Flags |= ephemeralFlag
	return r
}

func deferred(r *discordgo.InteractionResponse) *discordgo.InteractionResponse {
	r.Type = discordgo.InteractionResponseDeferredChannelMessageWithSource
	return r
}

func contentResponse(c string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: c,
		},
	}
}

func embedResponse(es ...*discordgo.MessageEmbed) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: es,
		},
	}
}

func successEmbed(m string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Success",
		Description: m,
		Color:       0x57F287,
	}
}

func errorEmbed(err error) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Error",
		Description: err.Error(),
		Color:       0xED4245,
	}
}

func contentMessage(c string) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Content: c,
	}
}

func embedMessage(es ...*discordgo.MessageEmbed) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embeds: es,
	}
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, r *discordgo.InteractionResponse) error {
	return s.InteractionRespond(i.Interaction, r)
}

func successRespond(s *discordgo.Session, i *discordgo.InteractionCreate, m string) error {
	return respond(s, i, ephemeralify(embedResponse(successEmbed(m))))
}

func errorRespond(s *discordgo.Session, i *discordgo.InteractionCreate, err error) error {
	return respond(s, i, ephemeralify(embedResponse(errorEmbed(err))))
}
