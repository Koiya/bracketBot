package util

import "github.com/bwmarrin/discordgo"

// RoleCheck Check the role of the user before performing action
func RoleCheck(i *discordgo.InteractionCreate) bool {
	var result = false
	for _, value := range i.Member.Roles {
		if value == GetTOML("Bot.ModRole_ID") {
			result = true
		}
	}
	return result
}
