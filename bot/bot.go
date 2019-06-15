package bot

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/acomagu/dicebot/soundplayer"
	"github.com/pkg/errors"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type DB interface {
	AddGuild(*Guild) error
}

type Bot struct {
	Session      Session
	SoundFrames  [][]byte
	SoundPlayers map[string]*soundplayer.SoundPlayer
	Joiner       soundplayer.VoiceChannelJoiner
	DB           DB
	me           *User
}

func New(session Session, joiner soundplayer.VoiceChannelJoiner, db DB, soundFrames [][]byte) (*Bot, error) {
	me, err := session.GetMe()
	if err != nil {
		return nil, errors.Wrap(err, "could not get own User info")
	}

	return &Bot{
		Session:      session,
		SoundFrames:  soundFrames,
		SoundPlayers: make(map[string]*soundplayer.SoundPlayer),
		Joiner:       joiner,
		me:           me,
	}, nil
}

func (b *Bot) Wait() {
	// Wait forever.
	<-(chan struct{})(nil)
}

func (b *Bot) OnReady(event *ReadyEvent) error {
	b.Session.UpdateStatus(0, "!dice")
	return nil
}

func (b *Bot) OnGuildCreate(event *GuildCreateEvent) error {
	if event.Guild.Unavailable {
		return nil
	}

	var channel *Channel
	for _, c := range event.Guild.Channels {
		if c.IsGuildText && b.Session.HasSendMessagePermission(b.me.ID, c.ID) {
			channel = c
			break
		}
	}
	if channel == nil {
		return nil
	}

	_, err := b.Session.ChannelMessageSend(channel.ID, "DICEBOTの準備ができました! ボイスチャンネルに参加して\"!dice 2D10\"などとタイプしてみましょう!")
	if err != nil {
		return err
	}

	if err := b.DB.AddGuild(&Guild{
		ID:event.Guild.ID,
	}); err != nil {
		return err
	}

	return nil
}

func (b *Bot) OnMessageCreate(e *MessageCreateEvent) error {
	if e.Author.ID == b.me.ID {
		return nil
	}

	channelID := e.ChannelID
	guildID := e.GuildID

	if !strings.HasPrefix(e.Content, "!dice") {
		return nil
	}

	vss, err := b.Session.VoiceStates(guildID)
	if err != nil {
		return err
	}

	if _, ok := b.SoundPlayers[guildID]; !ok {
		b.SoundPlayers[guildID] = soundplayer.NewSoundPlayer(b.Joiner, guildID)
	}
	player := b.SoundPlayers[guildID]

	var voiceChannelID string
	for _, vs := range vss {
		if vs.UserID == e.Author.ID {
			voiceChannelID = vs.ChannelID
		}
	}
	if channelID == "" {
		return nil
	}

	// Response
	answer, ok := calc(e.Content)
	if !ok {
		_, err := b.Session.ChannelMessageSend(e.ChannelID, help)
		return err
	}

	if err := player.PlaySound(context.TODO(), voiceChannelID, b.SoundFrames); err != nil {
		return err
	}

	if _, err := b.Session.ChannelMessageSend(e.ChannelID, answer); err != nil {
		return err
	}

	return nil
}

var diceregexp = regexp.MustCompile(`(\d+)[dD](\d+)`)

var help = `?? よくわかりません :sweat_drops: "!dice 2D10"などとタイプしてみてください :bow:`

func calc(msg string) (string, bool) {
	groups := diceregexp.FindStringSubmatch(msg)
	if groups == nil || len(groups) < 3 {
		return "", false
	}

	n, err := strconv.Atoi(groups[1])
	if err != nil {
		return "", false
	}

	h, err := strconv.Atoi(groups[2])
	if err != nil {
		return "", false
	}

	ans := 0
	for i := 0; i < n; i++ {
		ans += rand.Intn(h) + 1
	}
	return fmt.Sprint(ans), true
}
