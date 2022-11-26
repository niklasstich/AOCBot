package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/niklasstich/AOCBot/resources"
	"github.com/niklasstich/AOCBot/svg"

	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/AOCBot/aoc"
)

type imageCacheEntry struct {
	pngPath   string
	createdAt time.Time
}

const commandCD = 1 * time.Minute // TODO: maybe cooldown should be based on action (e.g. image more cd than invalid command)
const cacheDuration = 15 * time.Minute

var (
	imageCache              = map[int]imageCacheEntry{}
	userIDToLastCommandTime = map[string]time.Time{}
)

func CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := strings.TrimSpace(m.Content)

	// Ignore own messages and those that don't start with "/aoc"
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(msg, "/aoc") {
		return
	}

	// User based cooldown
	if time.Since(userIDToLastCommandTime[m.Author.ID]) < commandCD {
		return // Too early
	} else {
		userIDToLastCommandTime[m.Author.ID] = time.Now()
	}

	if msg == "/aoc-info" {
		replyWithInfo(s, m.ChannelID)
		return
	}

	msgParts := strings.Split(msg, " ")
	if len(msgParts) == 1 && msgParts[0] == "/aoc" {
		// No year -> use year from the last event
		handleLeaderboardCommand(s, m.ChannelID, aoc.LastYearAvailable())
	} else if len(msgParts) == 2 && msgParts[0] == "/aoc" {
		year, err := strconv.Atoi(msgParts[1])
		if err != nil || year < aoc.FirstYear || year > time.Now().Year() {
			replyWithInvalidYear(s, m.ChannelID)
		} else {
			handleLeaderboardCommand(s, m.ChannelID, year)
		}
	} else {
		// Doesn't fit /aoc-info or /aoc <year> format
		replyWithUnrecognizedCommand(s, m.ChannelID)
	}
}

func handleLeaderboardCommand(s *discordgo.Session, channelID string, year int) {
	cacheEntry, ok := imageCache[year]
	if ok && time.Since(cacheEntry.createdAt) < cacheDuration {
		replyWithLeaderboardImage(s, channelID, cacheEntry)
		return
	}

	currentTimeString := time.Now().Format("2006-01-02_15-04-05")
	config, err := resources.Config()
	if err != nil {
		replyWithError(s, channelID, err)
		return
	}

	leaderboard, err := aoc.FetchLeaderboard(config, year)
	if err != nil {
		replyWithError(s, channelID, err)
		return
	}

	sortedMembers, err := aoc.Top(leaderboard, 50)
	if err != nil {
		replyWithError(s, channelID, err)
		return
	}

	svgPath := fmt.Sprintf("/out/%s.svg", currentTimeString)
	err = svg.GenerateSvg(year, sortedMembers, svgPath)
	if err != nil {
		replyWithError(s, channelID, err)
		return
	}

	pngPath := fmt.Sprintf("/out/%s.png", currentTimeString)
	err = svg.ConvertSvgToPng(svgPath, pngPath)
	if err != nil {
		replyWithError(s, channelID, err)
		return
	}

	// Save in cache
	cacheEntry = imageCacheEntry{
		pngPath:   pngPath,
		createdAt: time.Now(),
	}
	imageCache[year] = cacheEntry

	replyWithLeaderboardImage(s, channelID, cacheEntry)
}

func replyWithInfo(s *discordgo.Session, channelID string) {
	s.ChannelMessageSend(
		channelID,
		"WTF is this? https://adventofcode.com/about\n\n"+
			"Join our leaderboard at https://adventofcode.com/leaderboard/private with the code: `784176-b767a0f2`\n"+
			"Type `/aoc <year>` (e.g. `/aoc 2022`) to see the standings",
	)
}

func replyWithInvalidYear(s *discordgo.Session, channelID string) {
	s.ChannelMessageSend(
		channelID,
		fmt.Sprintf("`/aoc <year>` needs a valid year in the range [%d..%d]", aoc.FirstYear, time.Now().Year()),
	)
}

func replyWithUnrecognizedCommand(s *discordgo.Session, channelID string) {
	s.ChannelMessageSend(
		channelID,
		"Unrecognized command. Try `/aoc-info` or `/aoc <year>` (e.g. `/aoc 2022`)",
	)
}

func replyWithError(s *discordgo.Session, channelID string, err error) {
	s.ChannelMessageSend(
		channelID,
		fmt.Sprintf("❌ %s", err.Error()),
	)
}

func replyWithLeaderboardImage(session *discordgo.Session, channelID string, imageEntry imageCacheEntry) {
	png, err := os.Open(imageEntry.pngPath)
	if err != nil {
		panic(err)
	}
	session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   "leaderboard.png",
				Reader: png,
			},
		},
		Content: fmt.Sprintf(
			"Updated <t:%d:R> - Next update <t:%d:R>",
			imageEntry.createdAt.Unix(),
			imageEntry.createdAt.Add(cacheDuration).Unix(),
		),
	})
}
