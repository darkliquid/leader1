package bot

import(
	"github.com/darkliquid/go-ircevent"
	"strings"
	"fmt"
	"github.com/darkliquid/leader1/utils"
	"time"
)

func (bot *Bot) RunBuiltinCommands(event *irc.Event) {
	args := strings.Split(strings.TrimSpace(event.Message()), " ")
	command := args[0]

	// Bin the command from the arg list
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}

	// Get the current channel privileges for the nick sending this command
	privs, ok := bot.state.GetPrivs(event.Arguments[0], event.Nick)
	
	// Commands must be run by known users
	if ok {
		switch {
		case command == "!reload":
			if privs.Owner || privs.Admin || privs.Op {
				utils.IRCAction(bot.conn, event.Arguments[0], "reloads it's plugins")
				bot.pm.InitJS()
			} else {
				utils.IRCAction(bot.conn, event.Arguments[0], fmt.Sprintf("slaps %s's hands away from the op only controls", event.Nick))
			}
		case command == "!ping":
			bot.conn.Privmsg(event.Arguments[0], fmt.Sprintf("%s: PONG!", event.Nick))
		case command == "!quit":
			if privs.Owner || privs.Admin || privs.Op {
				bot.Quit()
			} else {
				utils.IRCAction(bot.conn, event.Arguments[0], fmt.Sprint("slaps %s's hands away from the op only controls", event.Nick))
			}
		case command == "!help":
			if len(args) == 0 {
				bot.ShowCommandList(event.Arguments[0], event.Nick)
			} else {
				bot.ShowCommandHelp(event.Arguments[0], event.Nick, args[0])
			}
		}
	}
}

// Handler for reclaiming a stolen nick
func (bot *Bot) ReclaimNick(event *irc.Event) {
	if thief := bot.state.GetNick(bot.cfg.Irc.Nick); thief != nil {
		// Recover nick from thieves
		bot.conn.Privmsg("NickServ", fmt.Sprintf("RECOVER %s %s", bot.cfg.Irc.Nick, bot.cfg.Irc.NickPass))
		time.Sleep(time.Second)
		bot.conn.Privmsg("NickServ", fmt.Sprintf("RELEASE %s %s", bot.cfg.Irc.Nick, bot.cfg.Irc.NickPass))
	}
	bot.SetBotState(event)
}

// Automatically give voice to users in channels in which the bot has Op
func (bot *Bot) AutoVoice(event *irc.Event) {
	if bot.cfg.Irc.AutoVoice {
		time.Sleep(time.Second) // Wait a second before we bother

		privs, ok := bot.state.GetPrivs(event.Arguments[0], event.Nick)

		// No need to autovoice
		if ok && (privs.Owner || privs.Admin || privs.Op || privs.HalfOp || privs.Voice) {
			return
		}

		// Set mode
		bot.conn.Mode(event.Arguments[0], fmt.Sprintf("+v %s", event.Nick))
	}
}

// Function for setting up the botstate
func (bot *Bot) SetBotState(event *irc.Event) {
	if ghost := bot.state.GetNick(bot.cfg.Irc.Nick); ghost != nil && ghost != bot.state.Me() {
		// GHOST the old nick
		bot.conn.Privmsg("NickServ", fmt.Sprintf("GHOST %s %s", bot.cfg.Irc.Nick, bot.cfg.Irc.NickPass))
	}

	// Set up the nick
	bot.conn.Nick(bot.cfg.Irc.Nick)

	// Identify as the nick owner
	bot.conn.Privmsg("NickServ", fmt.Sprintf("IDENTIFY %s", bot.cfg.Irc.NickPass))

	// Tell IRC I'm a bot
	bot.conn.Mode(bot.cfg.Irc.Nick, "+B")
}

// Function for re-joining channels
func (bot *Bot) JoinChannels(event *irc.Event) {
	time.Sleep(time.Second) // Wait a second before we bother

	if _, ok := bot.state.GetPrivs(bot.cfg.Irc.NormalChannel, bot.state.Me().Nick) ; !ok {
		bot.conn.Join(bot.cfg.Irc.NormalChannel)
	}

	if _, ok := bot.state.GetPrivs(bot.cfg.Irc.StaffChannel, bot.state.Me().Nick) ; !ok {
		bot.conn.Join(bot.cfg.Irc.StaffChannel)
	}
}

// Print out the commands available
func (bot *Bot) ShowCommandList(source, nick string) {
	var commands []string = make([]string, 0)
	commands = append(commands, "!reload", "!ping", "!quit", "!help")

	for cmd, _ := range bot.pm.CommandHelp() {
		commands = append(commands, cmd)
	}
	bot.conn.Privmsg(source, fmt.Sprintf("%s: available commands are - %s", nick, strings.Join(commands, ", ")))
}

// Print out the commands available
func (bot *Bot) ShowCommandHelp(source, nick, cmd string) {
	message := "unknown command, run `!help` to see what commands are available."
	switch {
	case cmd == "!ping":
		message = fmt.Sprintf("makes `%s` reply with PONG!", bot.state.Me().Nick)
	case cmd == "!reload":
		message = "reloads the bots plugins and config file"
	case cmd == "!quit":
		message = fmt.Sprintf("makes `%s` quit IRC", bot.state.Me().Nick)
	case cmd == "!help":
		message = "shows this message, smart ass"
	default:
		if help, ok := bot.pm.CommandHelp()[cmd] ; ok {
			message = help
		}
	}

	bot.conn.Privmsg(source, fmt.Sprintf("%s: %s - %s", nick, cmd, message))
}