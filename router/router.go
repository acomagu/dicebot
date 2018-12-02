package router

import (
	"fmt"
	"os"

	"github.com/acomagu/dicebot/bot"
	"github.com/bwmarrin/discordgo"
)

type Router struct {
	Session *discordgo.Session
}

func (r *Router) Handle(handler bot.Handler) {
	logIf := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	for _, f := range []interface{}{
		func(_ *discordgo.Session, e *discordgo.Ready) {
			logIf(handler.OnReady(&bot.ReadyEvent{}))
		},
		func(_ *discordgo.Session, e *discordgo.MessageCreate) {
			logIf(handler.OnMessageCreate(&bot.MessageCreateEvent{
				Author: &bot.User{
					ID: e.Author.ID,
				},
				GuildID:   e.GuildID,
				ChannelID: e.ChannelID,
				Content:   e.Content,
			}))
		},
		func(_ *discordgo.Session, e *discordgo.GuildCreate) {
			var channels []*bot.Channel
			for _, c := range e.Guild.Channels {
				channels = append(channels, &bot.Channel{
					ID:          c.ID,
					IsGuildText: c.Type == discordgo.ChannelTypeGuildText,
				})
			}

			logIf(handler.OnGuildCreate(&bot.GuildCreateEvent{
				Guild: &bot.Guild{
					Unavailable: e.Guild.Unavailable,
					Channels:    channels,
				},
			}))
		},
	} {
		r.Session.AddHandler(f)
	}
}
