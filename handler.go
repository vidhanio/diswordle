package discordle

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
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
	correctRegex := regexp.MustCompile(`([a-z])_green`)
	wrongPositionRegex := regexp.MustCompile(`([a-z])_yellow`)
	wrongRegex := regexp.MustCompile(`([a-z])_black`)

	for _, g := range r.Guilds {
		if g.ID == wb.emojiGuilds[0] {
			emojis, err := s.GuildEmojis(g.ID)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get guild emojis")
			}

			for _, e := range emojis {
				match := correctRegex.FindStringSubmatch(e.Name)
				if match != nil {
					wb.emojiMap[0][match[1][0]-'a'] = e
				}
			}
		}

		if g.ID == wb.emojiGuilds[1] {
			emojis, err := s.GuildEmojis(g.ID)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get guild emojis")
			}

			for _, e := range emojis {
				match := wrongPositionRegex.FindStringSubmatch(e.Name)
				if match != nil {
					wb.emojiMap[1][match[1][0]-'a'] = e
				}
			}
		}

		if g.ID == wb.emojiGuilds[2] {
			emojis, err := s.GuildEmojis(g.ID)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get guild emojis")
			}

			for _, e := range emojis {
				match := wrongRegex.FindStringSubmatch(e.Name)
				if match != nil {
					wb.emojiMap[2][match[1][0]-'a'] = e
				}
			}
		}

		if g.ID == wb.emptyEmojiGuild {
			emojis, err := s.GuildEmojis(wb.emptyEmojiGuild)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get guild emojis")
			}

			for _, e := range emojis {
				if e.Name == "empty" {
					wb.emptyEmoji = e
				}
			}
		}
	}
}
