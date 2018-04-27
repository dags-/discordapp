package command

import "github.com/bwmarrin/discordgo"

type Context struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
	Args    map[string]string
}

func NewCtx(s *discordgo.Session, m *discordgo.MessageCreate) *Context {
	return &Context{Session: s, Message: m, Args: make(map[string]string)}
}
