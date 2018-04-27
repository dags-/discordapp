package command

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"sort"
	"bytes"
	"github.com/dags-/discordapp/util"
	"github.com/pkg/errors"
	"fmt"
)

type Manager struct {
	commands []*Command
}

func NewManager() *Manager {
	mgr := &Manager{}
	mgr.Add(New("!help <command...>", nil, mgr.help()))
	return mgr
}

// adds a command to the manager
func (mgr *Manager) Add(c *Command) *Manager {
	mgr.commands = append(mgr.commands, c)
	return mgr
}

// invokes a command based on the input message
func (mgr *Manager) Invoke(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var errs []error

	for _, c := range mgr.commands {
		ex, err := invoke(s, m, c)
		if ex {
			return err
		}

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		// return last error
		return errs[len(errs)-1]
	}

	// should only occur if no matching root commands
	return nil
}

// the !help executor
func (mgr *Manager) help() Executor {
	var ex Executor = func(ctx *Context) error {
		term := ctx.Args["command"]

		var lines []string
		for _, c := range mgr.commands {
			if strings.HasPrefix(c.Usage, term) {
				lines = append(lines, c.Usage)
			}
		}

		sort.Strings(lines)

		var buf bytes.Buffer
		for i := 0; i < len(lines); i++ {
			l := lines[i]
			buf.WriteString(l)
			if i < len(lines)-1 {
				buf.WriteRune('\n')
			}
		}

		_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, buf.String())
		return err
	}
	return ex
}

// attempts to invoke the given command, returning true if the executor was invoked & any errors thrown
func invoke(s *discordgo.Session, m *discordgo.MessageCreate, c *Command) (bool, error) {
	if c.Roles != nil && !hasPermission(s, m, *c.Roles) {
		return false, errors.New("no permission")
	}

	in := NewInput(m.Content)
	ctx := NewCtx(s, m)
	err := c.Parse(in, ctx)

	if in.pos == 0 {
		// no matching root command
		return false, nil
	}

	if err != nil {
		// command failed
		return false, err
	}

	// command passed
	return true, (*c.Executor)(ctx)
}

// check if the author of the message has one of the given roles
func hasPermission(s *discordgo.Session, m *discordgo.MessageCreate, roles []string) bool {
	c, e := s.Channel(m.ChannelID)
	if e != nil {
		fmt.Println("get channel error", e)
		return false
	}
	return util.HasAnyRole(s, c.GuildID, m.Author.ID, roles...)
}