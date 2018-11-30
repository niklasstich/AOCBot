Advent Of Code bot for discord!

This is a bit for discord that can print your private leaderboard to discord. 

## Installation

1) Have Go working. 
2) install discordgo: `go get github.com/bwmarrin/discordgo`
3) Alter the file `resources/data.json`:
    * channel: your private channel ID
    * bot_token: when you create a discord bot, it will give you a token. (More information
      here)[https://discordapp.com/developers/docs/topics/oauth2]
    * session_token: When inspecting the json data from your leaderboard (which can be found on
      adventofcode.com) you can look at the headers used to give your json data. Copy the sessionId
      to this field. (It's the only way of authenticating I could find)


## Usage

In discord, there are several commands available to get the data from this year and previous years.

For this year use:

* /aoc
* /aoc2018

for previous years, replace 2018 with the correct year such as:

* /aoc2016
* /aoc2017

If you want to limit the output (by default it will return everyone) you can append a number of max
output, such as:

* /aoc 10 // returns the top 10
* /aoc2015 100 // return the top 100


