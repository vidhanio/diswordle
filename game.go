package diswordle

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vidhanio/wordle"
)

type wordleGame struct {
	*wordle.Wordle
	session     *discordgo.Session
	interaction *discordgo.Interaction

	emojiMap   [3][26]*discordgo.Emoji
	emptyEmoji *discordgo.Emoji
}

func (wb *WordleBot) getWordleGame(i *discordgo.InteractionCreate) (*wordleGame, bool) {
	g, ok := wb.wordles[i.GuildID]
	if !ok {
		g = make(map[string]*wordleGame)
		wb.wordles[i.GuildID] = g
	}

	w, ok := g[i.Member.User.ID]

	return w, ok
}

func (wb *WordleBot) newWordleGame(i *discordgo.InteractionCreate, wordLength int) (*wordleGame, error) {
	wg, ok := wb.getWordleGame(i)
	if ok {
		return wg, nil
	}

	w, err := wordle.New(wordLength, wb.guessesAllowed, wb.dictionary, wb.common)
	if err != nil {
		return nil, err
	}

	wb.wordles[i.GuildID][i.Member.User.ID] = &wordleGame{
		Wordle:      w,
		session:     wb.session,
		interaction: i.Interaction,
		emojiMap:    wb.emojiMap,
		emptyEmoji:  wb.emptyEmoji,
	}

	return wb.wordles[i.GuildID][i.Member.User.ID], nil
}

func (wg *wordleGame) emoji(guessType wordle.GuessType, c byte) string {
	return wg.emojiMap[guessType][c-'a'].MessageFormat()
}

func (wg *wordleGame) String() string {
	builder := &strings.Builder{}

	for i, wordGuessTypes := range wg.GuessTypes() {
		guess := wg.Guesses()[i]
		for j, guessType := range wordGuessTypes {
			builder.WriteString(wg.emoji(guessType, guess[j]))
		}
		builder.WriteRune('\n')
	}

	i := 0
	for i < wg.GuessesLeft() {
		j := 0
		for j < wg.WordLength() {
			builder.WriteString(wg.emptyEmoji.MessageFormat())
			j++
		}
		builder.WriteRune('\n')
		i++
	}

	return builder.String()
}

func (wg *wordleGame) embed() *discordgo.MessageEmbed {
	user, err := wg.session.User(wg.interaction.Member.User.ID)
	if err != nil {
		return errorEmbed(err)
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: user.AvatarURL(""),
			Name:    user.Username,
		},
		Title:       "Wordle",
		Description: wg.String(),
		Color:       wordleBlack,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Guesses left: %d | Made with ❤️ & Go by vidhan#0001", wg.GuessesLeft()),
		},
	}

	if wg.Won() {
		embed.Title = "Wordle - Won"
		embed.Color = wordleGreen
	} else if wg.Cancelled() || wg.Lost() {
		if wg.Cancelled() {
			embed.Title = "Wordle - Cancelled"
		} else {
			embed.Title = "Wordle - Lost"
		}

		embed.Color = wordleRed

		builder := &strings.Builder{}

		builder.WriteString(embed.Description)

		builder.WriteRune('\n')
		builder.WriteString("The word was: ")
		builder.WriteRune('\n')

		for _, char := range wg.Word() {
			builder.WriteString(wg.emojiMap[wordle.GuessTypeCorrect][char-'a'].MessageFormat())
		}

		embed.Description = builder.String()
	}

	return embed
}

func (wg *wordleGame) responseCreate() error {
	return wg.session.InteractionRespond(wg.interaction, embedResponse(wg.embed()))
}

func (wg *wordleGame) responseUpdate() error {
	_, err := wg.session.InteractionResponseEdit(
		wg.session.State.User.ID,
		wg.interaction,
		&discordgo.WebhookEdit{
			Embeds: []*discordgo.MessageEmbed{
				wg.embed(),
			},
		},
	)

	return err
}

func (wg *wordleGame) responseDelete() error {
	return wg.session.InteractionResponseDelete(wg.session.State.User.ID, wg.interaction)
}

func (wg *wordleGame) setInteraction(i *discordgo.InteractionCreate) {
	wg.interaction = i.Interaction
}
