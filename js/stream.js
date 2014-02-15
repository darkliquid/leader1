RegisterCommand("listeners", function(){
	var args = this.event.message.split(" "),
		source = this.event.args[0],
		cmd = args.shift(),
		nick = this.event.nick,
		stats = UTILS.GetShoutcastStats();

	if(stats.UniqueListeners == 0){
		IRC.Privmsg(source, nick+": there are currently no listeners on the stream :'(");
	} else if(stats.UniqueListeners == 1) {
		IRC.Privmsg(source, nick+": there is currently 1 listener on the stream");
	} else if(stats.UniqueListeners == stats.MaxListeners) {
		IRC.Privmsg(source, nick+": holy crap! We have maxed out at "+stats.UniqueListeners+" listeners on the stream!!!");
	} else {
		IRC.Privmsg(source, nick+": there are currently "+stats.UniqueListeners+" listeners on the stream");
	}
	
}, "returns the number of listeners on the stream");

RegisterCommand("g3song", function(){
	var args = this.event.message.split(" "),
		source = this.event.args[0],
		cmd = args.shift(),
		nick = this.event.nick,
		stats = UTILS.GetShoutcastStats();

	IRC.Privmsg(source, nick+": current song is - `"+stats.SongTitle+"`");
}, "returns the currently playing song");

RegisterCommand("+", function(){
	var args = this.event.message.split(" "),
		source = this.event.args[0],
		cmd = args.shift(),
		nick = this.event.nick,
		stats = UTILS.GetShoutcastStats(),
		cfg = GetConfig();

	if(UTILS.LikeTrack(nick, stats.SongTitle)) {
		IRC.Privmsg(source, nick+": thanks, you liked `"+stats.SongTitle+"`");
		IRC.Privmsg(cfg.Irc.StaffChannel, "LIKE: "+nick+" liked `"+stats.SongTitle+"`");
	} else {
		IRC.Privmsg(cfg.Irc.StaffChannel, "LIKE: "+nick+" liked `"+args.join(" ")+"` (database insert failed)");
	}
}, "likes the currently playing song");

RegisterCommand("-", function(){
	var args = this.event.message.split(" "),
		source = this.event.args[0],
		cmd = args.shift(),
		nick = this.event.nick,
		stats = UTILS.GetShoutcastStats(),
		cfg = GetConfig();

	if(UTILS.HateTrack(nick, stats.SongTitle)) {
		IRC.Privmsg(source, nick+": thanks, you hated `"+stats.SongTitle+"`");
		IRC.Privmsg(cfg.Irc.StaffChannel, "HATE: "+nick+" hated `"+stats.SongTitle+"`");
	} else {
		IRC.Privmsg(cfg.Irc.StaffChannel, "HATE: "+nick+" hated `"+args.join(" ")+"` (database insert failed)");
	}
}, "hates the currently playing song");

RegisterCommand("request", function(){
	var args = this.event.message.split(" "),
		source = this.event.args[0],
		cmd = args.shift(),
		nick = this.event.nick,
		stats = UTILS.GetShoutcastStats(),
		cfg = GetConfig();

	if(args.length == 0) {
		IRC.Privmsg(source, nick+": usage is - !request [track name]");
	}

	if(UTILS.Request(nick, args.join(" "))) {
		IRC.Privmsg(source, nick+": thanks, you requested `"+args.join(" ")+"`");
		IRC.Privmsg(cfg.Irc.StaffChannel, "REQ: "+nick+" requests `"+args.join(" ")+"`");
	} else {
		IRC.Privmsg(cfg.Irc.StaffChannel, "REQ: "+nick+" requests `"+args.join(" ")+"` (database insert failed)");
	}
}, "logs a music request");