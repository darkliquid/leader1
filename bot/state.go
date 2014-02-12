package bot

import (
	"github.com/darkliquid/go-ircevent"
)

type StateTracker struct {
	channels map[string]*Channel
	nicks    map[string]*Nick
	me       *Nick
	conn     *irc.Connection
}

func (st *StateTracker) Joined(event *irc.Event) {
	channel := st.channels[event.Arguments[0]]

	if channel == nil {
		st.channels[event.Arguments[0]] = &Channel{
			Name: event.Arguments[0],
		}
	}
}

func (bot *Bot) InitStateTracking() {
	bot.state = &StateTracker{
		channels: make(map[string]*Channel),
		nicks:    make(map[string]*Nick),
		conn:     bot.conn,
	}
	bot.conn.AddCallback("JOIN", bot.state.Joined)
	/*bot.conn.AddCallback("KICK", bot.state.Kicked)
	bot.conn.AddCallback("NICK", bot.state.NickChanged)
	bot.conn.AddCallback("PART", bot.state.Parted)
	bot.conn.AddCallback("QUIT", bot.state.Quitted)
	bot.conn.AddCallback("TOPIC", bot.state.TopicSet)
	bot.conn.AddCallback("311", bot.state.WhoisReply)
	bot.conn.AddCallback("324", bot.state.ModeReply)
	bot.conn.AddCallback("332", bot.state.TopicReply)
	bot.conn.AddCallback("352", bot.state.WhoReply)
	bot.conn.AddCallback("353", bot.state.NamesReply)
	bot.conn.AddCallback("671", bot.state.SSLWhoisReply)*/
}