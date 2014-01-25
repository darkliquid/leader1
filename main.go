package main;

import (
    "flag"
    "os"
    "path/filepath"
)

func main() {
    // Gets the current executable path for use a cfgfile path
    cwd, err := os.Getwd()

    // Bail out if the executable can not be found
    if err != nil {
        panic("Can not find current working directory!")
    }

    // Default config path
    config_path := filepath.Join(cwd, "leader-1.json")

    // Allow setting of the cfg file to load stuff from
    flag.StringVar(&config_path, "cfg", config_path, "sets the config file")

    flag.Parse()
}