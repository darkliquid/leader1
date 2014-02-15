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