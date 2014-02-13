package bot

import (
	"github.com/darkliquid/go-ircevent"
	"sync"
	"strings"
)

type StateTracker struct {
	channels map[string]*Channel
	nicks    map[string]*Nick
	me       *Nick
	conn     *irc.Connection
	mutex    sync.Mutex
}

func (st *StateTracker) joined(event *irc.Event) {
	// Wait until ready
	st.mutex.Lock()

	channel := st.channels[event.Arguments[0]]
	nick := st.nicks[event.Nick]

	// Haven't seen this channel before
	if channel == nil {
		// Create a channel object
		channel = &Channel{
			Name: event.Arguments[0],
			Nicks: make(map[string]*ChannelPrivileges),
		}

		// Put it in the channels map
		st.channels[event.Arguments[0]] = channel

		// Get some initial info about it
		st.conn.Mode(channel.Name)
		st.conn.Who(channel.Name)
	}

	// Not seen this nick before
	if nick == nil {
		// Create a nick object
		nick = &Nick{
			Nick: event.Nick,
			User: event.User,
			Host: event.Host,
			Channels: make(map[string]*ChannelPrivileges),
		}

		// Put it in the nicks map
		st.nicks[event.Nick] = nick

		// Get some inital info about it
		st.conn.Who(nick.Nick)
	}

	// Associate the nick with the channel
	st.associate(event.Nick, event.Arguments[0])

	// Ready for next event
	st.mutex.Unlock()
}

func (st *StateTracker) kicked(event *irc.Event) {
	st.mutex.Lock()
	st.disassociate(event.Nick, event.Arguments[0])
	st.mutex.Unlock()
}

func (st *StateTracker) nickChanged(event *irc.Event) {
	st.mutex.Lock()
	st.changeNick(event.Nick, event.Arguments[0])
	st.mutex.Unlock()
}

func (st *StateTracker) parted(event *irc.Event) {
	st.mutex.Lock()
	st.disassociate(event.Nick, event.Arguments[0])
	st.mutex.Unlock()
}

func (st *StateTracker) quitted(event *irc.Event) {
	st.mutex.Lock()
	st.deleteNick(event.Nick)
	st.mutex.Unlock()
}

func (st *StateTracker) topicSet(event *irc.Event) {
	st.mutex.Lock()
	st.setTopic(event.Arguments[0], event.Arguments[1])
	st.mutex.Unlock()
}

func (st *StateTracker) whoisReply(event *irc.Event) {
	st.mutex.Lock()
	nick := st.nicks[event.Arguments[1]]
	if nick != nil && nick != st.me {
		nick.User = event.Arguments[2]
		nick.Host = event.Arguments[3]
		nick.Name = event.Arguments[5]
	}
	st.mutex.Unlock()
}

func (st *StateTracker) topicReply(event *irc.Event) {
	st.mutex.Lock()
	if channel := st.channels[event.Arguments[1]] ; channel != nil {
		st.setTopic(channel.Name, event.Arguments[2])
	}
	st.mutex.Unlock()
}

func (st *StateTracker) whoReply(event *irc.Event) {
	st.mutex.Lock()
	if nick, ok := st.nicks[event.Arguments[5]]; ok {
		nick.User = event.Arguments[2]
		nick.Host = event.Arguments[3]
		if idx := strings.Index(event.Arguments[6], "*"); idx != -1 {
			nick.Modes.Oper = true
		}
		if idx := strings.Index(event.Arguments[6], "H"); idx != -1 {
			nick.Modes.Invisible = true
		}
	}
	st.mutex.Unlock()
}

func (st *StateTracker) modeReply(event *irc.Event) {
	st.mutex.Lock()
	if channel, ok := st.channels[event.Arguments[1]] ; ok {
		channel.ParseModes(event.Arguments[2], event.Arguments[3:]...)
	}
	st.mutex.Unlock()
}

func (st *StateTracker) namesReply(event *irc.Event) {
	st.mutex.Lock()
	if channel, ok := st.channels[event.Arguments[2]] ; ok {
		names := strings.Split(strings.TrimSpace(event.Arguments[len(event.Arguments)-1]), " ")
		for _, name := range names {
			switch priv := name[0]; priv {
			case '~', '&', '@', '%', '+':
				name = name[1:]
				fallthrough
			default:
				nick := st.nicks[name]

				if nick == nil {
					st.nicks[name] = &Nick{
						Nick: name,
						Channels: make(map[string]*ChannelPrivileges),
					}
				}
				privs, ok := channel.Nicks[name]
				if !ok {
					privs = st.associate(name, channel.Name)
				}

				switch priv {
				case '~':
					privs.Owner = true
				case '&':
					privs.Admin = true
				case '@':
					privs.Op = true
				case '%':
					privs.HalfOp = true
				case '+':
					privs.Voice = true
				}
			}
		}
	}
	st.mutex.Unlock()
}

func (st *StateTracker) whoisReplySSL(event *irc.Event) {
	st.mutex.Lock()
	if nick, ok := st.nicks[event.Arguments[1]] ; ok && nick != st.me {
		nick.User = event.Arguments[2]
		nick.Host = event.Arguments[3]
		nick.Name = event.Arguments[5]
		nick.Modes.SSL = true
	}
	st.mutex.Unlock()
}


func (st *StateTracker) associate(nick, channel string) *ChannelPrivileges {
	channelObj := st.channels[channel]
	nickObj := st.nicks[nick]

	// Haven't seen this channel before
	if channelObj == nil {
		// Create a channel object
		channelObj = &Channel{
			Name: channel,
			Nicks: make(map[string]*ChannelPrivileges),
		}

		// Put it in the channels map
		st.channels[channel] = channelObj

		// Get some initial info about it
		st.conn.Mode(channelObj.Name)
		st.conn.Who(channelObj.Name)
	}

	// Not seen this nick before
	if nickObj == nil {
		// Create a nick object
		nickObj = &Nick{
			Nick: nick,
			Channels: make(map[string]*ChannelPrivileges),
		}

		// Put it in the nicks map
		st.nicks[nick] = nickObj

		// Get some inital info about it
		st.conn.Who(nickObj.Nick)
	}

	privs := new(ChannelPrivileges)
	nickObj.Channels[channel] = privs
	channelObj.Nicks[nick] = privs
	return privs
}

func (st *StateTracker) disassociate(nick, channel string) {
	if channelObj, ok := st.channels[channel]; ok {
		delete(channelObj.Nicks, nick)
	}
	if nickObj, ok := st.nicks[nick]; ok {
		delete(nickObj.Channels, channel)	
	}
}

func (st *StateTracker) changeNick(oldNick, newNick string) {
	if nick, ok := st.nicks[oldNick]; ok {
		nick.Nick = newNick
		st.nicks[newNick] = nick
		delete(st.nicks, oldNick)
	}
	for _, channel := range st.channels {
		if privs, ok := channel.Nicks[oldNick]; ok {
			channel.Nicks[newNick] = privs
			delete(channel.Nicks, oldNick)
		}
	}
}

func (st *StateTracker) deleteNick(nick string) {
	if _, ok := st.nicks[nick]; ok {
		delete(st.nicks, nick)
	}
	for _, channel := range st.channels {
		if _, ok := channel.Nicks[nick]; ok {
			delete(channel.Nicks, nick)
		}
	}
}

func (st *StateTracker) setTopic(channel, topic string) {
	channelObj := st.channels[channel]

	// Haven't seen this channel before
	if channelObj == nil {
		// Create a channel object
		channelObj = &Channel{
			Name: channel,
			Nicks: make(map[string]*ChannelPrivileges),
		}

		// Put it in the channels map
		st.channels[channel] = channelObj

		// Get some initial info about it
		st.conn.Mode(channelObj.Name)
		st.conn.Who(channelObj.Name)
	}

	channelObj.Topic = topic
}

func (bot *Bot) InitStateTracking() {
	bot.state = &StateTracker{
		channels: make(map[string]*Channel),
		nicks:    make(map[string]*Nick),
		conn:     bot.conn,
		me:       &Nick{
			Name: bot.cfg.Irc.Nick,
			Channels: make(map[string]*ChannelPrivileges),
		},
	}
	bot.state.nicks[bot.cfg.Irc.Nick] = bot.state.me
	bot.conn.AddCallback("JOIN", bot.state.joined)
	bot.conn.AddCallback("KICK", bot.state.kicked)
	bot.conn.AddCallback("NICK", bot.state.nickChanged)
	bot.conn.AddCallback("PART", bot.state.parted)
	bot.conn.AddCallback("QUIT", bot.state.quitted)
	bot.conn.AddCallback("TOPIC", bot.state.topicSet)
	bot.conn.AddCallback("311", bot.state.whoisReply)
	bot.conn.AddCallback("324", bot.state.modeReply)
	bot.conn.AddCallback("332", bot.state.topicReply)
	bot.conn.AddCallback("352", bot.state.whoReply)
	bot.conn.AddCallback("353", bot.state.namesReply)
	bot.conn.AddCallback("671", bot.state.whoisReplySSL)
}