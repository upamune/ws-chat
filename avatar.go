package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: アバターのURLを返します")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	GetAvatarURL(ChatUser) (string, error)
}

// TryAvatars is Slice Avatar
type TryAvatars []Avatar

// AuthAvatar implemet Avatar
type AuthAvatar struct{}

// GravatarAvatar is Gravatar's Avatar
type GravatarAvatar struct{}

// FileSystemAvatar is User Upload Avatar
type FileSystemAvatar struct{}

// UseAuthAvatar is AuthAvatar
var UseAuthAvatar AuthAvatar

// UseGravatar is using Gravatar
var UseGravatar GravatarAvatar

// UseFileSystemAvatar is using User Upload Avatar
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL is returns string and error
func (a AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

// GetAvatarURL is returns GravatarURL and error
func (g GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// GetAvatarURL is returns local url and error
func (f FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}

// GetAvatarURL try to get some avatar types
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}
