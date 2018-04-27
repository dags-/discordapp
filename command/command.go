package command

import (
	"strings"
)

type Executor func(ctx *Context) error

type Command struct {
	Usage    string
	Roles    *[]string
	Nodes    []*Node
	Executor *Executor
}

// Creates a new Command
//  - pattern 	- the pattern of the command input (see below)
//  - roles		- if not nil, the user must have one of these roles to use the command
//  - executor 	- the callback that is executed if the command successfully parses
//
// Pattern components:
//  - root command: `!command`
//  - parameters:
//    - sub: 		`name`		- the argument at this position must equal 'name'
//    - single: 	`<name>`	- the argument at this position will be the value of 'name'
//    - remaining: 	`<name...>` - the arguments after this position will be joined to form the value of 'name'
//    - user: 		`<@name>` 	- the user @mentioned at this position will be the value of 'name'
//    - bot:       	`<@bot>		- require the bot to be mentioned at this position in the command
//
// Examples:
// pattern: `!help <command>`, usage: `!help user` (show help about the '!user' command
// pattern: `!user <@user> slap`, usage: `!user @dags slap` (slap dags)
// pattern: `!user <@user> tell <message...>`, usage: `!user @dags tell yo, what's up bro! (message dags with 'yo, what's up bro!')
func New(pattern string, roles *[]string, executor Executor) *Command {
	pattern = trimPattern(pattern)

	var in = NewInput(pattern)
	var nodes []*Node

	for in.HasNext() {
		n, e := NewNode(in)
		if e != nil {
			return nil
		}
		nodes = append(nodes, n)
	}

	return &Command{Nodes: nodes, Usage: pattern, Roles: roles, Executor: &executor}
}

func (cmd *Command) Parse(i *Input, c *Context) (error) {
	for _, n := range cmd.Nodes {
		e := n.Parser(n, i, c)
		if e != nil {
			return e
		}
	}
	return nil
}

func trimPattern(s string) string {
	if strings.HasPrefix(s, "!") {
		return strings.TrimPrefix(s, "!")
	}
	return s
}
