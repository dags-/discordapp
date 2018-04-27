package util

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"fmt"
)

func GetRole(s *discordgo.Session, guild, role string) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(guild)
	if err != nil {
		return nil, err
	}
	for _, r := range roles {
		if r.Name == role {
			return r, nil
		}
	}
	return nil, errors.New("role not found")
}

func HasAnyRole(s *discordgo.Session, guild, user string, roles ...string) bool {
	m, e := s.GuildMember(guild, user)
	if e != nil {
		fmt.Println("get member error:", e)
		return false
	}

	for _, r := range roles {
		role, err := GetRole(s, guild, r)
		if err != nil {
			continue
		}
		if indexOf(role.ID, m.Roles) != -1 {
			return true
		}
	}

	return false
}

func HasAllRoles(s *discordgo.Session, guild, user string, roles ...string) bool {
	m, e := s.GuildMember(guild, user)
	if e != nil {
		fmt.Println("get member error:", e)
		return false
	}

	for _, r := range roles {
		role, err := GetRole(s, guild, r)
		if err != nil {
			continue
		}
		if indexOf(role.ID, m.Roles) == -1 {
			return false
		}
	}

	return true
}

func indexOf(s string, in []string) int {
	for i := 0; i < len(in); i++ {
		if s == in[i] {
			return i
		}
	}
	return -1
}