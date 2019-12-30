package sapphire

import (
	"github.com/bwmarrin/discordgo"
)

// Utility to help calculate permissions. Since discordgo is too damn low-level

// Permissions represent permission bits for a discord entity.
type Permissions int

// PermissionsForRole returns a permissions instance for a role.
func PermissionsForRole(role *discordgo.Role) Permissions {
	return Permissions(role.Permissions)
}

// PermissionsForMember returns a permissions instance for a member.
func PermissionsForMember(guild *discordgo.Guild, member *discordgo.Member) Permissions {
	// Owners have all permissions.
	if member.User.ID == guild.OwnerID {
		return Permissions(discordgo.PermissionAll)
	}
	bits := 0
	// Combine all permissions from every role.
	for _, rID := range member.Roles {
		var role *discordgo.Role
		for _, gRole := range guild.Roles {
			if gRole.ID == rID {
				role = gRole
				break
			}
		}
		// Do we have a choice if it was nil?
		if role != nil {
			bits |= role.Permissions
		}
	}
	return Permissions(bits)
}

func (perms Permissions) Has(bits int) bool {
	return (int(perms) & bits) == bits
}
