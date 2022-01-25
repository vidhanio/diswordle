package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var WordleApplicationCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "start",
		Description: "Start a new wordle",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "word_length",
				Description: "The length of the word to be used in the wordle",
				Type:        discordgo.ApplicationCommandOptionInteger,
			},
		},
	},
	{
		Name:        "guess",
		Description: "Guess a word in the wordle",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "word",
				Description: "The word to guess",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
	{
		Name:        "stop",
		Description: "Stop the current wordle",
		Type:        discordgo.ChatApplicationCommand,
	},
}

func (wb *WordleBot) HandleWordle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	c := i.ApplicationCommandData()
	switch c.Name {
	case "start":
		wb.WordleStart(s, i)
	case "guess":
		wb.WordleGuess(s, i)
	case "stop":
		wb.WordleStop(s, i)
	}
}

func (wb *WordleBot) WordleStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wordLength := 5
	if len(i.ApplicationCommandData().Options) > 1 {
		wordLength = int(i.ApplicationCommandData().Options[0].IntValue())
	}

	wg, err := wb.newWordleGame(i.GuildID, i.Member.User.ID, wordLength)
	if err != nil {
		err := errorRespond(s, i, err)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to respond to interaction")
		}

		return
	}

	err = respond(s, i, embedResponse(wg.Embed()))
	if err != nil {
		err := errorRespond(s, i, err)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to respond to interaction")
		}

		return
	}
}

func (wb *WordleBot) WordleGuess(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wg, ok := wb.wordles[i.GuildID][i.Member.User.ID]
	if !ok {
		err := errorRespond(s, i, errors.New("No wordle in progress. Use `/start` to start a new game."))
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to respond to interaction")
		}

		return
	}

	guess := i.ApplicationCommandData().Options[0].StringValue()
	_, err := wg.Guess(guess)
	if err != nil {
		err := errorRespond(s, i, err)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to respond to interaction")
		}

		return
	}

	err = respond(s, i, embedResponse(wg.Embed()))
	if err != nil {
		err := errorRespond(s, i, err)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to respond to interaction")
		}

		return
	}

	if wg.Won() {
		delete(wb.wordles[i.GuildID], i.Member.User.ID)
	}
}

func (wb *WordleBot) WordleStop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, ok := wb.wordles[i.GuildID][i.Member.User.ID]
	if !ok {
		err := errorRespond(s, i, errors.New("No wordle in progress. Use `/start` to start a new game."))
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to respond to interaction")
		}

		return
	}

	delete(wb.wordles[i.GuildID], i.Member.User.ID)

	err := successRespond(s, i, "Wordle stopped")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to respond to interaction")

		return
	}
}
