package main

import (
	"bufio"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vidhanio/discordle"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load .env file")
	}

	log.Debug().Msg("Reading word lists...")
	validWords := make([]string, 370103)

	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open words.txt")
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		validWords = append(validWords, scanner.Text())
	}

	err = file.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to close words.txt")
	}

	commonWords := make([]string, 10000)

	file, err = os.Open("commonwords.txt")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open commonwords.txt")
	}

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		commonWords = append(commonWords, scanner.Text())
	}

	err = file.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to close commonwords.txt")
	}

	log.Debug().Msg("Done reading word lists.")

	log.Debug().Msg("Starting bot...")
	bot, err := discordle.New(
		commonWords,
		validWords,
		6,
		os.Getenv("DISCORD_BOT_TOKEN"),
		[3]string{os.Getenv("CORRECT_EMOJI_GUILD"), os.Getenv("WRONG_POSITION_EMOJI_GUILD"), os.Getenv("WRONG_EMOJI_GUILD")},
		os.Getenv("EMPTY_EMOJI_GUILD"),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start bot")
	}

	err = bot.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start bot")
	}

	log.Debug().Msg("Bot started. Press Ctrl+C to stop.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Debug().Msg("Stopping bot...")

	err = bot.Stop()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to stop bot")
	}

	log.Debug().Msg("Bot stopped.")
}
