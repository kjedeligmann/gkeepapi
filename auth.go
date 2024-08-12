package gkeepapi

import (
	"fmt"

	"github.com/kjedeligmann/gpsoauth"
)

const OAuthScopes string = "oauth2:https://www.googleapis.com/auth/memento https://www.googleapis.com/auth/reminders"

type Auth struct {
	email       string
	gaid        string
	masterToken string
	accessToken string
}

func (s *Auth) Load(email, gaid, masterToken string) error {
	s.email = email
	s.gaid = gaid
	s.masterToken = masterToken
	if err := s.Refresh(); err != nil {
		return err
	}
	return nil
}

func (s *Auth) Refresh() error {
	resp, err := gpsoauth.PerformOAuthWithDefaults(s.email, s.masterToken, s.gaid, OAuthScopes, "com.google.android.keep")
	if err != nil {
		return err
	}
	if token, ok := resp["Auth"]; ok {
		s.accessToken = token
		return nil
	}
	return fmt.Errorf("No token was fetched")
}
