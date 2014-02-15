RegisterCommand("urban", function(){
	var args = this.event.message.split(" "),
		source = this.event.args[0],
		cmd = args.shift(),
		nick = this.event.nick,
		url = "http://www.urbandictionary.com/define.php?term=" + encodeURIComponent(args.join(" "));
	
	IRC.Privmsg(source, nick+": Urban dictionary says - " + url);
}, "returns an urban dictionary search url for the given query");