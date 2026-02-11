package cacheentity

import (
	"time"
)

type HealthcheckNotifEvent struct {
	Event       string    `json:"event"`
	Ts          time.Time `json:"ts"`
	EmailSent   bool      `json:"emailSent"`
	SlackSent   bool      `json:"slackSent"`
	DiscordSent bool      `json:"discordSent"`
}
