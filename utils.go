package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

const ephemeralFlag = 1 << 6

const (
	wordleGreen  = 0x7AB457
	wordleYellow = 0xF2E8B3
	wordleBlack  = 0x293137
	wordleRed    = 0xDF2640
)

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

func embed(title string, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       wordleBlack,
	}
}

func successEmbed(m string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Success",
		Description: m,
		Color:       wordleGreen,
	}
}

func warningEmbed(m string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Warning",
		Description: m,
		Color:       wordleYellow,
	}
}

func errorEmbed(err error) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Error",
		Description: err.Error(),
		Color:       wordleRed,
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

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, r *discordgo.InteractionResponse) {
	err := s.InteractionRespond(i.Interaction, r)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to respond to interaction")
	}
}

func successRespond(s *discordgo.Session, i *discordgo.InteractionCreate, m string) {
	respond(s, i, ephemeralify(embedResponse(successEmbed(m))))
}

func warningRespond(s *discordgo.Session, i *discordgo.InteractionCreate, m string) {
	respond(s, i, ephemeralify(embedResponse(warningEmbed(m))))
}

func errorRespond(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	respond(s, i, ephemeralify(embedResponse(errorEmbed(err))))
}
