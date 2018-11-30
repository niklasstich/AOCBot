package resources

import (
    "io/ioutil"
    "encoding/json"
)


// path relative to main
const datapath = "../resources/data.json"

// data for authentication + config (edit data.json to change)
type Data struct {
    Channel string `json:"channel"`
    BotToken string `json:"bot_token"`
    SessionToken string `json:"session_token"`
}


// read the config file
func Config() (*Data, error) {
    d := Data{}
    bytes, err := ioutil.ReadFile(datapath)
    if err != nil {
        return &d, err
    }
    json.Unmarshal(bytes, &d)
    return &d, nil
}
