package commands

import (
	"fmt"
	irc "github.com/darkliquid/goirc/client"
)

// Accepts conn so we can use the client to respond,
// Accepts line for details about the raw response line
// Accepts target for being able to respond back to the same place the command was sent to
func Listeners(conn *irc.Conn, line *irc.Line, target string) {
	stats, err := getShoutcastStats()
	if err != nil {
		conn.Privmsg(target, fmt.Sprintf("%s: %s", line.Nick, err.Error()))
		return
	}

	switch {
	case stats.UniqueListeners == 0:
		conn.Privmsg(target, fmt.Sprintf("%s: there are currently no listeners on the stream :'(", line.Nick))
	case stats.UniqueListeners == 1:
		conn.Privmsg(target, fmt.Sprintf("%s: there is currently 1 listener on the stream", line.Nick))
	case stats.UniqueListeners == stats.MaxListeners:
		conn.Privmsg(target, fmt.Sprintf("%s: holy crap! We have maxed out at %d listeners on the stream!!!", line.Nick, stats.UniqueListeners))
	default:
		conn.Privmsg(target, fmt.Sprintf("%s: there are currently %d listeners on the stream", line.Nick, stats.UniqueListeners))
	}
}