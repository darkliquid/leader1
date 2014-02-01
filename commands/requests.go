package commands

import (
	"fmt"
	irc "github.com/darkliquid/goirc/client"
	"github.com/darkliquid/leader1/config"
	"github.com/darkliquid/leader1/database"
	"github.com/fluffle/golog/logging"
	"time"
)

// Accepts conn so we can use the client to respond,
// Accepts line for details about the raw response line
// Accepts target for being able to respond back to the same place the command was sent to
func Request(conn *irc.Conn, line *irc.Line, target string, request string) {
	cfg := config.Config

	// Get db conn
	db, err := database.DB()
	if err != nil {
		conn.Privmsg(target, fmt.Sprintf("%s: sorry, can't connect to db", line.Nick))
		return
	}

	rows, err := db.Query("INSERT INTO requests (user, song, date) VALUES (?, ?, ?)", line.Nick, request, time.Now().Unix())
	defer rows.Close()

	switch {
	case err != nil:
		logging.Error(fmt.Sprintf("Failed to add like to db :( - %s", err.Error()))
		conn.Privmsg(cfg.Irc.StaffChannel, fmt.Sprintf("REQ: db conn failed on request from %s of %s", line.Nick, request))
		conn.Notice(line.Nick, fmt.Sprintf("Sorry, I couldn't register your request of %s", request))
	default:
		// Send a message to staff
		conn.Privmsg(cfg.Irc.StaffChannel, fmt.Sprintf("REQ: %s requested %s", line.Nick, request))
		conn.Notice(line.Nick, fmt.Sprintf("Thanks for that! You requested %s", request))
	}
}
