package handlers

import (
    "bytes"
    "strconv"
    "github.com/bwmarrin/discordgo"
    "strings"
    "../aoc"
    "../resources"
    "sort"
    "time"
)

func CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
    msgContent := strings.TrimSpace(message.Content)
    parts := strings.Split(msgContent, " ")
    if strings.HasPrefix(msgContent, "/aoc2015") {
        parse(session, message, 2015, parts)
    } else if strings.HasPrefix(msgContent, "/aoc2016") {
        parse(session, message, 2016, parts)
    } else if strings.HasPrefix(msgContent, "/aoc2017") {
        parse(session, message, 2017, parts)
    } else if strings.HasPrefix(msgContent, "/aoc") {
        parse(session, message, time.Now().Year(), parts)
    }
}

func parse(session *discordgo.Session, message *discordgo.MessageCreate, year int, parts []string) {
    config, _ := resources.Config()
    if len(parts) <= 1 {
        session.ChannelMessageSend(message.ChannelID, format(top(config, year, 200), year)) // 200 = max members
    } else {
        topAmount,_ := strconv.Atoi(parts[1])
        session.ChannelMessageSend(message.ChannelID, format(top(config, year, topAmount), year))
    }
}

// return the top x members. 
func top(config *resources.Data, year int, x int) []aoc.Member{
    memMap := aoc.FetchLeaderboard(config, year).Members
    memArr := make([]aoc.Member, 0)
    for _,v := range memMap {
        memArr = append(memArr, v)
    }
    sort.Slice(memArr, func(i, j int) bool {
        return memArr[i].Stars > memArr[j].Stars
    })
    if x < len(memArr) {
        return memArr[:x]
    }
    return memArr
}


func format(members []aoc.Member, year int) string {
    // return a list of members as string
    strYear := strconv.Itoa(year)
    var buffer bytes.Buffer
    buffer.WriteString("Programmingcord Leaderboard (" + strYear + ") :\n================\nStars:\n")
    for _, mem := range members {
        score := strconv.Itoa(mem.Stars)
        name := mem.Name
        if &name == nil || name == "" {
            name = "Anonymous"
        }
        buffer.WriteString(name + ": " + score + "!\n")
    }
    return buffer.String()
}
