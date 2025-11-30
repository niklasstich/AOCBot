package svg

import (
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"

	svg "github.com/ajstarks/svgo"
	"github.com/niklasstich/AOCBot/aoc"
)

var (
	nameWidth = 140
	starGap   = 22
	canvas    *svg.SVG
)

func GenerateSvg(year int, members []aoc.Member, nextPuzzle *aoc.NextPuzzle, filepath string) (err error) {
	fmt.Println("Generating new image...")
	numDays := aoc.DaysForYear(year)

	// Adjust width based on number of days
	width := nameWidth + starGap*(numDays+1)
	baseHeight := 60 + 40 + len(members)*20

	height := baseHeight
	if nextPuzzle != nil {
		height += 30
	}
	file, err := os.Create(filepath)

	if err != nil {
		return
	}
	defer file.Close()
	canvas = svg.New(file)
	canvas.Start(width, height)

	// Calculate where the vertical line should end (before countdown area)
	verticalLineEnd := baseHeight

	//header and lines
	canvas.Rect(0, 0, width, height, "fill: rgb(15,15,35);")
	canvas.Text(width/2, 30,
		"Advent of Code Leaderboard",
		"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
			"dominant-baseline:middle; text-anchor:middle;"+
			"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00;")
	AddRandomYearMarkup(year, width)
	canvas.Line(0, 80, width, 80, "stroke: rgb(0, 99, 0); stroke-width: 2;")
	canvas.Line(nameWidth, 80, nameWidth, verticalLineEnd, "stroke: rgb(0, 99, 0); stroke-width: 2;")

	//numbers for stars
	for i := 1; i <= numDays; i++ {
		canvas.Text(nameWidth+starGap*i, 70, fmt.Sprint(i), "fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 12px;"+
			"dominant-baseline:middle; text-anchor:middle;")
	}

	PrintMembers(canvas, members, numDays)

	// Add countdown text if there's a next puzzle
	if nextPuzzle != nil {
		PrintNextPuzzleCountdown(canvas, nextPuzzle, baseHeight+20, width)
	}

	canvas.End()
	return
}

// PrintNextPuzzleCountdown renders the countdown text at the bottom of the image
func PrintNextPuzzleCountdown(canvas *svg.SVG, nextPuzzle *aoc.NextPuzzle, y int, width int) {
	duration := time.Until(nextPuzzle.UnlockAt)

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	var countdownText string
	if hours > 0 {
		countdownText = fmt.Sprintf("Day %d puzzle unlocks in %dh %dm", nextPuzzle.Day, hours, minutes)
	} else if minutes > 0 {
		countdownText = fmt.Sprintf("Day %d puzzle unlocks in %dm", nextPuzzle.Day, minutes)
	} else {
		countdownText = fmt.Sprintf("Day %d puzzle unlocks soon!", nextPuzzle.Day)
	}

	canvas.Text(width/2, y, countdownText,
		"fill: rgb(204, 204, 0); font-family: Fira Code; font-size: 14px;"+
			"dominant-baseline:middle; text-anchor:middle;"+
			"text-shadow: 0 0 2px #cccc00, 0 0 5px #cccc00;")
}

func ConvertSvgToPng(svg string, png string) (err error) {
	command := exec.Command("inkscape",
		"--export-type=png",
		fmt.Sprintf("--export-filename=%s", png),
		"--export-dpi=200",
		svg)
	err = command.Run()
	return
}

func PrintMembers(canvas *svg.SVG, members []aoc.Member, numDays int) {
	i := 0
	for _, member := range members {
		canvas.Text(5, 100+i*20, member.Name, "fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 12px;")
		for j := 1; j <= numDays; j++ {
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

/*
Add random year markup to the image
Just like on the real leaderboard
*/
func AddRandomYearMarkup(year int, width int) {
	//If we don't add random seed after a certain amount of request
	//the random numbers will be the same
	rand.NewSource(time.Now().UnixNano())
	min := 0
	max := 5
	rand := rand.Intn(max-min+1) + min

	// Center position for the year display
	centerX := width / 2

	switch rand {
	case 1:
		canvas.Text(centerX-32, 50,
			"$year=",
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:end;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00; opacity: 0.4; letter-spacing: -0.05em;")
		canvas.Text(centerX-32, 50,
			fmt.Sprintf("%d;", year),
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:start;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00;")
	case 2:
		canvas.Text(centerX-32, 50,
			"0xffff&",
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:end;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00; opacity: 0.4; letter-spacing: -0.05em;")
		canvas.Text(centerX-32, 50,
			fmt.Sprintf("%d", year),
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:start;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00;")
	case 3:
		canvas.Text(centerX-32, 50,
			"/*",
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:end;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00; opacity: 0.4; letter-spacing: -0.05em;")
		canvas.Text(centerX-32, 50,
			fmt.Sprintf("%d", year),
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:start;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00;")
		canvas.Text(centerX+20, 50,
			"*/",
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:start;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00; opacity: 0.4; letter-spacing: -0.05em;")
	case 4:
		canvas.Text(centerX-32, 50,
			"λy.",
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:end;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00; opacity: 0.4; letter-spacing: -0.05em;")
		canvas.Text(centerX-32, 50,
			fmt.Sprintf("%d", year),
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:start;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00;")
	case 5:
		canvas.Text(centerX-32, 50,
			"/^",
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:end;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00; opacity: 0.4; letter-spacing: -0.05em;")
		canvas.Text(centerX-32, 50,
			fmt.Sprintf("%d", year),
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:start;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00;")
		canvas.Text(centerX+20, 50,
			"$/",
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:start;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00; opacity: 0.4; letter-spacing: -0.05em;")
	case 6:
		canvas.Text(centerX-32, 50,
			"//",
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:end;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00; opacity: 0.4; letter-spacing: -0.05em;")
		canvas.Text(centerX-32, 50,
			fmt.Sprintf("%d", year),
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:start;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00;")
	default:
		canvas.Text(centerX-32, 50,
			"var y=",
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:end;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00; opacity: 0.4; letter-spacing: -0.05em;")
		canvas.Text(centerX-32, 50,
			fmt.Sprintf("%d;", year),
			"fill: rgb(0, 144, 0); font-family: Fira Code; font-size: 20.7px;"+
				"dominant-baseline:middle; text-anchor:start;"+
				"text-shadow: 0 0 1px #00cc00, 0 0 7px #00cc00;")
	}
}
