package commands

import (
	"fmt"
	irc "github.com/darkliquid/goirc/client"
	"net/url"
)

// Accepts conn so we can use the client to respond,
// Accepts line for details about the raw response line
// Accepts target for being able to respond back to the same place the command was sent to
// Accepts query as the search term
func LMGTFY(conn *irc.Conn, line *irc.Line, target string, query string) {
	url := "http://lmgtfy.com/?q=" + url.QueryEscape(query)
	conn.Privmsg(target, fmt.Sprintf("%s: Let me google that for you - %s", line.Nick, url))
}
