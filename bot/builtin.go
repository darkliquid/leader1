package bot

import(
	"github.com/darkliquid/go-ircevent"
	"strings"
)

func (bot *Bot) RunBuiltinCommands(event *irc.Event) {
	if strings.TrimSpace(event.Message()) == "!reload" {
		bot.pm.LoadPlugins()
	}
}