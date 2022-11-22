package handlers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/niklasstich/AOCBot/resources"
	"github.com/niklasstich/AOCBot/svg"

	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/AOCBot/aoc"
)

type imageCacheEntry struct {
	pngPath string
	created time.Time
}

var (
	imageCache map[int]imageCacheEntry
)

func init() {
	imageCache = map[int]imageCacheEntry{}
}

const dayStarFormat = "%3v"

func CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	msgContent := strings.TrimSpace(message.Content)
	//TODO: refactor this with slash commands and have year be an argument, current year as default
	if strings.HasPrefix(msgContent, "/aoc2020") {
		parse(session, message, 2020)
	} else if strings.HasPrefix(msgContent, "/aoc2021") {
		parse(session, message, 2021)
	} else if strings.HasPrefix(msgContent, "/aoc2022") {
		parse(session, message, 2022)
	} else if strings.HasPrefix(msgContent, "/aoc") {
		if time.Now().Month() < 12 { //if it's not yet december, just take last years leaderboard
			parse(session, message, time.Now().Year()-1)
		} else {
			parse(session, message, time.Now().Year())
		}
	}
}

// gets top 200 (or if otherwise specified) members and sends a message highlighting their progress
func parse(session *discordgo.Session, message *discordgo.MessageCreate, year int) {
	d, _ := time.ParseDuration("15m")
	cacheEntry, ok := imageCache[year]
	if ok && time.Since(cacheEntry.created) < d {
		sendPng(cacheEntry.pngPath, session, message, cacheEntry)
		return
	}
	config, _ := resources.Config()
	var err error

	//get leaderboard
	leaderboard, err := aoc.FetchLeaderboard(config, year)
	if err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, "❌"+err.Error())
		return
	}

	sortedMembers, err := aoc.Top(leaderboard, 50)

	if err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, "❌"+err.Error())
		return
	}

	currentTime := time.Now().Format("2006-01-02_15-04-05")
	svgPath := fmt.Sprintf("/out/%s.svg", currentTime)
	pngPath := fmt.Sprintf("/out/%s.png", currentTime)

	err = svg.GenerateSvg(sortedMembers, svgPath)
	if err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, "❌"+err.Error())
		return
	}

	err = ConvertSvgToPng(svgPath, pngPath)
	if err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, "❌"+err.Error())
		return
	}

	cacheEntry = imageCacheEntry{
		pngPath: pngPath,
		created: time.Now(),
	}
	imageCache[year] = cacheEntry

	if err != nil {
		_, _ = session.ChannelMessageSend(message.ChannelID, "❌"+err.Error())
	} else {
		sendPng(pngPath, session, message, cacheEntry)
	}
}

func sendPng(pngPath string, session *discordgo.Session, message *discordgo.MessageCreate, entry imageCacheEntry) {
	png, err := os.Open(pngPath)
	if err != nil {
		panic(err)
	}
	_, _ = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   "leaderboard.png",
				Reader: png,
			},
		},
		Content: fmt.Sprintf("Leaderboard last updated: <t:%d> Next update: <t:%d:R>", 
			entry.created.Unix(), 
			entry.created.Add(15*time.Minute).Unix()),
	})
}

func ConvertSvgToPng(svg string, png string) (err error) {
	command := exec.Command("inkscape",
		"--export-type=png",
		fmt.Sprintf("--export-filename=%s", png),
		"--export-dpi=200",
		fmt.Sprintf("%s", svg))
	err = command.Run()
	return
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
