package plugins

import (
	"github.com/darkliquid/go-ircevent"
	"github.com/darkliquid/leader1/config"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"runtime"
)

type PluginManager struct {
	plugins map[string]*Plugin
	log     *log.Logger
	js      *otto.Otto
	cfg     *config.Settings
	conn    *irc.Connection
}

// Walker func
func (pm *PluginManager) traversePluginDir(path string, info os.FileInfo, err error) error {
	// If errors, just stop
	if err != nil {
		return err
	}

	// Skip hidden files and folders
	if filepath.HasPrefix(info.Name(), ".") {
		return nil
	}

	if ext := strings.ToLower(filepath.Ext(info.Name())); info.Mode().IsRegular() && ext == ".js" {
		pErr := pm.LoadPlugin(path)
		if pErr != nil {
			pm.log.Printf("Skipping plugin file `%s`: %s\n", path, err)
		}
	}
	return nil
}

func (pm *PluginManager) LoadPlugins() {
	// Ditch existing plugins by redeclaring
	pm.plugins = make(map[string]*Plugin)
	// Init js env now (no need for it until we load plugins!)
	pm.js = otto.New()
	// Force a GC
	runtime.GC()

	_, err := os.Stat(pm.cfg.Irc.PluginsDir)
	if err == nil {
		// Go through each plugin file and load it
		err = filepath.Walk(pm.cfg.Irc.PluginsDir, pm.traversePluginDir)
		if err != nil {
			pm.log.Printf("Error loading plugins: %s\n", err)
		}
	} else if os.IsNotExist(err) {
		pm.log.Printf("Plugin directory does not exist: %s\n", err)
	} else {
		pm.log.Printf("Error accessing plugin directory: %s\n", err)
	}
}

func (pm *PluginManager) LoadPlugin(path string) error {
	plugin, err := ioutil.ReadFile(path)
	if err != nil {
		pm.log.Printf("Couldn't read plugin file `%s`: %s\n", path, err)
		return err
	}

	name := filepath.Base(path)
	log := log.New(os.Stdout, "["+name+"] ", log.LstdFlags)
	pm.plugins[name] = &Plugin{
		commands:  make(map[string]*pluginFunc),
		callbacks: make(map[string][]*pluginFunc),
		log:       log,
		js:        pm.js,
		cfg:       pm.cfg,
	}

	pm.js.Set("log", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
			log.Println(call.ArgumentList[0].String())
			return otto.TrueValue()
		} else {
			return otto.FalseValue()
		}
	})

	// Add in function to Register !commands
	pm.js.Set("RegisterCommand", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) >= 3 && call.ArgumentList[0].IsString() && call.ArgumentList[1].IsFunction() && call.ArgumentList[2].IsString() {
			command := call.ArgumentList[0].String()
			f := call.ArgumentList[1]
			help := call.ArgumentList[2].String()
			pm.plugins[name].SetCommand(command, f, help)
			pm.log.Printf("Registered command `%s` from plugin `%s`\n", command, name)
			return otto.TrueValue()
		} else {
			return otto.FalseValue()
		}
	})

	// Add in function to register callbacks
	pm.js.Set("RegisterCallback", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) >= 3 && call.ArgumentList[0].IsString() && call.ArgumentList[1].IsString() && call.ArgumentList[2].IsFunction() {
			eventCode := call.ArgumentList[0].String()
			callbackName := call.ArgumentList[1].String()
			f := call.ArgumentList[2]
			pm.plugins[name].AddCallback(eventCode, callbackName, f)
			pm.log.Printf("Registered callback `%s` from plugin `%s`\n", callbackName, name)
			return otto.TrueValue()
		} else {
			return otto.FalseValue()
		}
	})

	// Now we have defined the required registration commands in the JS
	// execution context, we run the plugin file contents
	_, err = pm.js.Run(string(plugin))

	// And now we remove the functions, since they aren't needed during
	// normal operations
	pm.js.Set("RegisterCommand", nil)
	pm.js.Set("RegisterCallback", nil)
	pm.js.Set("log", nil)

	if err != nil {
		pm.log.Printf("Error interpreting plugin code for plugin `%s`: %s\n", name, err)
		delete(pm.plugins, name)
		return err
	}
	return nil
}

func (pm *PluginManager) RunCallbacks(event *irc.Event) {
	pm.log.Printf("Looking for plugin callbacks for event `%s`...\n", event.Code)
	for name, plugin := range pm.plugins {
		if pm.cfg.Irc.Debug || pm.cfg.Debug {
			pm.log.Printf("Dispatching event `%s` to plugin `%s` callbacks\n", event.Code, name)
		}
		plugin.RunCallbacks(event)
	}
}

func (pm *PluginManager) RunCommands(event *irc.Event) {
	for name, plugin := range pm.plugins {
		if pm.cfg.Irc.Debug || pm.cfg.Debug {
			pm.log.Printf("Dispatching event `%s` to plugin `%s` commands\n", event.Code, name)
		}
		// Stop processing as soon as a command is found in a plugin
		if plugin.RunCommand(event) {
			return
		}
	}
}

func New(cfg *config.Settings, conn *irc.Connection) *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*Plugin),
		log:     log.New(os.Stdout, "[plugins] ", log.LstdFlags),
		cfg:     cfg,
		conn:    conn,
	}
}
