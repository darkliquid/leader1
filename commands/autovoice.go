package commands

import (
	"fmt"
	irc "github.com/darkliquid/goirc/client"
	"github.com/darkliquid/leader1/config"
	"github.com/fluffle/golog/logging"
)

// Accepts conn so we can use the client to respond,
// Accepts line for details about the raw response line
// Accepts target for being able to respond back to the same place the command was sent to
// Accepts message as the message to send
func AutoVoice(conn *irc.Conn, line *irc.Line, target string, state string) {
	cfg := config.Config

	if target != cfg.Irc.StaffChannel {
		logging.Warn(fmt.Sprintf("Attempt to use !autovoice by %s from %s", line.Nick, target))
		return
	}

	// Op only
	privs, ok := conn.ST.GetNick(line.Nick).PrivsOnStr(target)
	if !ok || (!privs.Owner && !privs.Admin && !privs.Op && !privs.HalfOp) {
		conn.Privmsg(target, fmt.Sprintf("%s: oi! ops only! !autovoice is not for you", line.Nick))
		logging.Warn(fmt.Sprintf("Attempt to use !autovoice by non-op %s", line.Nick))
		return
	}

	switch state {
	case "on":
		config.Config.Irc.AutoVoice = true
		conn.Privmsg(target, fmt.Sprintf("%s: autovoice enabled", line.Nick))
		logging.Info(fmt.Sprintf("AutoVoice enabled by %s", line.Nick))
	case "off":
		config.Config.Irc.AutoVoice = false
		conn.Privmsg(target, fmt.Sprintf("%s: autovoice disabled", line.Nick))
		logging.Info(fmt.Sprintf("AutoVoice disabled by %s", line.Nick))
	default:
		conn.Privmsg(target, fmt.Sprintf("%s: unknown argument, call !autovoice [on|off]", line.Nick))
	}

}
