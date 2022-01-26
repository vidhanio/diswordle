package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var wordleApplicationCommands = []*discordgo.ApplicationCommand{
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
		Name:        "show",
		Description: "Show the wordle again",
		Type:        discordgo.ChatApplicationCommand,
	},
	{
		Name:        "cancel",
		Description: "Cancel the current wordle",
		Type:        discordgo.ChatApplicationCommand,
	},
}

func (wb *WordleBot) handleWordle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		d := i.ApplicationCommandData()
		switch d.Name {
		case "start":
			wb.handleWordleStart(s, i)
		case "guess":
			wb.handleWordleGuess(s, i)
		case "cancel":
			wb.handleWordleCancel(s, i)
		case "show":
			wb.handleWordleShow(s, i)
		}
	case discordgo.InteractionMessageComponent:
		d := i.MessageComponentData()
		switch d.CustomID {
		}
	}
}

func (wb *WordleBot) handleWordleStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, ok := wb.wordles[i.GuildID][i.Member.User.ID]
	if ok {
		err := errorRespond(s, i, errors.New("A wordle is already in progress. Use `/cancel` to cancel it."))
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to respond to interaction")
		}

		return
	}

	wordLength := 5
	if len(i.ApplicationCommandData().Options) > 0 {
		wordLength = int(i.ApplicationCommandData().Options[0].IntValue())
	}

	wg, err := wb.newWordleGame(i, wordLength)
	if err != nil {
		err := errorRespond(s, i, err)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to respond to interaction")
		}

		return
	}

	err = wg.responseCreate()
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

func (wb *WordleBot) handleWordleGuess(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	err = wg.responseUpdate()
	if err != nil {
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to edit interaction")
		}

		return
	}

	err = successRespond(s, i, "Guess accepted")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to respond to interaction")

		return
	}

	if wg.Done() {
		delete(wb.wordles[i.GuildID], i.Member.User.ID)
	}
}

func (wb *WordleBot) handleWordleCancel(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	wg.Cancel()

	err := wg.responseUpdate()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to delete interaction")
	}

	delete(wb.wordles[i.GuildID], i.Member.User.ID)

	err = successRespond(s, i, "Wordle canceled")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to respond to interaction")

		return
	}
}

func (wb *WordleBot) handleWordleShow(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	err := wg.responseDelete()
	if err != nil {
		err := errorRespond(s, i, err)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to delete interaction")
		}

		return
	}

	wg.setInteraction(i)
	err = wg.responseCreate()
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
