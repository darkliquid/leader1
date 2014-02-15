RegisterCommand("debug", function() {
	this.log(JSON.stringify(IRC.Nicks()))
}, "prints out debug info");

log("Loaded")