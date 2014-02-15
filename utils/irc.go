package utils

import (
	"fmt"
	"github.com/darkliquid/go-ircevent"
	"strings"
)

func IRCAction(conn *irc.Connection, channel, action string) {
	conn.Privmsg(channel, fmt.Sprintf("\001ACTION %s\001", action))
}

func IRCInvite(conn *irc.Connection, nick, channel string) {
	conn.SendRawf("INVITE %s %s", nick, channel)
}

func IRCOper(conn *irc.Connection, user, pass string) {
	conn.SendRawf("OPER %s %s", user, pass)
}

func IRCAway(conn *irc.Connection, message ...string) {
	msg := strings.Join(message, " ")
	if msg != "" {
		msg = " :" + msg
	}
	conn.SendRawf("AWAY%s", msg)
}

func IRCTopic(conn *irc.Connection, channel string, topic ...string) {
	msg := strings.Join(topic, " ")
	if msg != "" {
		msg = " :" + msg
	}
	conn.SendRawf("TOPIC %s%s", channel, msg)
}
