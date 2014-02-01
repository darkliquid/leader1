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
func Announce(conn *irc.Conn, line *irc.Line, target string, message string) {
	cfg := config.Config

	if target != cfg.Irc.StaffChannel {
		logging.Warn(fmt.Sprintf("Attempt to use !announce by %s out of %s", line.Nick, cfg.Irc.StaffChannel))
		return
	}
	conn.Notice(cfg.Irc.NormalChannel, message)
}