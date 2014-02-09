package commands

import (
	"fmt"
	irc "github.com/darkliquid/goirc/client"
	"github.com/darkliquid/leader1/config"
	"github.com/darkliquid/leader1/database"
	"github.com/fluffle/golog/logging"
	"strings"
	"time"
)

// Accepts conn so we can use the client to respond,
// Accepts line for details about the raw response line
// Accepts target for being able to respond back to the same place the command was sent to
func LikeTrack(conn *irc.Conn, line *irc.Line, target string) {
	doLikeOrHate(conn, line, target, "like")
}

// Accepts conn so we can use the client to respond,
// Accepts line for details about the raw response line
// Accepts target for being able to respond back to the same place the command was sent to
func HateTrack(conn *irc.Conn, line *irc.Line, target string) {
	doLikeOrHate(conn, line, target, "dislike")
}

func doLikeOrHate(conn *irc.Conn, line *irc.Line, target string, likeType string) {
	cfg := config.Config

	// Get db conn
	db, err := database.DB()
	if err != nil {
		conn.Privmsg(target, fmt.Sprintf("%s: sorry, can't connect to db", line.Nick))
		return
	}

	// Get some stats
	stats, err := getShoutcastStats()
	if err != nil {
		conn.Privmsg(target, fmt.Sprintf("%s: %s", line.Nick, err.Error()))
		return
	}

	if strings.TrimSpace(stats.SongTitle) == "" {
		logging.Error(fmt.Sprintf("Song title missing for some reason when do a %s for %s at %s", likeType, line.Nick, time.Now()))
		return
	}

	rows, err := db.Query("INSERT INTO likelogs (type, user, song, date) VALUES (?, ?, ?, ?)", likeType, line.Nick, stats.SongTitle, time.Now().Unix())
	defer rows.Close()

	switch {
	case err != nil:
		logging.Error(fmt.Sprintf("Failed to add %s to db :( - %s", likeType, err.Error()))
		conn.Privmsg(cfg.Irc.StaffChannel, fmt.Sprintf("%s: db conn failed on %s from %s of %s", strings.ToUpper(likeType), likeType, line.Nick, stats.SongTitle))
		conn.Notice(line.Nick, fmt.Sprintf("Sorry, I couldn't register your %s of %s", likeType, stats.SongTitle))
	default:
		// Send a message to staff
		conn.Privmsg(cfg.Irc.StaffChannel, fmt.Sprintf("%s: %s %sd %s", strings.ToUpper(likeType), line.Nick, likeType, stats.SongTitle))
		conn.Notice(line.Nick, fmt.Sprintf("Thanks for that! You %sd %s", likeType, stats.SongTitle))
	}
}
