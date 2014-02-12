package bot

type NickModes struct {
	Bot, Invisible, Oper, WallOps, HiddenHost, SSL bool
}

type Nick struct {
	Nick, Ident, Host, Name string
	Modes                   *NickModes
	Channels                []string
}