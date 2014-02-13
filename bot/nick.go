package bot

type NickModes struct {
	Bot, Invisible, Oper, WallOps, HiddenHost, SSL bool
}

type Nick struct {
	Nick, User, Host, Name string
	Modes                  NickModes
	Channels               map[string]*ChannelPrivileges
}