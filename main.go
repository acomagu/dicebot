package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/acomagu/dicebot/bot"
	"github.com/acomagu/dicebot/discord"
	"github.com/acomagu/dicebot/router"
	_ "github.com/acomagu/dicebot/statik"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
)

var token = os.Getenv("DISCORD_TOKEN")
var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "80"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if token == "" {
		return errors.New("No token provided. Set DISCORD_TOKEN.")
	}

	files, err := fs.New()
	if err != nil {
		return err
	}

	soundfile, err := files.Open("/dicesound.dca")
	if err != nil {
		return err
	}

	frames, err := loadSound(soundfile)
	if err != nil {
		return err
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	joiner := discord.NewVoiceChannelJoiner(session)

	r := &router.Router{
		Session: session,
	}

	firestore := firestore.NewFirestore()

	bot, err := bot.New(discord.NewSession(session), joiner, firestore, frames)
	if err != nil {
		return err
	}

	r.Handle(bot)

	if err := session.Open(); err != nil {
		return err
	}
	defer session.Close()
	bot.Wait()

	return nil
}

func loadSound(file io.Reader) ([][]byte, error) {
	decoder := dca.NewDecoder(file)

	var frames [][]byte
	for {
		frame, err := decoder.OpusFrame()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		frames = append(frames, frame)
	}

	return frames, nil
}
