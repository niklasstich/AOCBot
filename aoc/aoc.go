// package for handling Advent Of Code data
package aoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/niklasstich/AOCBot/resources"
)

type Leaderboard struct {
	OwnerId int               `json:"owner_id"`
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

const FirstYear = 2015

func LastYearAvailable() int {
	if time.Now().Month() < 12 {
		return time.Now().Year() - 1
	} else {
		return time.Now().Year()
	}
}

func FetchLeaderboard(config *resources.Data, year int) (*Leaderboard, error) {
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
	return &leaderb, nil
}

// return the top x members.
func Top(leaderboard *Leaderboard, x int) ([]Member, error) {
	for key, mem := range leaderboard.Members {
		if mem.Name == "" {
			mem.Name = fmt.Sprintf("Anonymous %s", key)
			leaderboard.Members[key] = mem
		}
	}
	memMap := leaderboard.Members
	memArr := make([]Member, 0)
	for _, v := range memMap {
		memArr = append(memArr, v)
	}
	sort.Slice(memArr, func(i, j int) bool {
		if memArr[i].Stars == memArr[j].Stars {
			return memArr[i].LocalScore > memArr[j].LocalScore
		}
		return memArr[i].Stars > memArr[j].Stars
	})
	if x < len(memArr) {
		return memArr[:x], nil
	}
	return memArr, nil
}
