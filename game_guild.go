package diswordle

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/vidhanio/wordle"
)

type guildWordleGame struct {
	*wordle.Wordle
	session     *discordgo.Session
	interaction *discordgo.Interaction

	emojiMap   map[wordle.GuessResult][26]*discordgo.Emoji
	blankEmoji *discordgo.Emoji
	emptyEmoji *discordgo.Emoji

	cancelVotes map[string]bool
}

func (wb *WordleBot) guildWordleGame(i *discordgo.InteractionCreate) (*guildWordleGame, bool) {
	wg, ok := wb.guildWordles[i.GuildID]

	return wg, ok
}

func (wb *WordleBot) newGuildWordleGame(i *discordgo.InteractionCreate, wordLength int) (*guildWordleGame, error) {
	wg, ok := wb.guildWordleGame(i)
	if ok {
		return wg, nil
	}

	w, err := wordle.New(wordLength, wb.guildGuessesAllowed, wb.dictionary, wb.common)
	if err != nil {
		return nil, err
	}

	wb.guildWordles[i.GuildID] = &guildWordleGame{
		Wordle:      w,
		session:     wb.session,
		interaction: i.Interaction,
		emojiMap:    wb.emojiMap,
		blankEmoji:  wb.blankEmoji,
		emptyEmoji:  wb.emptyEmoji,

		cancelVotes: make(map[string]bool),
	}

	wb.guildWordles[i.GuildID].addPlayer(i.Member.User.ID)

	return wb.guildWordles[i.GuildID], nil
}

func (wg *guildWordleGame) emoji(guessResult wordle.GuessResult, c rune) string {
	return wg.emojiMap[guessResult][c-'a'].MessageFormat()
}

func (wg *guildWordleGame) String() string {
	builder := &strings.Builder{}

	for i, wordGuessResults := range wg.GuessResults() {
		guess := wg.Guesses()[i]
		for j, guessResult := range wordGuessResults {
			builder.WriteString(wg.emoji(guessResult, rune(guess[j])))
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

func (wg *guildWordleGame) embed() *discordgo.MessageEmbed {
	guild, err := wg.session.Guild(wg.interaction.GuildID)
	if err != nil {
		return errorEmbed(err)
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: guild.IconURL(),
			Name:    guild.Name,
		},
		Title:       "Guild Wordle",
		Description: wg.String(),
		Color:       wordleBlack,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Guesses left: %d | Made with ❤️ & Go by vidhan#0001", wg.GuessesLeft()),
		},
	}

	if wg.Won() {
		embed.Title = "Guild Wordle - Won"
		embed.Color = wordleGreen
	} else if wg.Cancelled() || wg.Lost() {
		if wg.Cancelled() {
			embed.Title = "Guild Wordle - Cancelled"
		} else {
			embed.Title = "Guild Wordle - Lost"
		}

		embed.Color = wordleRed

		builder := &strings.Builder{}

		builder.WriteString(embed.Description)

		builder.WriteRune('\n')
		builder.WriteString("The word was: ")
		builder.WriteRune('\n')

		for _, char := range wg.Word() {
			builder.WriteString(wg.emoji(wordle.GuessResultCorrect, char))
		}

		embed.Description = builder.String()
	}

	return embed
}

func (wg *guildWordleGame) lettersEmbed() *discordgo.MessageEmbed {
	keyboard := [26]rune{
		'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p',
		'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l',
		'z', 'x', 'c', 'v', 'b', 'n', 'm',
	}

	embed := &discordgo.MessageEmbed{
		Title: "Letters",
		Color: wordleBlack,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Made with ❤️ & Go by vidhan#0001",
		},
	}

	letters := wg.Letters()
	builder := &strings.Builder{}

	for _, key := range keyboard {
		for c, result := range letters {
			if c+'a' == int(key) {
				builder.WriteString(wg.emoji(result, key))

				if key == 'p' {
					builder.WriteRune('\n')
				}
				if key == 'l' {
					builder.WriteRune('\n')
					builder.WriteString(wg.blankEmoji.MessageFormat())
				}
			}
		}
	}

	embed.Description = builder.String()

	return embed
}

func (wg *guildWordleGame) voteEmbed() *discordgo.MessageEmbed {
	guild, err := wg.session.Guild(wg.interaction.GuildID)
	if err != nil {
		return errorEmbed(err)
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: guild.IconURL(),
			Name:    guild.Name,
		},
	}

	if wg.votes() >= wg.votesNeeded() {
		embed.Title = "Cancel Vote - Passed"
		embed.Color = wordleRed
	} else if wg.votes() > 0 {
		embed.Title = "Cancel Vote - In Progress"
		embed.Color = wordleBlack
	} else {
		embed.Title = "Cancel Vote - Not Started"
		embed.Color = wordleGreen
	}

	builder := &strings.Builder{}
	for player, vote := range wg.cancelVotes {
		if vote {
			builder.WriteString(":red_square:")
		} else {
			builder.WriteString(":black_large_square:")
		}
		builder.WriteString(" ")

		builder.WriteString("<@")
		builder.WriteString(player)
		builder.WriteString(">")

		builder.WriteRune('\n')
	}

	embed.Description = builder.String()

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Cancel votes: %d/%d | Made with ❤️ & Go by vidhan#0001", wg.votes(), len(wg.cancelVotes)),
	}

	return embed
}

func (wg *guildWordleGame) addPlayer(userID string) {
	wg.cancelVotes[userID] = false
}

func (wg *guildWordleGame) votes() int {
	count := 0
	for _, vote := range wg.cancelVotes {
		if vote {
			count++
		}
	}

	return count
}

func (wg *guildWordleGame) votesNeeded() int {
	return int(math.Ceil(float64(len(wg.cancelVotes)) / 2.0))
}

func (wg *guildWordleGame) voteCancel(userID string) error {
	_, ok := wg.cancelVotes[userID]
	if !ok {
		return errors.New("invalid cancel vote: at least one guess must be made by a user for them to vote to cancel")
	} else {
		wg.cancelVotes[userID] = true
	}

	if wg.votes() >= wg.votesNeeded() {
		wg.Cancel()
	}

	return nil
}

func (wg *guildWordleGame) responseCreate() error {
	return wg.session.InteractionRespond(wg.interaction, embedResponse(wg.embed()))
}

func (wg *guildWordleGame) responseUpdate() error {
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

func (wg *guildWordleGame) responseDelete() error {
	return wg.session.InteractionResponseDelete(wg.session.State.User.ID, wg.interaction)
}

func (wg *guildWordleGame) setInteraction(i *discordgo.InteractionCreate) {
	wg.interaction = i.Interaction
}
