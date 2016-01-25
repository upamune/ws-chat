package main

import "errors"

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: アバターのURLを返します")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	GetAvatarURL(*client) (string, error)
}

// AuthAvatar implemet Avatar
type AuthAvatar struct{}

// GravatarAvatar is Gravatar's Avatar
type GravatarAvatar struct{}

// UseGravatar is using Gravatar
var UseGravatar GravatarAvatar

// UseAuthAvatar is AuthAvatar
var UseAuthAvatar AuthAvatar

// GetAvatarURL is returns string and error
func (a AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// GetAvatarURL is returns GravatarURL and error
func (g GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
