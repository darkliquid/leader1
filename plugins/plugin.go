package plugins

import (
	"github.com/darkliquid/go-ircevent"
	"github.com/darkliquid/leader1/config"
	"github.com/robertkrimen/otto"
	"log"
	"strings"
)

type pluginFunc struct {
	function func(*pluginEnvironment)
	help     string
}

type pluginEnvironment struct {
	Event *irc.Event
	Log   func(call otto.FunctionCall) otto.Value
}

type Plugin struct {
	commands  map[string]*pluginFunc
	callbacks map[string][]*pluginFunc
	log       *log.Logger
	js        *otto.Otto
	cfg       *config.Settings
}

func (p *Plugin) SetCommand(name string, command otto.Value, help string) {
	if _, ok := p.commands[name]; ok {
		p.log.Printf("Warning: Command `%s` was already defined. Overriding...", name)
	}
	wrappedCommand := func(env *pluginEnvironment) {
		_, err := command.Call(p.js.ToValue(env))
		if err != nil {
			p.log.Printf("Command `%s` errored: %s", name, err)
		}
	}
	p.commands[name] = &pluginFunc{
		function: wrappedCommand,
		help:     help,
	}
}

func (p *Plugin) AddCallback(eventCode string, name string, callback otto.Value) {
	wrappedCallback := func(env *pluginEnvironment) {
		_, err := callback.Call(p.js.ToValue(env))
		if err != nil {
			p.log.Printf("Callback `%s` (%#v) for event code `%s` errored: %s", name, callback, eventCode, err)
		}
	}
	p.callbacks[eventCode] = append(p.callbacks[eventCode], &pluginFunc{
		function: wrappedCallback,
		help:     name,
	})
}

func (p *Plugin) RunCallbacks(event *irc.Event) {
	if callbacks, ok := p.callbacks[event.Code]; ok {
		if p.cfg.Irc.Debug || p.cfg.Debug {
			p.log.Printf("%v (%v) >> %#v\n", event.Code, len(callbacks), event)
		}

		for _, callback := range callbacks {
			callback.function(p.jsEnv(event))
		}
	}
}

func (p *Plugin) RunCommand(event *irc.Event) bool {
	if event.Message()[0] == '!' && len(event.Message()) > 1 {
		call := strings.SplitN(event.Message()[1:], " ", 2)
		command := call[0]

		var ok bool
		if _, ok = p.commands[command]; ok {
			if p.cfg.Irc.Debug || p.cfg.Debug {
				p.log.Printf("%v (!%v) >> %#v\n", event.Code, command, event)
			}

			cmd := p.commands[command].function
			cmd(p.jsEnv(event))
		}

		return ok
	}
	return false
}

func (p *Plugin) jsEnv(event *irc.Event) *pluginEnvironment {
	return &pluginEnvironment{
		Event: event,
		Log: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
				p.log.Println(call.ArgumentList[0].String())
				return otto.TrueValue()
			} else {
				return otto.FalseValue()
			}
		},
	}
}