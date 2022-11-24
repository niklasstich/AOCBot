package svg

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"

	svg "github.com/ajstarks/svgo"
	"github.com/niklasstich/AOCBot/aoc"
)

var (
	nameWidth = 140
	starGap   = 22
)

func GenerateSvg(year int, members []aoc.Member, filepath string) (err error) {
	fmt.Println("Generating new image...")
	width := 720
	height := 60 + 40 + len(members)*20
	file, err := os.Create(filepath)

	if err != nil {
		return
	}
	defer file.Close()
	canvas := svg.New(file)
	canvas.Start(width, height)

	//header and lines
	canvas.Rect(0, 0, width, height, "fill: rgb(15,15,35);")
	canvas.Text(width/2, 30,
		"Advent of Code Leaderboard",
		"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
			"dominant-baseline:middle; text-anchor:middle;"+
			"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00;")
	//Align //Year right underneath Advent
	canvas.Text(301, 50,
		"//",
		"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
			"dominant-baseline:middle; text-anchor:start;"+
			"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00; opacity: 0.4; letter-spacing: -0.05em;")
	canvas.Text(320, 50,
		fmt.Sprintf("%d", year),
		"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
			"dominant-baseline:middle; text-anchor:start;"+
			"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00;")
	canvas.Line(0, 80, width, 80, "stroke: rgb(0, 99, 0); stroke-width: 2;")
	canvas.Line(nameWidth, 80, nameWidth, height, "stroke: rgb(0, 99, 0); stroke-width: 2;")

	//numbers for stars
	for i := 1; i < 26; i++ {
		canvas.Text(nameWidth+starGap*i, 70, fmt.Sprint(i), "fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 12px;"+
			"dominant-baseline:middle; text-anchor:middle;")
	}

	PrintMembers(canvas, members)

	canvas.End()
	return
}

func PrintMembers(canvas *svg.SVG, members []aoc.Member) {
	i := 0
	for _, member := range members {
		canvas.Text(5, 100+i*20, member.Name, "fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 12px;")
		for j := 1; j < 26; j++ {
			day, ok := member.CompletionDayLevel[strconv.FormatInt(int64(j), 10)]
			if !ok {
				PrintStar(canvas, i, j, "#3f3f46")
			} else if _, ok := day["1"]; ok {
				if _, ok := day["2"]; ok {
					PrintStar(canvas, i, j, "#F2C83B")
				} else {
					PrintStar(canvas, i, j, "#d4d4d8")
				}
			}
		}
		i++
	}
}

func PrintStar(canvas *svg.SVG, i int, j int, colour string) {
	canvas.Text(nameWidth+starGap*j,
		96+i*20+3,
		"*",
		fmt.Sprintf("fill: %s; font-family: Fira Code; font-size: 12px;"+
			"dominant-baseline:middle; text-anchor:middle;",
			colour),
	)

}
