package diswordle

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/vidhanio/wordle"
)

func (wb *WordleBot) setGuilds(s *discordgo.Session, r *discordgo.Ready) {
	for _, g := range r.Guilds {
		wb.wordles[g.ID] = make(map[string]*wordleGame)
	}
}

func (wb *WordleBot) createGuild(s *discordgo.Session, g *discordgo.GuildCreate) {
	_, ok := wb.wordles[g.ID]
	if !ok {
		wb.wordles[g.ID] = make(map[string]*wordleGame)
	}
}

func (wb *WordleBot) setEmojis(s *discordgo.Session, r *discordgo.Ready) {
	notGuessedRegex := regexp.MustCompile(`([a-z])_grey`)
	wrongRegex := regexp.MustCompile(`([a-z])_black`)
	wrongPositionRegex := regexp.MustCompile(`([a-z])_yellow`)
	correctRegex := regexp.MustCompile(`([a-z])_green`)

	wb.emojiMap = make(map[wordle.GuessResult][26]*discordgo.Emoji)

	notGuessedEmojis := [26]*discordgo.Emoji{}
	wrongEmojis := [26]*discordgo.Emoji{}
	wrongPositionEmojis := [26]*discordgo.Emoji{}
	correctEmojis := [26]*discordgo.Emoji{}

	for _, g := range r.Guilds {
		if g.ID == wb.emojiGuilds[wordle.GuessResultNotGuessed] {
			emojis, err := s.GuildEmojis(g.ID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get guild emojis")
			}

			for _, e := range emojis {
				match := notGuessedRegex.FindStringSubmatch(e.Name)
				if match != nil {
					notGuessedEmojis[match[1][0]-'a'] = e
				}
			}
		}

		if g.ID == wb.emojiGuilds[wordle.GuessResultCorrect] {
			emojis, err := s.GuildEmojis(g.ID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get guild emojis")
			}

			for _, e := range emojis {
				match := correctRegex.FindStringSubmatch(e.Name)
				if match != nil {
					correctEmojis[match[1][0]-'a'] = e
				}
			}
		}

		if g.ID == wb.emojiGuilds[wordle.GuessResultWrongPosition] {
			emojis, err := s.GuildEmojis(g.ID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get guild emojis")
			}

			for _, e := range emojis {
				match := wrongPositionRegex.FindStringSubmatch(e.Name)
				if match != nil {
					wrongPositionEmojis[match[1][0]-'a'] = e
				}
			}
		}

		if g.ID == wb.emojiGuilds[wordle.GuessResultWrong] {
			emojis, err := s.GuildEmojis(g.ID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get guild emojis")
			}

			for _, e := range emojis {
				match := wrongRegex.FindStringSubmatch(e.Name)
				if match != nil {
					wrongEmojis[match[1][0]-'a'] = e
				}
			}
		}

		if g.ID == wb.miscEmojiGuild {
			emojis, err := s.GuildEmojis(wb.miscEmojiGuild)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get guild emojis")
			}

			for _, e := range emojis {
				if e.Name == "blank" {
					wb.blankEmoji = e
				} else if e.Name == "empty" {
					wb.emptyEmoji = e
				}
			}
		}
	}

	wb.emojiMap[wordle.GuessResultNotGuessed] = notGuessedEmojis
	wb.emojiMap[wordle.GuessResultWrong] = wrongEmojis
	wb.emojiMap[wordle.GuessResultWrongPosition] = wrongPositionEmojis
	wb.emojiMap[wordle.GuessResultCorrect] = correctEmojis
}
