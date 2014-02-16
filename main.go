package main

import (
	"fmt"
	"github.com/darkliquid/leader1/bot"
	"github.com/darkliquid/leader1/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load config
	cfg := config.Load()

	// Set up Irc Client
	client, err := bot.New(cfg)

	err = client.Connect()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	trap := make(chan os.Signal, 1)
	quit := make(chan bool)

	signal.Notify(trap, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-trap
		close(trap)
		fmt.Printf("Caught: %s - quitting\n", sig)
		client.Quit()
		quit <- true
	}()

	// Block here for quit channel
	<-quit
	close(quit)

	fmt.Println("Quitting...")
	os.Exit(0)
}
