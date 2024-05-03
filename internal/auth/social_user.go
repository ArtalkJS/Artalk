package auth

import "github.com/markbates/goth"

type SocialUser struct {
	goth.User
	RemoteUID string
	Link      string
}

func GetSocialUser(u goth.User) SocialUser {
	var link string
	if u.Provider == "github" {
		if l, ok := u.RawData["blog"].(string); ok && l != "" {
			link = l
		} else if l, ok := u.RawData["html_url"].(string); ok && l != "" {
			link = l
		}
	}

	// Email patch
	if u.Provider == "steam" {
		// @see https://stackoverflow.com/questions/31571267/steam-get-users-email-address
		u.Email = u.UserID + "@steam.com"
	}
	if u.Email == "" {
		u.Email = u.UserID + "@" + u.Provider + ".com"
	}

	return SocialUser{
		User:      u,
		RemoteUID: u.UserID,
		Link:      link,
	}
}
