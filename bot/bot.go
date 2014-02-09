package bot

import (
	"github.com/darkliquid/leader1/config"
	"github.com/darkliquid/go-ircevent"
	"time"
	"net"
	"github.com/robertkrimen/otto"
	"log"
	"os"
)

type Bot struct {
	conn *irc.Connection
	cfg config.Settings
	js *otto.Otto
}

func (bot *Bot) Quit() {
	bot.conn.Quit()
}

func (bot *Bot) Connect() error {
	bot.InitCallbacks()
	// Connect
	if bot.cfg.Irc.Port != "" {
		return bot.conn.Connect(net.JoinHostPort(bot.cfg.Irc.Host, bot.cfg.Irc.Port))
	}
	return bot.conn.Connect(bot.cfg.Irc.Host)
}

func (bot *Bot) InitCallbacks() error {
	// Setup callbacks
	bot.conn.AddCallback("001", func(event *irc.Event) {
		bot.conn.Join("#dl-dev-test")
	});
	return nil
}

func New(cfg config.Settings) (*Bot, error) {
	// Set up Irc Client
	client := irc.IRC(cfg.Irc.Nick, "leader-1")

	if cfg.Irc.Version != "" {
		client.Version = cfg.Irc.Version		
	}
	
	if cfg.Irc.Debug || cfg.Debug {
		client.Debug = true		
	}
	if cfg.Irc.Timeout > 0 {
		// Set client timeout to configured amount in seconds
		client.Timeout = time.Duration(cfg.Irc.Timeout) * time.Second
	}
	if cfg.Irc.KeepAlive > 0 {
		// Set client keepalive duration to configured amount in seconds
		client.KeepAlive = time.Duration(cfg.Irc.KeepAlive) * time.Second
	}
	if cfg.Irc.PingFreq > 0 {
		// Set client pingfreq duration to configured amount in seconds
		client.PingFreq = time.Duration(cfg.Irc.PingFreq) * time.Second
	}

	// Optionally, enable SSL
	client.UseTLS = cfg.Irc.Ssl

	// Setup IRC logger
	client.Log = log.New(os.Stdout, "[irc] ", log.LstdFlags)

	return &Bot{
		cfg: cfg,
		conn: client,
		js: otto.New(),
	}, nil
}