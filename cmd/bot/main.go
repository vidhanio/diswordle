package main

import (
	"bufio"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vidhanio/diswordle"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	dictionaryFile := ""
	flag.StringVar(&dictionaryFile, "dictionary", "dictionary.txt", "filename of dictionary to use")

	commonFile := ""
	flag.StringVar(&commonFile, "common", "common.txt", "filename of common words to use")

	guesses := 0
	flag.IntVar(&guesses, "guesses", 6, "Number of guesses allowed")

	guildGuesses := 0
	flag.IntVar(&guildGuesses, "guild-guesses", 10, "Number of guesses allowed for guild wordle")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load .env file")
	}

	log.Debug().Msg("Reading word lists...")

	dictionary := []string{}

	file, err := os.Open(dictionaryFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open dictionary.txt")
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dictionary = append(dictionary, scanner.Text())
	}

	err = file.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to close dictionary.txt")
	}

	common := []string{}

	file, err = os.Open(commonFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open common.txt")
	}

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		common = append(common, scanner.Text())
	}

	err = file.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to close common.txt")
	}

	log.Debug().Msg("Done reading word lists.")

	log.Debug().Msg("Starting bot...")
	bot, err := diswordle.New(
		dictionary,
		common,
		guesses,
		guildGuesses,
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
