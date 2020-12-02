// package for handling Advent Of Code data
package aoc

import (
    "fmt"
    "net/http"
    "encoding/json"
    "strconv"
    "../resources"
    "time"
)

var (
    lastHit time.Time
    lastLeaderboard Leaderboard
)
type Leaderboard struct {
    OwnerId string `json:"owner_id"`
    Event string `json:"event"`
    Members map[string]Member `json:"members"`
}

type Member struct {
    Name string `json:"name"`
    LocalScore int `json:"local_score"`
    GlobalScore int `json:"global_score"`
    Stars int `json:"stars"`
}

func FetchLeaderboard(config *resources.Data, year int) *Leaderboard {
    //Limit requests to leaderboard to every 15 minutes
    d, _ := time.ParseDuration("15m")
    if time.Now().Sub(lastHit) < d {
        return &lastLeaderboard
    }
    fmt.Println("fetching data..")
    url := "https://adventofcode.com/" + strconv.Itoa(year) + "/leaderboard/private/view/" + config.Channel + ".json"
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        panic(err)
    }
    req.Header.Set("Cookie", "session=" + config.SessionToken)
    client := &http.Client{}
    response, e := client.Do(req)
    if e != nil {
        panic(e)
    }
    defer response.Body.Close()
    leaderb := Leaderboard{}
    err = json.NewDecoder(response.Body).Decode(&leaderb)
    if err != nil {
        panic(err)
    }
    lastLeaderboard = leaderb
    return &leaderb
}
