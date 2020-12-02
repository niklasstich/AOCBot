package handlers

import (
    "../aoc"
    "../resources"
    "bytes"
    "fmt"
    "github.com/bwmarrin/discordgo"
    "sort"
    "strconv"
    "strings"
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

func formatDays(daysFormat string, startDay int, endDay int) string {
    var out string
    for i := startDay; i <= endDay; i++ {
        dayN := strconv.Itoa(i)
        out += fmt.Sprintf(daysFormat, " "+dayN)
    }
    return out
}

func formatMemberStars(mem *aoc.Member, starFormat string, startDay int, endDay int) string {
    var out string
    for i := startDay; i < endDay; i++ {
        dayKey := strconv.Itoa(i)
        if day, dayOk := mem.CompletionDayLevel[dayKey]; dayOk {
            if len(day) == 2 {
                out += fmt.Sprintf(starFormat, "[*]")
            } else {
                out += fmt.Sprintf(starFormat, "(*)")
            }
        } else {
            out += fmt.Sprintf(starFormat, "")
        }
    }
    return out
}

// format a list of members as string
func format(members []aoc.Member, year int) string {
    strYear := strconv.Itoa(year)
    var buffer bytes.Buffer
    buffer.WriteString(
        "Programmingcord Leaderboard (" + strYear + ") :\n" +
            "```css\n",
    )

    nameScoreFormat := "%-10v%-4v"
    starFormat := "%-3v"

    firstHalf := fmt.Sprintf(
        nameScoreFormat+
            formatDays(starFormat, 1, 13)+
            "\n",
        "", "",
    )
    secondHalf := fmt.Sprintf(
        nameScoreFormat+
            formatDays(starFormat, 13, 25)+
            "\n",
        "", "",
    )

    for _, mem := range members {

        score := strconv.Itoa(mem.Stars)
        name := mem.Name

        if &name == nil || name == "" {
            name = "Anonymous"
        }

        firstHalf += fmt.Sprintf(
            nameScoreFormat+
                formatMemberStars(&mem, starFormat, 1, 13)+
                "\n",
            name,
            "#"+score,
        )
        secondHalf += fmt.Sprintf(
            nameScoreFormat+
                formatMemberStars(&mem, starFormat, 13, 25)+
                "\n",
            name,
            "#"+score,
        )
    }
    buffer.WriteString(firstHalf + secondHalf)
    buffer.WriteString("```")
    return buffer.String()
}
