package diswordle

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/vidhanio/wordle"
)

type WordleBot struct {
	wordles      map[string]map[string]*wordleGame
	guildWordles map[string]*guildWordleGame
	session      *discordgo.Session

	dictionary          []string
	common              []string
	guessesAllowed      int
	guildGuessesAllowed int

	emojiGuilds    map[wordle.GuessResult]string
	miscEmojiGuild string

	emojiMap   map[wordle.GuessResult][26]*discordgo.Emoji
	blankEmoji *discordgo.Emoji
	emptyEmoji *discordgo.Emoji
}

func New(dictionary, common []string, guessesAllowed, guildGuessesAllowed int, botToken string, emojiGuilds map[wordle.GuessResult]string, miscEmojiGuild string) (*WordleBot, error) {
	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	wb := &WordleBot{
		wordles:             make(map[string]map[string]*wordleGame),
		guildWordles:        make(map[string]*guildWordleGame),
		session:             session,
		dictionary:          dictionary,
		common:              common,
		guessesAllowed:      guessesAllowed,
		guildGuessesAllowed: guildGuessesAllowed,
		miscEmojiGuild:      miscEmojiGuild,
		emojiGuilds:         emojiGuilds,
	}

	wb.session.AddHandler(wb.setEmojis)
	wb.session.AddHandler(wb.setGuilds)
	wb.session.AddHandler(wb.createGuild)

	wb.session.AddHandler(wb.handleWordle)
	wb.session.AddHandler(wb.handleGuildWordle)

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
		append(wordleApplicationCommands, guildWordleApplicationCommands...),
	)

	return err
}

func (wb *WordleBot) Stop() error {
	return wb.session.Close()
}
