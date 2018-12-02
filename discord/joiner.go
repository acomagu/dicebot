package discord

import (
	"github.com/acomagu/dicebot/soundplayer"
	"github.com/bwmarrin/discordgo"
)

type VoiceChannelJoiner struct {
	session *discordgo.Session
}

func NewVoiceChannelJoiner(session *discordgo.Session) *VoiceChannelJoiner {
	return &VoiceChannelJoiner{
		session: session,
	}
}

func (c *VoiceChannelJoiner) JoinVoiceChannel(guildID, channelID string) (soundplayer.VC, error) {
	vc, err := c.session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return nil, err
	}

	return &VC{
		conn: vc,
	}, nil
}
