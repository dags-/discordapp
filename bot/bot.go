package bot

import (
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/dags-/discordapp/command"
	"github.com/pkg/errors"
)

type Bot struct {
	cmdLock  sync.RWMutex
	session  *discordgo.Session
	commands *command.Manager
}

func New(auth *string) *Bot {
	if auth == nil || *auth == "" {
		panic(errors.New("invalid token"))
	}

	s, e := discordgo.New("Bot " + *auth)
	if e != nil {
		panic(e)
	}

	return &Bot{session: s, commands: command.NewManager()}
}

func (b *Bot) Connect() error {
	b.session.AddHandler(b.cmdListener)

	e := b.session.Open()
	if e != nil {
		return e
	}

	c := make(chan os.Signal, 1)
	<-c

	b.session.Close()
	b.session = nil
	b.commands = nil

	return nil
}

func (b *Bot) AddCommand(c *command.Command) *Bot {
	b.cmdLock.Lock()
	defer b.cmdLock.Unlock()

	if c != nil {
		b.commands.Add(c)
	}

	return b
}

func (b *Bot) AddListener(hndlr ...interface{}) *Bot {
	for _, h := range hndlr {
		b.session.AddHandler(h)
	}
	return b
}

func (b *Bot) cmdListener(s *discordgo.Session, m *discordgo.MessageCreate) {
	b.cmdLock.Lock()
	defer b.cmdLock.Unlock()

	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	m.Content = strings.TrimPrefix(m.Content, "!")
	err := b.commands.Invoke(s, m)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
	}
}
