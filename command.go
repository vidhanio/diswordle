package diswordle

import (
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
	{
		Name:        "letters",
		Description: "Show the letters that can be used in the wordle",
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
		case "letters":
			wb.handleWordleLetters(s, i)
		}
	case discordgo.InteractionMessageComponent:
		d := i.MessageComponentData()
		switch d.CustomID {
		}
	}
}

func (wb *WordleBot) handleWordleStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, ok := wb.wordleGame(i)
	if ok {
		warningRespond(s, i, "A wordle is already in progress. Use `/cancel` to cancel it.")

		return
	}

	wordLength := 5
	if len(i.ApplicationCommandData().Options) > 0 {
		wordLength = int(i.ApplicationCommandData().Options[0].IntValue())
	}

	wg, err := wb.newWordleGame(i, wordLength)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to create wordle game")

		errorRespond(s, i, err)

		return
	}

	err = wg.responseCreate()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to create response")

		errorRespond(s, i, err)
	}
}

func (wb *WordleBot) handleWordleGuess(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wg, ok := wb.wordleGame(i)
	if !ok {
		warningRespond(s, i, "No wordle in progress. Use `/start` to start a new game.")

		return
	}

	guess := i.ApplicationCommandData().Options[0].StringValue()
	_, err := wg.Guess(guess)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	err = wg.responseUpdate()
	if err != nil {
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to edit interaction")
		}

		errorRespond(s, i, err)

		return
	}

	successRespond(s, i, "Guess accepted")

	if wg.Done() {
		delete(wb.wordles[i.GuildID], i.Member.User.ID)
	}
}

func (wb *WordleBot) handleWordleCancel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wg, ok := wb.wordleGame(i)
	if !ok {
		warningRespond(s, i, "No wordle in progress. Use `/start` to start a new game.")

		return
	}

	wg.Cancel()

	err := wg.responseUpdate()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to delete interaction")

		errorRespond(s, i, err)
	}

	delete(wb.wordles[i.GuildID], i.Member.User.ID)

	successRespond(s, i, "Wordle cancelled")
}

func (wb *WordleBot) handleWordleShow(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wg, ok := wb.wordleGame(i)
	if !ok {
		warningRespond(s, i, "No wordle in progress. Use `/start` to start a new game.")

		return
	}

	err := wg.responseDelete()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to delete response")

		errorRespond(s, i, err)
	}

	wg.setInteraction(i)
	err = wg.responseCreate()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to create response")

		errorRespond(s, i, err)
	}
}

func (wb *WordleBot) handleWordleLetters(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wg, ok := wb.wordleGame(i)
	if !ok {
		warningRespond(s, i, "No wordle in progress. Use `/start` to start a new game.")

		return
	}

	respond(s, i, ephemeralify(embedResponse(wg.lettersEmbed())))
}
