package commands

import (
	"fmt"
	irc "github.com/darkliquid/goirc/client"
	"strings"
)

func Hand(conn *irc.Conn, line *irc.Line, target string, args []string) {

	switch {
	case len(args) == 1:
		conn.Action(target, fmt.Sprintf("hands %s to %s", args[0], line.Nick))
		return
	case len(args) == 2 && args[0] == "me":
		conn.Action(target, fmt.Sprintf("hands %s to %s", args[1], line.Nick))
		return
	case len(args) == 3 && args[1] == "to":
		conn.Action(target, fmt.Sprintf("hands %s to %s, courtesy of %s", args[0], args[2], line.Nick))
		return
	case len(args) > 3:
		message := strings.Join(args, " ")
		if split := strings.LastIndex(message, " to "); split > -1 {
			conn.Action(target, fmt.Sprintf("hands %s to %s, courtesy of %s", message[:split], message[split+4:], line.Nick))
			return
		}
	}

	conn.Privmsg(target, fmt.Sprintf("%s: you what now? Use !hand [object] to [person] OR !hand me [object]", line.Nick))
}
