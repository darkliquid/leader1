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
		// Split command line on spaces
		raw_args := strings.Split(line.Args[1], " ")
		var args []string

		// Strip out empty strings
		for _, arg := range raw_args {
			if arg != "" {
				args = append(args, arg)
			}
		}

		logging.Debug(fmt.Sprintf("Got ! command \"%s\" with args %#v", args[0], args[1:]))

		// Switch on the firts part of the line (i.e. the actual command)
		switch args[0] {
		case "!ping":
			commands.Ping(conn, line, target)
		case "!lmgtfy":
			commands.LMGTFY(conn, line, target, strings.Join(args[1:], " "))
		case "!urban":
			commands.UrbanDictionary(conn, line, target, strings.Join(args[1:], " "))
		case "!time":
			commands.Time(conn, line, target)
		case "!g3song":
			commands.G3Song(conn, line, target)
		case "!listeners":
			commands.Listeners(conn, line, target)
		case "!+":
			commands.LikeTrack(conn, line, target)
		case "!-":
			commands.HateTrack(conn, line, target)
		case "!request":
			commands.Request(conn, line, target, strings.Join(args[1:], " "))
		case "!announce":
			commands.Announce(conn, line, target, strings.Join(args[1:], " "))
		case "!invite":
			commands.Invite(conn, line, target, args[1:])
		}
	}
}

// Register the defined commands with the client
func RegisterCommands() {
	client.AddHandler("PRIVMSG", privMsgHandler)
}
