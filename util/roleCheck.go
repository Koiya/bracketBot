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

func SendRoleCheckMessage(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if !RoleCheck(i) {
		cmd := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have permission to use this command!",
			},
		}
		return s.InteractionRespond(i.Interaction, cmd)
	}
	return nil
}
