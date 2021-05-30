// package for handling Advent Of Code data
package aoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/niklasstich/AOCBot/resources"
)

var (
	lastHit         time.Time
	lastLeaderboard *Leaderboard
)

type Leaderboard struct {
	OwnerId string            `json:"owner_id"`
	Event   string            `json:"event"`
	Members map[string]Member `json:"members"`
}

type Member struct {
	Name               string                            `json:"name"`
	LocalScore         int                               `json:"local_score"`
	GlobalScore        int                               `json:"global_score"`
	Stars              int                               `json:"stars"`
	CompletionDayLevel map[string]map[string]interface{} `json:"completion_day_level"`
}

func FetchLeaderboard(config *resources.Data, year int) (*Leaderboard, error) {
	//Limit requests to leaderboard to every 15 minutes
	d, _ := time.ParseDuration("15m")
	if time.Now().Sub(lastHit) < d {
		return lastLeaderboard, nil
	}
	fmt.Println("fetching data..")
	url := "https://adventofcode.com/" + strconv.Itoa(year) + "/leaderboard/private/view/" + config.LeaderboardID + ".json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating http request, %v", err)
	}
	req.Header.Set("Cookie", "session="+config.SessionToken)
	client := &http.Client{}
	response, e := client.Do(req)
	if e != nil {
		return nil, fmt.Errorf("failed filling http request to adventofcode.com, %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to server failed with status %d, header: %s", response.StatusCode, response.Status)
	}
	//make a copy of the stream and feed that to the decoder, so we can preserve body for error handling
	var buf bytes.Buffer
	tee := io.TeeReader(response.Body, &buf)
	leaderb := Leaderboard{}
	err = json.NewDecoder(tee).Decode(&leaderb)
	if err != nil {
		return nil, fmt.Errorf("failed decoding responce from adventofcode.com, %v\n<@136512985542819840> fix your session token babyrage", err)
	}
	lastLeaderboard = &leaderb
	return &leaderb, nil
}
