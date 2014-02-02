leader-1
========

Go IRC Bot, named after Leader-1 from GoBots.


## Configuration

Leader-1 expects there to be a leader-1.json in the current directory by default. The config file should be structured in the following way. All settings are required.

```json
{
    "irc": {
        "host": "irc.domain.tld",
        "port": "6667",
        "nick": "mightygobot",
        "nickpass": "MyIdentifyPassword",
        "pass": "ServerPassword",
        "ssl": false,
        "normal_channel": "#channel",
        "staff_channel": "#staff",
        "timeout": 30
    },
    "db": {
    	"dsn": "user:pass@tcp(localhost)/database"
    }
}
```
