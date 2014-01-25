package main

import (
	"fmt"
	"github.com/darkliquid/leader1/config"
	irc "github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"os"
	"os/signal"
)

func main() {
	// Load our configuration
	config.Load()

	// Alias it for easy use
	cfg := config.Config

	// Set up Irc Client
	c := irc.SimpleClient(cfg.Irc.Nick)

	// Optionally, enable SSL
	c.SSL = cfg.Irc.Ssl

	// Join channels on connect
	c.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		for _, channel := range cfg.Irc.Channels {
			conn.Join(channel)
			logging.Info(fmt.Sprintf("Joining channel %s", channel))
		}
	})

	// And a signal on disconnect
	quit := make(chan bool)
	c.AddHandler(irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) { quit <- true })

	logging.Info(fmt.Sprintf("Connection to %s:%s", cfg.Irc.Host, cfg.Irc.Port))

	// Tell client to connect
	if err := c.Connect(cfg.Irc.Host + ":" + cfg.Irc.Port); err != nil {
		// At the moment, just fail, but ideally we will retry until a maximum number of failures is exceeded
		logging.Fatal(fmt.Sprintf("Connection error: %s\n", err.Error()))
	}
	logging.Info(fmt.Sprintf("Connected to irc server as %s", cfg.Irc.Nick))

    // Trap interrupt signal so we can cleanly disconnect on fail
	trap := make(chan os.Signal, 1)
    signal.Notify(trap, os.Interrupt)

    go func(){
        for sig := range trap {
            c.Quit(fmt.Sprintf("Goodbye (%s)", sig))
            quit <- true
        }
    }()

	// Wait for disconnect
	<-quit
}
