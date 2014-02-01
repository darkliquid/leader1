package main

import (
	"fmt"
	irc "github.com/darkliquid/goirc/client"
	"github.com/darkliquid/leader1/config"
	"github.com/fluffle/golog/logging"
)

func JoinChannels(conn *irc.Conn, line *irc.Line) {
	cfg := config.Config

	if _, ok := conn.Me.IsOnStr(cfg.Irc.NormalChannel); !ok {
		conn.Join(cfg.Irc.NormalChannel)
		logging.Info(fmt.Sprintf("Joining channel %s", cfg.Irc.NormalChannel))
	}

	if _, ok := conn.Me.IsOnStr(cfg.Irc.StaffChannel); !ok {
		conn.Join(cfg.Irc.StaffChannel)
		logging.Info(fmt.Sprintf("Joining channel %s", cfg.Irc.StaffChannel))
	}
}
