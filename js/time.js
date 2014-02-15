RegisterCommand("time", function(){
	var args = this.event.message.split(" "),
		source = this.event.args[0],
		cmd = args.shift(),
		nick = this.event.nick,
		date = new Date();
	IRC.Privmsg(source, nick+": the time is " + date)
}, "returns the current time for the bot");