package domain

import (
	"sort"
	"time"
)

type Session struct {
	UserId       int
	RefreshToken string
	ExpiresAt    time.Time
}

func SortSessionsByTime(sessions []Session) {
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].ExpiresAt.Unix() > sessions[j].ExpiresAt.Unix()
	})
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
