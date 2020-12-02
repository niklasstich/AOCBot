package resources

import (
    "os"
)

// data for authentication + config (edit data.json to change)
type Data struct {
    Channel string `json:"channel"`
    BotToken string `json:"bot_token"`
    SessionToken string `json:"session_token"`
}


// read the config file
func Config() (*Data, error) {
    d := Data{}
    d.Channel = os.Getenv("LEADERBOARDID")
    d.BotToken = os.Getenv("BOT_TOKEN")
    d.SessionToken = os.Getenv("SESSION_TOKEN")
    if d.Channel == "" || d.BotToken == "" || d.SessionToken == "" {
        panic("check environment variables")
    }
    return &d, nil
}
