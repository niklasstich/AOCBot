package handlers

import (
    "../aoc"
    "../resources"
    "bytes"
    "fmt"
    "github.com/bwmarrin/discordgo"
    "io"
    "sort"
    "strconv"
    "strings"
    "text/tabwriter"
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

func parse(session *discordgo.Session, message *discordgo.MessageCreate, year int, parts []string) {
    config, _ := resources.Config()
    if len(parts) <= 1 {
        session.ChannelMessageSend(message.ChannelID, format(top(config, year, 200), year)) // 200 = max members
    } else {
        topAmount, _ := strconv.Atoi(parts[1])
        session.ChannelMessageSend(message.ChannelID, format(top(config, year, topAmount), year))
    }
}

// return the top x members.
func top(config *resources.Data, year int, x int) []aoc.Member {
    memMap := aoc.FetchLeaderboard(config, year).Members
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
        return memArr[:x]
    }
    return memArr
}

func formatDays(startDay int, endDay int) string {
    var out string
    for i := startDay; i <= endDay; i++ {
        dayN := strconv.Itoa(i)
        out += " " + dayN + " "
    }
    return out
}

func formatMemberStars(mem aoc.Member, startDay int, endDay int) string {
    var out string
    for i := startDay; i < endDay; i++ {
        dayKey := strconv.Itoa(i)
        if day, dayOk := mem.CompletionDayLevel[dayKey]; dayOk {
            if len(day) == 2 {
                out += "[*]"
            } else {
                out += "(*)"
            }
        } else {
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
func format(members []aoc.Member, year int) string {
    strYear := strconv.Itoa(year)
    var buffer bytes.Buffer
    buffer.WriteString(
        "Programmingcord Leaderboard (" + strYear + "):\n" +
            "```css\n",
    )
    formatLeaderboard(&members, &buffer)
    buffer.WriteString("```")
    return buffer.String()
}
