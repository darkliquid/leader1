RegisterCallback("PRIVMSG", "McIRCBot Redispatcher", function() {
	// If the original is the minecraft bot
	if(this.event.nick === "McIRCBot") {
		var matches = this.event.message.match(/^(<(.+)>) (.*)/);
		// If we can do a successful regex match, we rewrite the 
		// message into a native IRC event and redispatch it
		if(matches.length === 4) {
			IRC.Redispatch(
				this.event.code,    // Copy event code
				this.event.raw,     // Copy event raw
				matches[2],         // User minecraft nick as event nick
				this.event.host,    // Copy event host
				this.event.source,  // Copy event source
				this.event.user,    // Copy event user
				this.event.args[0], // Copy event arg[0] as first argument (for privmsg, this is the destination channel)
				matches[3]          // Set the message posted by the minecraft user as the second event argument
			)
		}
	}
});