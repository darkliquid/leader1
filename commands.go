package main

import (
	"fmt"
	irc "github.com/darkliquid/goirc/client"
	"github.com/darkliquid/leader1/commands"
	"github.com/fluffle/golog/logging"
	"strings"
)

// Generic PrivMsg handler
// This is the handler we use to pick up any old message being sent
// We then search these for !commands and dispatch via our switch
// statement
func privMsgHandler(conn *irc.Conn, line *irc.Line) {
	// PRIVMSG comes in the format of [target msg]
	// where destination could be a user or a #channel

	// Gets the target (user, #channel, etc)
	target := line.Args[0]

	// If it's a !command
	if strings.HasPrefix(line.Args[1], "!") {
		// Split command line into arguments
		args := strings.Fields(line.Args[1])

		if len(args) > 1 {
			logging.Debug(fmt.Sprintf("Got ! command \"%s\" with args %#v", args[0], args[1:]))
		} else {
			logging.Debug(fmt.Sprintf("Got ! command \"%s\" with no arguments", args[0]))
		}

		// Switch on the first part of the line (i.e. the actual command)
		switch {
		case args[0] == "!ping":
			commands.Ping(conn, line, target)
		case args[0] == "!lmgtfy" && len(args) > 1:
			commands.LMGTFY(conn, line, target, strings.Join(args[1:], " "))
		case args[0] == "!urban" && len(args) > 1:
			commands.UrbanDictionary(conn, line, target, strings.Join(args[1:], " "))
		case args[0] == "!time":
			commands.Time(conn, line, target)
		case args[0] == "!g3song":
			commands.G3Song(conn, line, target)
		case args[0] == "!listeners":
			commands.Listeners(conn, line, target)
		case args[0] == "!+":
			commands.LikeTrack(conn, line, target)
		case args[0] == "!-":
			commands.HateTrack(conn, line, target)
		case args[0] == "!request" && len(args) > 1:
			commands.Request(conn, line, target, strings.Join(args[1:], " "))
		case args[0] == "!announce" && len(args) > 1:
			commands.Announce(conn, line, target, strings.Join(args[1:], " "))
		case args[0] == "!invite" && len(args) > 1:
			commands.Invite(conn, line, target, args[1:])
		case args[0] == "!autovoice" && len(args) > 1:
			commands.AutoVoice(conn, line, target, args[1])
		case args[0] == "!hand":
			if len(args) > 1 {
				commands.Hand(conn, line, target, args[1:])
			} else {
				conn.Privmsg(target, fmt.Sprintf("%s: you what now? Use !hand [object] to [person] OR !hand me [object]", line.Nick))
			}
		}
		return
	}

	// Auto title URLs that are found in non-command lines
	if url, err := commands.ExtractURL(line.Args[1]); err == nil {
		logging.Debug(fmt.Sprintf("Extracted URL(%s) from %s", url, line.Args[1]))
		commands.URLTitler(conn, line, target, url)
	} else {
		logging.Debug(fmt.Sprintf("Couldn't extract URL from %s - %s", line.Args[1], err.Error()))
	}
}

// Register the defined commands with the client
func RegisterCommands() {
	client.AddHandler("PRIVMSG", privMsgHandler)
}
