package main

import (
	"fmt"
	"github.com/darkliquid/leader1/config"
	irc "github.com/darkliquid/goirc/client"
	"github.com/fluffle/golog/logging"
	"os"
	"os/signal"
	"time"
)

var client *irc.Conn

func main() {
	// Load our configuration
	config.Load()

	// Alias it for easy use
	cfg := config.Config

	// Set up Irc Client
	client = irc.SimpleClient(cfg.Irc.Nick, "leader-1", "A mighty, mighty Go Bot")

	// Set client timeout to 1 second
	client.Timeout = time.Second

	// Track the state of various things
	client.EnableStateTracking()

	// Optionally, enable SSL
	client.SSL = cfg.Irc.Ssl

	// Join channels on connect
	client.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		for _, channel := range cfg.Irc.Channels {
			conn.Join(channel)
			logging.Info(fmt.Sprintf("Joining channel %s", channel))
		}
	})

	// And a signal on disconnect
	quit := make(chan bool)
	client.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { quit <- true })

	// Register commands
	RegisterCommands()

	// Trap interrupt signal so we can cleanly disconnect on fail
	trap := make(chan os.Signal, 1)
	signal.Notify(trap, os.Interrupt)

	// Var we check to see if we want to actually quit the app
	really_quit := false

	// Setup connection failure count
	connection_failures := 0

	// concurrent handler for trapping SIGINT
	go func() {
		for sig := range trap {
		    really_quit = true
			client.Quit(fmt.Sprintf("Goodbye (%s)", sig))
		}
	}()

	for !really_quit {
		// connect to server
		logging.Info(fmt.Sprintf("Connection to %s:%s", cfg.Irc.Host, cfg.Irc.Port))

		// This connection blocks for an unknown number of seconds based on the system settings
		if err := client.Connect(cfg.Irc.Host + ":" + cfg.Irc.Port); err != nil {
			// At the moment, just fail, but ideally we will retry until a maximum number of failures is exceeded
			connection_failures++
			if connection_failures > cfg.Irc.MaxFailures {
				really_quit = true
			}
			quit <- true
		} else {
			logging.Info(fmt.Sprintf("Connected to irc server as %s", cfg.Irc.Nick))
			connection_failures = 0
		}

		// wait on quit channel
		<-quit
	}
}
