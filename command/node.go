package command

import (
	"bytes"
	"strings"

	"github.com/pkg/errors"
)

type Node struct {
	Name   string
	Parser func(n *Node, input *Input, ctx *Context) error
}

func (n *Node) String() string {
	return n.Name
}

func NewNode(input *Input) (*Node, error) {
	if !input.HasNext() {
		return nil, errors.New("not enough args")
	}

	name := input.Next()
	parser := sub

	if strings.HasPrefix(name, "<") && strings.HasSuffix(name, ">") {
		name = strings.Trim(name, "<>")
		parser = single
		if name == "@bot" {
			parser = bot
		} else if strings.HasPrefix(name, "@") {
			name = strings.TrimPrefix(name, "@")
			parser = user
		} else if strings.HasSuffix(name, "...") {
			name = strings.TrimSuffix(name, "...")
			parser = remaining
		}
	}

	return &Node{Name: name, Parser: parser}, nil
}

func single(n *Node, input *Input, ctx *Context) error {
	if !input.HasNext() {
		return errors.New("not enough args")
	}

	ctx.Args[n.Name] = input.Next()
	return nil
}

func remaining(n *Node, input *Input, ctx *Context) error {
	if !input.HasNext() {
		return errors.New("not enough args")
	}

	buf := bytes.Buffer{}
	for input.HasNext() {
		buf.WriteString(input.Next())
		if input.HasNext() {
			buf.WriteRune(' ')
		}
	}

	ctx.Args[n.Name] = buf.String()
	return nil
}

func user(n *Node, input *Input, ctx *Context) error {
	if !input.HasNext() {
		return errors.New("not enough args")
	}

	user := input.Next()
	if strings.HasPrefix(user, "<@") && strings.HasSuffix(user, ">") {
		user = strings.Trim(user, "<@>")
		ctx.Args[n.Name] = user
		return nil
	}

	return errors.New("user not mentioned")
}

func bot(n *Node, input *Input, ctx *Context) error {
	if !input.HasNext() {
		return errors.New("not enough args")
	}

	bot := input.Next()
	if strings.HasPrefix(bot, "<@") && strings.HasSuffix(bot, ">") {
		bot = strings.Trim(bot, "<@>")
		if bot != ctx.Session.State.User.ID {
			return errors.New("input does not @mention the bot")
		}
	}

	return nil
}

func sub(n *Node, input *Input, ctx *Context) error {
	if !input.HasNext() {
		return errors.New("not enough args")
	}

	name := input.Next()
	if name != n.Name {
		return errors.New("invalid arg: " + name + ", expected: " + n.Name)
	}

	return nil
}