# AOCBot
Fork of DylanMeeus Advent of Code Bot for our private discord servers. It will display the names and stars of every user on your leaderboard.

## Running the Bot
You can use the provided GHCR.io container or build it yourself.

You will have to set the following environment variables:

| Variable      | Description                                                                                                                                                                                           |
|---------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| LEADERBOARDID | AoC Leaderboard ID, found in the link to the leaderboard, e.g. https://adventofcode.com/2020/leaderboard/private/view/123456                                                                          |
| BOT_TOKEN     | Your Discord Bot token, see https://discord.com/developers/docs/intro                                                                                                                                 |
| SESSION_TOKEN | The Advent of Code website session token. You can get this token by logging into an account that has access to your leaderboard and extracting the session cookie from your browsers developer tools. |