package soundplayer

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
)

//go:generate moq -out vc_mock_test.go . VC
//go:generate moq -out voice_channel_joiner_mock_test.go . VoiceChannelJoiner

func TestSoundPlayer(t *testing.T) {
	ctx := context.Background()
	sound := [][]byte{[]byte("abc"), []byte("def"), []byte("ghi")}
	guildID := "guildid"
	channelID := "channelid"

	t.Run("first play", func(t *testing.T) {
		t.Parallel()
		is := is.New(t)

		vc := newVCMock(t, sound)
		joiner := newVoiceChannelJoinerMock(t, vc, guildID, channelID, sound)

		player := NewSoundPlayer(joiner, guildID)
		player.setIdleTime(500 * time.Millisecond)
		player.setSpeakBufTime(0)

		player.PlaySound(ctx, channelID, sound)
		time.Sleep(1 * time.Second) // Wait for closing

		is.Equal(len(joiner.JoinVoiceChannelCalls()), 1)
		is.Equal(len(vc.WriteCalls()), 1)
		is.Equal(len(vc.CloseCalls()), 1)
	})

	t.Run("second play not timeout", func(t *testing.T) {
		t.Parallel()
		is := is.New(t)

		vc := newVCMock(t, sound)
		joiner := newVoiceChannelJoinerMock(t, vc, guildID, channelID, sound)

		player := NewSoundPlayer(joiner, guildID)
		player.setIdleTime(500 * time.Millisecond)
		player.setSpeakBufTime(0)

		player.PlaySound(ctx, channelID, sound)
		time.Sleep(30 * time.Millisecond)
		player.PlaySound(ctx, channelID, sound)
		time.Sleep(30 * time.Millisecond)
		player.PlaySound(ctx, channelID, sound)
		time.Sleep(30 * time.Millisecond)
		player.PlaySound(ctx, channelID, sound)
		time.Sleep(1 * time.Second) // Wait for closing

		is.Equal(len(joiner.JoinVoiceChannelCalls()), 1)
		is.Equal(len(vc.WriteCalls()), 4)
		is.Equal(len(vc.CloseCalls()), 1)
	})

	t.Run("changing channelID", func(t *testing.T) {
		t.Parallel()
		is := is.New(t)

		channelID2 := "channelid2"

		vc1 := newVCMock(t, sound)
		vc2 := newVCMock(t, sound)
		vc3 := newVCMock(t, sound)

		var vc1called int

		joiner := &VoiceChannelJoinerMock{
			JoinVoiceChannelFunc: func(g, c string) (VC, error) {
				is.Equal(g, guildID)

				switch c {
				case channelID:
					vc1called++
					if vc1called == 1 {
						return vc1, nil
					}
					return vc3, nil
				case channelID2:
					return vc2, nil
				default:
					t.Error("unexpected channel id: ", c)
					return vc1, nil
				}
			},
		}

		player := NewSoundPlayer(joiner, guildID)
		player.setIdleTime(500 * time.Millisecond)
		player.setSpeakBufTime(0)

		player.PlaySound(ctx, channelID, sound)
		time.Sleep(30 * time.Millisecond)
		player.PlaySound(ctx, channelID2, sound)
		time.Sleep(30 * time.Millisecond)
		player.PlaySound(ctx, channelID, sound)
		time.Sleep(1 * time.Second) // Wait for closing

		is.Equal(len(joiner.JoinVoiceChannelCalls()), 3)
		is.Equal(len(vc1.WriteCalls()), 1)
		is.Equal(len(vc1.CloseCalls()), 1)
		is.Equal(len(vc2.WriteCalls()), 1)
		is.Equal(len(vc2.CloseCalls()), 1)
		is.Equal(len(vc3.WriteCalls()), 1)
		is.Equal(len(vc3.CloseCalls()), 1)
	})
}

func newVCMock(t *testing.T, sound [][]byte) *VCMock {
	is := is.New(t)

	closed := false
	speaking := false
	return &VCMock{
		CloseFunc: func() error {
			t.Log("connection closed")

			closed = true
			return nil
		},
		WriteFunc: func(ctx context.Context, frames [][]byte) error {
			t.Log("sound wrote")

			is.Equal(closed, false)
			is.Equal(speaking, true)
			is.Equal(frames, sound)
			return nil
		},
		SpeakingFunc: func(b bool) error {
			t.Logf("spaking: %v", b)

			speaking = b
			return nil
		},
	}
}

func newVoiceChannelJoinerMock(t *testing.T, vc *VCMock, guildID, channelID string, sound [][]byte) *VoiceChannelJoinerMock {
	is := is.New(t)

	return &VoiceChannelJoinerMock{
		JoinVoiceChannelFunc: func(g, c string) (VC, error) {
			is.Equal(g, guildID)
			is.Equal(c, channelID)

			return vc, nil
		},
	}
}
