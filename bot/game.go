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
}

func (wm wordleGame) String() string {
	builder := new(strings.Builder)

	for _, guess := range wm.CharGuesses() {
		for _, char := range guess {
			switch char {
			case wordle.CharCorrect:
				builder.WriteRune('🟩')
			case wordle.CharWrongPlace:
				builder.WriteRune('🟨')
			case wordle.CharWrong:
				builder.WriteRune('🟥')
			}
		}
		builder.WriteRune('\n')
	}

	return builder.String()
}

func (wm wordleGame) Embed() *discordgo.MessageEmbed {
	user, err := wm.session.User(wm.userID)
	if err != nil {
		return errorEmbed(err)
	}

	return &discordgo.MessageEmbed{
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
}
