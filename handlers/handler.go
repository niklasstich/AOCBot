package handlers

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/AOCBot/aoc"
	"github.com/niklasstich/AOCBot/resources"
	log "github.com/sirupsen/logrus"
)

const dayStarFormat = "%3v"

func CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	msgContent := strings.TrimSpace(message.Content)
	parts := strings.Split(msgContent, " ")
	//TODO: refactor this with slash commands and have year be an argument, current year as default
	if strings.HasPrefix(msgContent, "/aoc2015") {
		parse(session, message, 2015, parts)
	} else if strings.HasPrefix(msgContent, "/aoc2016") {
		parse(session, message, 2016, parts)
	} else if strings.HasPrefix(msgContent, "/aoc2017") {
		parse(session, message, 2017, parts)
	} else if strings.HasPrefix(msgContent, "/aoc2020") {
		parse(session, message, 2020, parts)
	} else if strings.HasPrefix(msgContent, "/aoc") {
		if time.Now().Month() < 12 { //if it's not yet december, just take last years leaderboard
			parse(session, message, time.Now().Year()-1, parts)
		} else {
			parse(session, message, time.Now().Year(), parts)
		}
	}
}

//gets top 200 (or if otherwise specified) members and sends a message highlighting their progress
func parse(session *discordgo.Session, message *discordgo.MessageCreate, year int, parts []string) {
	config, _ := resources.Config()
	var topMem []aoc.Member
	var err error
	//get guild name
	var guildName string
	guild, err := session.Guild(message.GuildID)
	if err != nil {
		log.Errorf("Failed to get guild for GuildID %v: %v", message.GuildID, err)
		guildName = "unknown"
	} else {
		guildName = guild.Name
	}

	if len(parts) <= 1 {
		topMem, err = top(config, year, 200)
	} else {
		topAmount, _ := strconv.Atoi(parts[1])
		topMem, err = top(config, year, topAmount)
	}
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "❌"+err.Error())
	} else {
		session.ChannelMessageSend(message.ChannelID, format(topMem, year, guildName))
	}
}

// return the top x members.
func top(config *resources.Data, year int, x int) ([]aoc.Member, error) {
	lb, err := aoc.FetchLeaderboard(config, year)
	if err != nil {
		return nil, err
	}
	memMap := lb.Members
	memArr := make([]aoc.Member, 0)
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

func formatDays(startDay int, endDay int) string {
	var out string
	for i := startDay; i <= endDay; i++ {
		if i == startDay {
			out += fmt.Sprintf("%2v", strconv.Itoa(i))
		} else {
			out += fmt.Sprintf(dayStarFormat, strconv.Itoa(i))
		}
	}
	return out
}

func formatMemberStars(mem aoc.Member, startDay int, endDay int) string {
	var out string
	for i := startDay; i <= endDay; i++ {
		dayKey := strconv.Itoa(i)
		if day, dayOk := mem.CompletionDayLevel[dayKey]; dayOk {
			switch len(day) {
			case 2:
				out += "[*]"
			case 1:
				out += "(*)"
			default: //len is either 0 or greater than 2, both of which make no sense, therefore space
				out += "   "
			}
		} else { //if !dayOk, then there was no entry for the day in the map, we assume no stars were gotten
			out += "   "
		}
	}
	return out
}

func formatLeaderboard(members *[]aoc.Member, w io.Writer) {
	const padding = 3
	tw := tabwriter.NewWriter(w, 0, 0, padding, ' ', 0)
	dayRanges := [][]int{{1, 14}, {15, 25}}
	for _, days := range dayRanges {
		startDay, endDay := days[0], days[1]
		fmt.Fprintln(tw, "\t\t"+formatDays(startDay, endDay)+"\t") // header
		for _, mem := range *members {
			score := strconv.Itoa(mem.LocalScore)
			fmt.Fprintln(tw, mem.Name+"\t#"+score+"\t"+formatMemberStars(mem, startDay, endDay)+"\t")
		}
	}
	tw.Flush()
}

// format a list of members as string
func format(members []aoc.Member, year int, guildName string) string {
	strYear := strconv.Itoa(year)
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s Leaderboard (Year %s)\n```css\n", guildName, strYear))
	formatLeaderboard(&members, &buffer)
	buffer.WriteString("```")
	return buffer.String()
}
