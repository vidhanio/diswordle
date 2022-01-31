package diswordle

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var guildWordleApplicationCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "guild",
		Description: "Play a guild wordle",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "start",
				Description: "Start a new guild wordle",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "word_length",
						Description: "The length of the word to be used",
						Type:        discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			{
				Name:        "guess",
				Description: "Guess a word in the guild wordle",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
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
				Description: "Show the guild wordle again",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "cancel",
				Description: "Vote to cancel the current guild wordle",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "votes",
				Description: "Show the cancel votes for the current guild wordle",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "letters",
				Description: "Show the letters used in the current guild wordle",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	},
}

func (wb *WordleBot) handleGuildWordle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		d := i.ApplicationCommandData()
		if d.Name == "guild" {
			switch d.Options[0].Name {
			case "start":
				wb.handleGuildWordleStart(s, i)
			case "guess":
				wb.handleGuildWordleGuess(s, i)
			case "cancel":
				wb.handleGuildWordleCancel(s, i)
			case "show":
				wb.handleGuildWordleShow(s, i)
			case "votes":
				wb.handleGuildWordleVotes(s, i)
			case "letters":
				wb.handleGuildWordleLetters(s, i)
			}
		}
	case discordgo.InteractionMessageComponent:
		d := i.MessageComponentData()
		switch d.CustomID {
		}
	}
}

func (wb *WordleBot) handleGuildWordleStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, ok := wb.guildWordleGame(i)
	if ok {
		warningRespond(s, i, "A guild wordle is already in progress. Use `/guild cancel` to vote to cancel it.")

		return
	}

	wordLength := 7
	if len(i.ApplicationCommandData().Options[0].Options) > 0 {
		wordLength = int(i.ApplicationCommandData().Options[0].Options[0].IntValue())
	}

	wg, err := wb.newGuildWordleGame(i, wordLength)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to create guild wordle game")

		errorRespond(s, i, err)

		return
	}

	err = wg.responseCreate()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to create response")

		errorRespond(s, i, err)

		return
	}
}

func (wb *WordleBot) handleGuildWordleGuess(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wg, ok := wb.guildWordleGame(i)
	if !ok {
		warningRespond(s, i, "No guild wordle in progress. Use `/guild start` to start a new game.")

		return
	}

	guess := i.ApplicationCommandData().Options[0].Options[0].StringValue()
	_, err := wg.Guess(guess)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	err = wg.responseUpdate()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to edit interaction")

		errorRespond(s, i, err)

		return
	}

	wg.addPlayer(i.Member.User.ID)

	successRespond(s, i, "Guess accepted")

	if wg.Done() {
		delete(wb.guildWordles, wg.interaction.GuildID)
	}
}

func (wb *WordleBot) handleGuildWordleCancel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wg, ok := wb.guildWordleGame(i)
	if !ok {
		warningRespond(s, i, "No guild wordle in progress. Use `/guild start` to start a new game.")

		return
	}

	err := wg.voteCancel(i.Member.User.ID)
	if err != nil {
		errorRespond(s, i, err)

		return
	}

	respond(s, i, ephemeralify(embedResponse(wg.voteEmbed())))

	if wg.Done() {
		err = wg.responseUpdate()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to update response")

			errorRespond(s, i, err)
		}

		delete(wb.guildWordles, wg.interaction.GuildID)
	}
}

func (wb *WordleBot) handleGuildWordleShow(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wg, ok := wb.guildWordleGame(i)
	if !ok {
		warningRespond(s, i, "No wordle in progress. Use `/guild start` to start a new game.")

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

		return
	}
}

func (wb *WordleBot) handleGuildWordleVotes(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wg, ok := wb.guildWordleGame(i)
	if !ok {
		warningRespond(s, i, "No guild wordle in progress. Use `/guild start` to start a new game.")

		return
	}

	respond(s, i, ephemeralify(embedResponse(wg.voteEmbed())))
}

func (wb *WordleBot) handleGuildWordleLetters(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wg, ok := wb.guildWordleGame(i)
	if !ok {
		warningRespond(s, i, "No wordle in progress. Use `/start` to start a new game.")

		return
	}

	respond(s, i, ephemeralify(embedResponse(wg.lettersEmbed())))
}
