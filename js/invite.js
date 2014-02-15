RegisterCommand("invite", function() {
	var args = this.event.message.split(" "),
		source = this.event.args[0],
		cmd = args.shift(),
		nick = this.event.nick,
		cfg = GetConfig(),
		privs = IRC.GetPrivs(source, nick);

	if(privs && (privs.Owner || privs.Admin || privs.Op) && source == cfg.Irc.StaffChannel) {
		IRC.Invite(args[0], cfg.Irc.StaffChannel)
	}
}, "invites a user to the staff channel");