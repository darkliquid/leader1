package plugins

import (
	"github.com/darkliquid/leader1/config"
	"github.com/darkliquid/go-ircevent"
	"strings"
	"github.com/robertkrimen/otto"
	"log"
)

type PluginFunc struct {
	function func(*irc.Event)
	help     string
}

type Plugin struct {
	commands     map[string]*PluginFunc
	callbacks    map[string][]*PluginFunc
	log          *log.Logger
	js			 *otto.Otto
	cfg          *config.Settings
}

func (p *Plugin) SetCommand(name string, command otto.Value, help string) {
	if _, ok := p.commands[name]; ok {
		p.log.Printf("Warning: Command `%s` was already defined. Overriding...", name)
	}
	wrappedCommand := func(event *irc.Event) {
		_, err := command.Call(p.js.ToValue(event))
		if err != nil {
			p.log.Printf("Command `%s` errored: %s", name, err)
		}
	}
	p.commands[name] = &PluginFunc{
		function: wrappedCommand,
		help:     help,
	}
}

func (p *Plugin) AddCallback(eventCode string, name string, callback otto.Value) {
	wrappedCallback := func(event *irc.Event) {
		_, err := callback.Call(p.js.ToValue(event))
		if err != nil {
			p.log.Printf("Callback `%s` (%#v) for event code `%s` errored: %s", name, callback, eventCode, err)
		}
	}
	p.callbacks[eventCode] = append(p.callbacks[eventCode], &PluginFunc{
		function: wrappedCallback,
		help: name,
	})
}

func (p *Plugin) RunCallbacks(event *irc.Event) {
	if callbacks, ok := p.callbacks[event.Code]; ok {
		if p.cfg.Irc.Debug || p.cfg.Debug {
			p.log.Printf("%v (%v) >> %#v\n", event.Code, len(callbacks), event)
		}

		for _, callback := range callbacks {
			go callback.function(event)
		}
	}
}

func (p *Plugin) RunCommand(event *irc.Event) bool {
	if event.Message[0] == '!' && len(event.Message) > 1 {
		call := strings.SplitN(event.Message[1:], " ", 2)
		command := call[0]

		var ok bool
		if _, ok = p.commands[command]; ok {
			if p.cfg.Irc.Debug || p.cfg.Debug {
				p.log.Printf("%v (!%v) >> %#v\n", event.Code, command, event)
			}
			cmd := p.commands[command].function
			go cmd(event)
		}

		return ok
	}
	return false
}