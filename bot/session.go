package bot

import (
	"github.com/acomagu/dicebot/soundplayer"
)

type Session interface {
	GetMe() (*User, error)
	User(userID string) (*User, error)
	ChannelMessageSend(guildID, channelID string) (*Message, error)
	ChannelVoiceJoin(guildID, channelID string, mute, deaf bool) (soundplayer.VC, error)
	UpdateStatus(idle int, game string) error
	VoiceStates(guildID string) ([]*VoiceState, error)
	HasSendMessagePermission(userID, channelID string) bool
}
