package bot

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/vidhanio/wordle"
)

type WordleBot struct {
	wordles map[string]map[string]*wordleGame
	session *discordgo.Session

	commonWords    []string
	validWords     []string
	guessesAllowed int
}

func New(commonWords, validWords []string, guessesAllowed int, botToken string) (*WordleBot, error) {
	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	wb := &WordleBot{
		wordles:        make(map[string]map[string]*wordleGame),
		session:        session,
		commonWords:    commonWords,
		validWords:     validWords,
		guessesAllowed: guessesAllowed,
	}

	wb.session.AddHandler(wb.Wordle)

	return wb, nil
}

func (wb *WordleBot) Start() error {
	err := wb.session.Open()
	if err != nil {
		return err
	}

	_, err = wb.session.ApplicationCommandBulkOverwrite(
		wb.session.State.User.ID,
		os.Getenv("DISCORD_GUILD_ID"),
		WordleApplicationCommands,
	)

	return err
}

func (wb *WordleBot) Stop() error {
	return wb.session.Close()
}

func (wb *WordleBot) newWordle(guildID, userID string, wordLength int) (*wordleGame, error) {
	w, err := wordle.New(wordLength, wb.guessesAllowed, wb.commonWords, wb.validWords)
	if err != nil {
		return nil, err
	}

	if wb.wordles[guildID] == nil {
		wb.wordles[guildID] = make(map[string]*wordleGame)
	}

	wb.wordles[guildID][userID] = &wordleGame{
		Wordle:  w,
		session: wb.session,
		guildID: guildID,
		userID:  userID,
	}

	return wb.wordles[guildID][userID], nil
}
