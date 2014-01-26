package commands

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
)

// Accepts conn so we can use the client to respond,
// Accepts line for details about the raw response line
// Accepts target for being able to respond back to the same place the command was sent to
func NotAPing(conn *irc.Conn, line *irc.Line, target string) {
	conn.Privmsg(target, fmt.Sprintf("%s: Not a pong", line.Nick))
}
