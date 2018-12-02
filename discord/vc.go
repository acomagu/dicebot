package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type VC struct {
	conn *discordgo.VoiceConnection
}

func (vc *VC) Close() error {
	return vc.conn.Disconnect()
}

func (vc *VC) Write(ctx context.Context, frames [][]byte) error {
	for _, frame := range frames {
		select {
		case vc.conn.OpusSend <- frame:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (vc *VC) Speaking(b bool) error {
	return vc.conn.Speaking(b)
}
