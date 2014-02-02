package main

import (
	"fmt"
	irc "github.com/darkliquid/goirc/client"
	"github.com/darkliquid/leader1/config"
	"github.com/fluffle/golog/logging"
	"time"
)

// Handler for reclaiming a stolen nick
func ReclaimNick(conn *irc.Conn, line *irc.Line) {
	cfg := config.Config

	if thief := client.ST.GetNick(cfg.Irc.Nick); thief != nil {
		// Recover nick from thieves
		logging.Info("Nick taken - issuing RECOVER/RELEASE...")
		client.Privmsg("NickServ", fmt.Sprintf("RECOVER %s %s", cfg.Irc.Nick, cfg.Irc.NickPass))
		client.Privmsg("NickServ", fmt.Sprintf("RELEASE %s %s", cfg.Irc.Nick, cfg.Irc.NickPass))
	}
	SetBotState()
}

// Automatically give voice to users in channels in which the bot has Op
func AutoVoice(conn *irc.Conn, line *irc.Line) {
	if config.Config.Irc.AutoVoice {
		target := line.Args[0]
		time.Sleep(time.Millisecond * 500) // Wait half a second before we bother

		privs, ok := conn.ST.GetNick(line.Nick).PrivsOnStr(target)

		if ok && (privs.Owner || privs.Admin || privs.Op || privs.HalfOp || privs.Voice) {
			logging.Debug(fmt.Sprintf("Skipping Auto Voicing %s on channel %s", line.Nick, target))
			return
		}

		logging.Debug(fmt.Sprintf("Auto Voicing %s on channel %s", line.Nick, target))
		conn.Mode(target, fmt.Sprintf("+v %s", line.Nick))
	}
}

// Function for setting up the botstate
func SetBotState() {
	cfg := config.Config

	if ghost := client.ST.GetNick(cfg.Irc.Nick); ghost != nil && ghost != client.Me {
		// GHOST the old nick
		logging.Info("Nick taken - issuing GHOST...")
		client.Privmsg("NickServ", fmt.Sprintf("GHOST %s %s", cfg.Irc.Nick, cfg.Irc.NickPass))
	}

	// Set up the nick
	client.Nick(cfg.Irc.Nick)

	// Identify as the nick owner
	client.Privmsg("NickServ", fmt.Sprintf("IDENTIFY %s", cfg.Irc.NickPass))

	// Tell IRC I'm a bot
	client.Mode(cfg.Irc.Nick, "+B")
}
