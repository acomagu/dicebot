package discord

import (
	"github.com/acomagu/dicebot/bot"
	"github.com/acomagu/dicebot/soundplayer"
	"github.com/bwmarrin/discordgo"
)

var _ bot.Session = (*Session)(nil)

type Session struct {
	session *discordgo.Session
}

func NewSession(dgsession *discordgo.Session) *Session {
	return &Session{
		session: dgsession,
	}
}

func (r *Session) GetMe() (*bot.User, error) {
	return r.User("@me")
}

func (r *Session) User(userID string) (*bot.User, error) {
	u, err := r.session.User(userID)
	if err != nil {
		return nil, err
	}

	return &bot.User{
		ID: u.ID,
	}, nil
}

func (r *Session) ChannelMessageSend(guildID, channelID string) (*bot.Message, error) {
	_, err := r.session.ChannelMessageSend(guildID, channelID)
	if err != nil {
		return nil, err
	}

	return &bot.Message{}, nil
}

func (r *Session) ChannelVoiceJoin(guildID, channelID string, mute, deaf bool) (soundplayer.VC, error) {
	vc, err := r.session.ChannelVoiceJoin(guildID, channelID, mute, deaf)
	if err != nil {
		return nil, err
	}

	return &VC{
		conn: vc,
	}, nil
}

func (r *Session) UpdateStatus(idle int, game string) error {
	return r.session.UpdateStatus(idle, game)
}

func (r *Session) VoiceStates(guildID string) ([]*bot.VoiceState, error) {
	guild, err := r.session.State.Guild(guildID)
	if err != nil {
		return nil, err
	}

	var vss []*bot.VoiceState
	for _, vs := range guild.VoiceStates {
		vss = append(vss, &bot.VoiceState{
			UserID:    vs.UserID,
			ChannelID: vs.ChannelID,
		})
	}
	return vss, nil
}

func (r *Session) HasSendMessagePermission(userID string, channelID string) bool {
	p, err := r.session.State.UserChannelPermissions(userID, channelID)
	return (err == nil) && (p&discordgo.PermissionSendMessages == discordgo.PermissionSendMessages)
}
