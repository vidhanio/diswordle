package discordle

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

type WordleBot struct {
	wordles map[string]map[string]*wordleGame
	session *discordgo.Session

	commonWords    []string
	validWords     []string
	guessesAllowed int

	emojiGuilds     [3]string
	emojiMap        [3][26]*discordgo.Emoji
	emptyEmojiGuild string
	emptyEmoji      *discordgo.Emoji
}

func New(commonWords, validWords []string, guessesAllowed int, botToken string, emojiGuilds [3]string, emptyEmojiGuild string) (*WordleBot, error) {
	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	wb := &WordleBot{
		wordles:         make(map[string]map[string]*wordleGame),
		session:         session,
		commonWords:     commonWords,
		validWords:      validWords,
		guessesAllowed:  guessesAllowed,
		emptyEmojiGuild: emptyEmojiGuild,
		emojiGuilds:     emojiGuilds,
	}

	wb.session.AddHandler(wb.setEmojis)
	wb.session.AddHandler(wb.setGuilds)
	wb.session.AddHandler(wb.createGuild)

	wb.session.AddHandler(wb.handleWordle)

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
		wordleApplicationCommands,
	)

	return err
}

func (wb *WordleBot) Stop() error {
	return wb.session.Close()
}
