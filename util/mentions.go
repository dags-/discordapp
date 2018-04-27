package util

import "github.com/bwmarrin/discordgo"

func Mentions(m *discordgo.MessageCreate, id string) bool {
	for _, u := range m.Mentions {
		if u.ID == id {
			return true
		}
	}
	return false
}