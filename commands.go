package main

import (
	"fmt"
	"github.com/darkliquid/leader1/commands"
	irc "github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"strings"
)

// Generic PrivMsg handler
// This is the handler we use to pick up any old message being sent
// We then search these for !commands and dispatch via our switch
// statement
func priv_msg_handler(conn *irc.Conn, line *irc.Line) {
	// PRIVMSG comes in the format of [target msg]
	// where destination could be a user or a #channel

	// Gets the target (user, #channel, etc)
	target := line.Args[0]

	// If it's a !command
	if strings.HasPrefix(line.Args[1], "!") {
		// Split command line on spaces
		args := strings.Split(line.Args[1], " ")

		logging.Debug(fmt.Sprintf("Got ! command \"%s\" with args %#v", args[0], args[1:]))

		// Switch on the firts part of the line (i.e. the actual command)
		switch args[0] {
		case "!not_a_ping":
			commands.NotAPing(conn, line, target)
		}
	}
}

// Register the defined commands with the client
func RegisterCommands() {
	client.AddHandler("PRIVMSG", priv_msg_handler)
}
