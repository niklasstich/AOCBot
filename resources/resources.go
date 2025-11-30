package resources

import (
	"os"
)

// data for authentication + configuration
type Data struct {
	LeaderboardID string `json:"channel"`
	BotToken      string `json:"bot_token"`
	SessionToken  string `json:"session_token"`
	InviteCode    string `json:"invite_code"`
}

// read configuration from environment
func Config() (*Data, error) {
	d := Data{}
	d.LeaderboardID = os.Getenv("LEADERBOARDID")
	d.BotToken = os.Getenv("BOT_TOKEN")
	d.SessionToken = os.Getenv("SESSION_TOKEN")
	d.InviteCode = os.Getenv("INVITE_CODE")
	if d.LeaderboardID == "" {
		panic("leaderboardID environment variable missing")
	} else if d.BotToken == "" {
		panic("bot token environment variable missing")
	} else if d.SessionToken == "" {
		panic("aoc session token environment variable missing")
	} else if d.InviteCode == "" {
		panic("invite code environment variable missing")
	}
	return &d, nil
}
