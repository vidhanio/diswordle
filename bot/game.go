package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vidhanio/wordle"
)

type wordleGame struct {
	*wordle.Wordle
	session *discordgo.Session

	guildID string
	userID  string

	emojiMap   [3][26]*discordgo.Emoji
	emptyEmoji *discordgo.Emoji
}

func (wb *WordleBot) newWordleGame(guildID, userID string, wordLength int) (*wordleGame, error) {
	w, err := wordle.New(wordLength, wb.guessesAllowed, wb.commonWords, wb.validWords)
	if err != nil {
		return nil, err
	}

	if wb.wordles[guildID] == nil {
		wb.wordles[guildID] = make(map[string]*wordleGame)
	}

	wb.wordles[guildID][userID] = &wordleGame{
		Wordle:     w,
		session:    wb.session,
		guildID:    guildID,
		userID:     userID,
		emojiMap:   wb.emojiMap,
		emptyEmoji: wb.emptyEmoji,
	}

	return wb.wordles[guildID][userID], nil
}

func (wm wordleGame) String() string {
	builder := new(strings.Builder)

	for i, guessType := range wm.GuessTypes() {
		guess := wm.Guesses()[i]
		for j, char := range guessType {
			switch char {
			case wordle.GuessTypeCorrect:
				builder.WriteString(wm.emojiMap[0][guess[j]-'a'].MessageFormat())
			case wordle.GuessTypeWrongPosition:
				builder.WriteString(wm.emojiMap[1][guess[j]-'a'].MessageFormat())
			case wordle.GuessTypeWrong:
				builder.WriteString(wm.emojiMap[2][guess[j]-'a'].MessageFormat())
			}
		}
		builder.WriteRune('\n')
	}

	i := 0
	for i < wm.GuessesLeft() {
		j := 0
		for j < wm.WordLength() {
			builder.WriteString(wm.emptyEmoji.MessageFormat())
			j++
		}
		builder.WriteRune('\n')
		i++
	}

	return builder.String()
}

func (wm wordleGame) embed() *discordgo.MessageEmbed {
	user, err := wm.session.User(wm.userID)
	if err != nil {
		return errorEmbed(err)
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: user.AvatarURL(""),
			Name:    user.Username,
		},
		Title:       "Wordle",
		Description: wm.String(),
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: "https://avatars.githubusercontent.com/u/41439633?v=4",
			Text:    fmt.Sprintf("Guesses left: %d | Made with ❤️ & Go by Vidhan", wm.GuessesLeft()),
		},
	}

	if wm.Won() {
		embed.Title = "Wordle - You won!"
		embed.Color = 0x57F287
	}

	return embed
}
