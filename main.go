package main;

import (
    "github.com/VividCortex/godaemon"
)

func main() {
    godaemon.MakeDaemon(&godaemon.DaemonAttr{})
}