package command

import "strings"

type Input struct {
	pos    int
	tokens []string
}

func NewInput(s string) *Input {
	return &Input{
		pos: -1,
		tokens: strings.Split(s, " "),
	}
}

func (i *Input) Reset() {
	i.pos = -1
}

func (i *Input) Size() int {
	return len(i.tokens)
}

func (i *Input) HasNext() bool {
	return i.pos + 1 < len(i.tokens)
}

func (i *Input) Next() string {
	i.pos += 1
	return i.tokens[i.pos]
}
