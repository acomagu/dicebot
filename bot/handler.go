package bot

type ReadyEvent struct{}

type Guild struct {
	Unavailable bool
	Channels    []*Channel
}

type MessageCreateEvent struct {
	Author    *User
	GuildID   string
	ChannelID string
	Content   string
}

type Channel struct {
	ID          string
	IsGuildText bool
}

type GuildCreateEvent struct {
	Guild *Guild
}

type Handler interface {
	OnReady(*ReadyEvent) error
	OnMessageCreate(*MessageCreateEvent) error
	OnGuildCreate(*GuildCreateEvent) error
}
