package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type IrcSettings struct {
	Host          string
	Port          string
	Nick          string
	NickPass      string
	Pass          string
	Ssl           bool
	NormalChannel string `json:"normal_channel"`
	StaffChannel  string `json:"staff_channel"`
	MaxFailures   int
	Timeout       uint
	KeepAlive     uint `json:"keep_alive"`
	PingFreq      uint `json:"ping_frequency"`
	AutoVoice     bool
	Version       string
	Debug         bool
	PluginsDir    string `json:"plugin_dir"`
}

type DbSettings struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
}

type StreamSettings struct {
	StatsURL  string `json:"stats_url"`
	StatsUser string `json:"stats_user"`
	StatsPass string `json:"stats_pass"`
}

type Settings struct {
	Irc    IrcSettings
	Db     DbSettings
	Stream StreamSettings
	Debug  bool
}

func Load() *Settings {
	log := log.New(os.Stdout, "[config] ", log.LstdFlags)
	// Gets the current executable path for use a cfgfile path
	cwd, err := os.Getwd()

	// Bail out if the executable can not be found
	if err != nil {
		// Init here since we haven't done so yet
		log.Fatal("Can not find current working directory!")
	}

	// Default config path
	config_path := filepath.Join(cwd, "leader-1.json")

	// Allow setting of the cfg file to load stuff from
	flag.StringVar(&config_path, "cfg", config_path, "sets the config file")

	// Parse command line settings
	if !flag.Parsed() {
		flag.Parse()
	}

	log.Println("Loading...")

	// Load config file
	file, err := ioutil.ReadFile(config_path)

	// Bail out if reading the config file errored
	if err != nil {
		log.Fatal("Couldn't read config file")
	}

	log.Println("Reading config file")

	// Unmarshal json config file into our config structure
	var cfg Settings
	err = json.Unmarshal(file, &cfg)

	// Bail out if the demarshalling fails
	if err != nil {
		log.Fatalln("Couldn't parse config!")
	}

	// Set default Max/Idle DB Conns
	if cfg.Db.MaxOpenConns == 0 {
		cfg.Db.MaxOpenConns = 5
	}
	if cfg.Db.MaxIdleConns == 0 {
		cfg.Db.MaxIdleConns = 5
	}

	if cfg.Irc.PluginsDir == "" {
		cfg.Irc.PluginsDir = filepath.Join(cwd, "plugins")
	}

	log.Println("Loaded config")

	return &cfg
}
