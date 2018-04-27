package main

import (
	"fmt"

	"github.com/dags-/discordapp/bot"
	"github.com/dags-/discordapp/command"
	"github.com/dags-/discordapp/util"
	"flag"
)

func main() {
	token := flag.String("token", "", "Discord auth token")
	flag.Parse()
	b := bot.New(token)
	b.AddCommand(command.New("!user <@user> role add <role>", &[]string{"admin"}, addRole))
	b.AddCommand(command.New("!user <@user> role rem <role>", &[]string{"admin"}, remRole))
	b.Connect()
}

func addRole(ctx *command.Context) error {
	user := ctx.Args["user"]
	role := ctx.Args["role"]

	c, e := ctx.Session.Channel(ctx.Message.ChannelID)
	if e != nil {
		return e
	}

	r, e := util.GetRole(ctx.Session, c.GuildID, role)
	if e != nil {
		return e
	}

	e = ctx.Session.GuildMemberRoleAdd(c.GuildID, user, r.ID)
	if e != nil {
		return e
	}

	content := fmt.Sprintf("Added role `%s` to user <@%s>", role, user)
	_, e = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, content)
	return e
}

func remRole(ctx *command.Context) error {
	user := ctx.Args["user"]
	role := ctx.Args["role"]

	c, e := ctx.Session.Channel(ctx.Message.ChannelID)
	if e != nil {
		return e
	}

	r, e := util.GetRole(ctx.Session, c.GuildID, role)
	if e != nil {
		return e
	}

	e = ctx.Session.GuildMemberRoleRemove(c.GuildID, user, r.ID)
	if e != nil {
		return e
	}

	content := fmt.Sprintf("Removed role `%s` from user <@%s>", role, user)
	_, e = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, content)
	return e
}