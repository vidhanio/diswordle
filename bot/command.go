package bot

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func (wb *WordleBot) ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	correctRegex := regexp.MustCompile(`([a-z])_green`)
	wrongPositionRegex := regexp.MustCompile(`([a-z])_yellow`)
	wrongRegex := regexp.MustCompile(`([a-z])_black`)

	for _, g := range r.Guilds {
		wb.wordles[g.ID] = make(map[string]*wordleGame)
		switch g.ID {
		case wb.emojiGuilds[0]:
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
		case wb.emojiGuilds[1]:
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
		case wb.emojiGuilds[2]:
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
	}
}
