package commands

import (
	"fmt"
	irc "github.com/darkliquid/goirc/client"
)

// Accepts conn so we can use the client to respond,
// Accepts line for details about the raw response line
// Accepts target for being able to respond back to the same place the command was sent to
func G3Song(conn *irc.Conn, line *irc.Line, target string) {
	stats, err := getShoutcastStats()
	if err != nil {
		conn.Privmsg(target, fmt.Sprintf("%s: %s", line.Nick, err.Error()))
		return
	}

	conn.Privmsg(target, fmt.Sprintf("%s: current song is - `%s`", line.Nick, stats.SongTitle))
}