package bot

import (
	"context"
	"math/rand"
	"strings"
	"testing"

	"github.com/acomagu/dicebot/soundplayer"
	"github.com/matryer/is"
)

//go:generate moq -out session_mock_test.go . Session
//go:generate moq -pkg bot -out joiner_mock_test.go ../soundplayer VoiceChannelJoiner
//go:generate moq -pkg bot -out vc_mock_test.go ../soundplayer VC

func TestBot_OnGuildCreate(t *testing.T) {
	is := is.New(t)

	newSessionMock := func() *SessionMock {
		return &SessionMock{
			GetMeFunc: func() (*User, error) {
				return &User{
					ID: "meid",
				}, nil
			},
			ChannelMessageSendFunc: func(cid, content string) (*Message, error) {
				is.Equal(cid, "fuga")
				is.True(strings.Contains(content, "DICEBOTの準備ができました"))
				return nil, nil
			},
			HasSendMessagePermissionFunc: func(userID, channelID string) bool {
				return true
			},
		}
	}

	newJoinerMock := func() *VoiceChannelJoinerMock {
		vc := &VCMock{
			CloseFunc:    func() error { return nil },
			SpeakingFunc: func(_ bool) error { return nil },
			WriteFunc:    func(_ context.Context, _ [][]byte) error { return nil },
		}
		return &VoiceChannelJoinerMock{
			JoinVoiceChannelFunc: func(guildID, channelID string) (soundplayer.VC, error) {
				return vc, nil
			},
		}
	}

	newEventMock := func() *GuildCreateEvent {
		return &GuildCreateEvent{
			Guild: &Guild{
				Unavailable: false,
				Channels: []*Channel{
					{
						ID:          "fuga",
						IsGuildText: true,
					},
				},
			},
		}
	}

	frames := [][]byte{
		[]byte("abc"),
		[]byte("def"),
		[]byte("ghi"),
	}

	t.Run("normal", func(t *testing.T) {
		session, joiner, event := newSessionMock(), newJoinerMock(), newEventMock()

		bot, err := New(session, joiner, frames)
		is.NoErr(err)

		is.NoErr(bot.OnGuildCreate(event))

		is.Equal(len(session.ChannelMessageSendCalls()), 1)
	})

	t.Run("select channel can be sent message", func(t *testing.T) {
		session, joiner, event := newSessionMock(), newJoinerMock(), newEventMock()

		session.HasSendMessagePermissionFunc = func(userID, channelID string) bool {
			return channelID == "fuga"
		}
		event.Guild.Channels = []*Channel{
			{
				ID:          "hoge",
				IsGuildText: true,
			},
			{
				ID:          "fuga",
				IsGuildText: true,
			},
			{
				ID:          "hage",
				IsGuildText: true,
			},
		}

		bot, err := New(session, joiner, frames)
		is.NoErr(err)

		is.NoErr(bot.OnGuildCreate(event))

		is.Equal(len(session.ChannelMessageSendCalls()), 1)
	})

	t.Run("select text channel", func(t *testing.T) {
		session, joiner, event := newSessionMock(), newJoinerMock(), newEventMock()

		event.Guild.Channels = []*Channel{
			{
				ID:          "hoge",
				IsGuildText: false,
			},
			{
				ID:          "fuga",
				IsGuildText: true,
			},
			{
				ID:          "hage",
				IsGuildText: false,
			},
		}

		bot, err := New(session, joiner, frames)
		is.NoErr(err)

		is.NoErr(bot.OnGuildCreate(event))

		is.Equal(len(session.ChannelMessageSendCalls()), 1)
	})
}

func TestBot_OnMessageCreate(t *testing.T) {
	is := is.New(t)

	channelID := "channelid"
	voiceChannelID := "voicechannelid"
	guildID := "guildid"
	userID := "userid"

	frames := [][]byte{
		[]byte("abc"),
		[]byte("def"),
		[]byte("ghi"),
	}

	newSessionMock := func() *SessionMock {
		return &SessionMock{
			GetMeFunc: func() (*User, error) {
				return &User{
					ID: "meid",
				}, nil
			},
			VoiceStatesFunc: func(gid string) ([]*VoiceState, error) {
				return []*VoiceState{
					{
						UserID:    userID,
						ChannelID: voiceChannelID,
					},
				}, nil
			},
			ChannelMessageSendFunc: func(cid, content string) (*Message, error) {
				is.Equal(cid, channelID)
				is.Equal(content, "10")
				return nil, nil
			},
		}
	}

	newVCMock := func() *VCMock {
		return &VCMock{
			CloseFunc:    func() error { return nil },
			SpeakingFunc: func(_ bool) error { return nil },
			WriteFunc: func(_ context.Context, fs [][]byte) error {
				is.Equal(fs, frames)
				return nil
			},
		}
	}

	newJoinerMock := func(vc *VCMock) *VoiceChannelJoinerMock {
		return &VoiceChannelJoinerMock{
			JoinVoiceChannelFunc: func(gid, cid string) (soundplayer.VC, error) {
				is.Equal(gid, guildID)
				is.Equal(cid, voiceChannelID)
				return vc, nil
			},
		}
	}

	newEventMock := func() *MessageCreateEvent {
		return &MessageCreateEvent{
			ChannelID: channelID,
			GuildID:   guildID,
			Content:   "!dice 2d10",
			Author: &User{
				ID: userID,
			},
		}
	}

	t.Run("normal", func(t *testing.T) {
		rand.Seed(0)

		vc := newVCMock()
		session, joiner, event := newSessionMock(), newJoinerMock(vc), newEventMock()

		bot, err := New(session, joiner, frames)
		is.NoErr(err)

		is.NoErr(bot.OnMessageCreate(event))

		is.Equal(len(session.ChannelMessageSendCalls()), 1)
		is.Equal(len(vc.WriteCalls()), 1)
		is.Equal(vc.WriteCalls()[0].In2, frames)
	})

	t.Run("help", func(t *testing.T) {
		vc := newVCMock()
		session, joiner, event := newSessionMock(), newJoinerMock(vc), newEventMock()

		event.Content = "!dice"
		session.ChannelMessageSendFunc = func(cid, content string) (*Message, error) {
			is.Equal(cid, channelID)
			is.True(strings.Contains(content, "よくわかりません"))
			return nil, nil
		}

		bot, err := New(session, joiner, frames)
		is.NoErr(err)

		is.NoErr(bot.OnMessageCreate(event))

		is.Equal(len(joiner.JoinVoiceChannelCalls()), 0)
		is.Equal(len(session.ChannelMessageSendCalls()), 1)
		is.Equal(len(vc.WriteCalls()), 0)
	})

	t.Run("reply in a proper channel", func(t *testing.T) {
		rand.Seed(0)

		vc := newVCMock()
		session, joiner, event := newSessionMock(), newJoinerMock(vc), newEventMock()

		session.VoiceStatesFunc = func(gid string) ([]*VoiceState, error) {
			return []*VoiceState{
				{
					UserID:    "userid2",
					ChannelID: "voicechannelid2",
				},
				{
					UserID:    userID,
					ChannelID: voiceChannelID,
				},
				{
					UserID:    "userid3",
					ChannelID: "voicechannelid3",
				},
			}, nil
		}

		bot, err := New(session, joiner, frames)
		is.NoErr(err)

		is.NoErr(bot.OnMessageCreate(event))

		is.Equal(len(joiner.JoinVoiceChannelCalls()), 1)
		is.Equal(len(vc.WriteCalls()), 1)
		is.Equal(len(session.ChannelMessageSendCalls()), 1)
	})
}
