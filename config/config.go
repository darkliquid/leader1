package config

import (
	"encoding/json"
	"flag"
	"github.com/fluffle/golog/logging"
	"io/ioutil"
	"os"
	"path/filepath"
)

type IrcSettings struct {
	Host     string
	Port     string
	Nick     string
	Pass     string
	Ssl      bool
	Channels []string
	MaxFailures int
}

type DbSettings struct {
	User string
	Pass string
	Host string
	Port string
}

type Settings struct {
	Irc IrcSettings
	Db  DbSettings
}

var Config Settings

func Load() {
	// Gets the current executable path for use a cfgfile path
	cwd, err := os.Getwd()

	// Bail out if the executable can not be found
	if err != nil {
		// Init here since we haven't done so yet
		logging.InitFromFlags()
		logging.Fatal("Can not find current working directory!")
	}

	// Default config path
	config_path := filepath.Join(cwd, "leader-1.json")

	// Allow setting of the cfg file to load stuff from
	flag.StringVar(&config_path, "cfg", config_path, "sets the config file")

	// Parse command line settings
	if !flag.Parsed() {
		flag.Parse()
	}

	// Need to init logging engine
	logging.InitFromFlags()

	logging.Info("Loading...")

	// Load config file
	file, err := ioutil.ReadFile(config_path)

	// Bail out if reading the config file errored
	if err != nil {
		logging.Fatal("Couldn't read config file")
	}

	logging.Info("Reading config file")

	// Unmarshal json config file into our config structure
	err = json.Unmarshal(file, &Config)

	// Bail out if the demarshalling fails
	if err != nil {
		logging.Fatal("Couldn't parse config!")
	}

	logging.Info("Loaded config")

	return
}
