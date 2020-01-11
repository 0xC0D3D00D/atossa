package main

type CommandFlag string

const (
	CommandFlagWrite         CommandFlag = "write"
	CommandFlagReadonly      CommandFlag = "readonly"
	CommandFlagDenyOOM       CommandFlag = "denyoom"
	CommandFlagAdmin         CommandFlag = "admin"
	CommandFlagPubSub        CommandFlag = "pubsub"
	CommandFlagNoscript      CommandFlag = "noscript"
	CommandFlagRandom        CommandFlag = "random"
	CommandFlagSortForScript CommandFlag = "sort_for_script"
	CommandFlagLoading       CommandFlag = "loading"
	CommandFlagStale         CommandFlag = "stale"
	CommandFlagSkipMonitor   CommandFlag = "skip_monitor"
	CommandFlagAsking        CommandFlag = "asking"
	CommandFlagFast          CommandFlag = "fast"
	CommandFlagMovableKeys   CommandFlag = "movablekeys"
)

type command struct {
	name        string
	arity       int8
	flags       []CommandFlag
	firstKeyPos int8
	lastKeyPos  int8
	stepCount   int8
	handler     commandHandler
}

func (c command) Slice() []interface{} {
	flags := []interface{}{}
	for _, flag := range c.flags {
		flags = append(flags, string(flag))
	}
	return []interface{}{
		c.name,
		c.arity,
		flags,
		c.firstKeyPos,
		c.lastKeyPos,
		c.stepCount,
	}
}

var commandMap = map[string]command{
	"PING": command{
		name:  "ping",
		arity: -1,
		flags: []CommandFlag{
			CommandFlagFast,
			CommandFlagStale,
		},
		firstKeyPos: 0,
		lastKeyPos:  0,
		stepCount:   0,
		handler:     ping,
	},
	"COMMAND": command{
		name:  "command",
		arity: 0,
		flags: []CommandFlag{
			CommandFlagRandom,
			CommandFlagLoading,
			CommandFlagStale,
		},
		firstKeyPos: 0,
		lastKeyPos:  0,
		stepCount:   0,
	},
	"INFO": command{
		name:  "info",
		arity: -1,
		flags: []CommandFlag{
			CommandFlagRandom,
			CommandFlagLoading,
			CommandFlagStale,
		},
		firstKeyPos: 0,
		lastKeyPos:  0,
		stepCount:   0,
		handler:     info,
	},
	"GET": command{
		name:  "get",
		arity: 2,
		flags: []CommandFlag{
			CommandFlagReadonly,
			CommandFlagFast,
		},
		firstKeyPos: 1,
		lastKeyPos:  1,
		stepCount:   1,
		handler:     get,
	},
	"SET": command{
		name:  "set",
		arity: -3,
		flags: []CommandFlag{
			CommandFlagWrite,
			CommandFlagDenyOOM,
		},
		firstKeyPos: 1,
		lastKeyPos:  1,
		stepCount:   1,
		handler:     set,
	},
	"SEQ": command{
		name:  "seq",
		arity: 1,
		flags: []CommandFlag{
			CommandFlagWrite,
			CommandFlagDenyOOM,
		},
		firstKeyPos: 1,
		lastKeyPos:  1,
		stepCount:   1,
		handler:     seq,
	},
	"KEYS": command{
		name:  "keys",
		arity: 2,
		flags: []CommandFlag{
			CommandFlagReadonly,
			CommandFlagSortForScript,
		},
		firstKeyPos: 0,
		lastKeyPos:  0,
		stepCount:   0,
		handler:     keys,
	},
	"LPUSH": command{
		name:  "lpush",
		arity: -3,
		flags: []CommandFlag{
			CommandFlagWrite,
			CommandFlagDenyOOM,
			CommandFlagFast,
		},
		firstKeyPos: 1,
		lastKeyPos:  1,
		stepCount:   1,
		handler:     lpush,
	},
	"RPUSH": command{
		name:  "rpush",
		arity: -3,
		flags: []CommandFlag{
			CommandFlagWrite,
			CommandFlagDenyOOM,
			CommandFlagFast,
		},
		firstKeyPos: 1,
		lastKeyPos:  1,
		stepCount:   1,
		handler:     rpush,
	},
	"LPOP": command{
		name:  "lpop",
		arity: 2,
		flags: []CommandFlag{
			CommandFlagWrite,
			CommandFlagFast,
		},
		firstKeyPos: 1,
		lastKeyPos:  1,
		stepCount:   1,
		handler:     lpop,
	},
	"RPOP": command{
		name:  "rpop",
		arity: 2,
		flags: []CommandFlag{
			CommandFlagWrite,
			CommandFlagFast,
		},
		firstKeyPos: 1,
		lastKeyPos:  1,
		stepCount:   1,
		handler:     rpop,
	},
	"LSET": command{
		name:  "lset",
		arity: 4,
		flags: []CommandFlag{
			CommandFlagFast,
			CommandFlagWrite,
			CommandFlagDenyOOM,
		},
		firstKeyPos: 1,
		lastKeyPos:  1,
		stepCount:   1,
		handler:     lset,
	},
	"LINDEX": command{
		name:  "lindex",
		arity: 3,
		flags: []CommandFlag{
			CommandFlagReadonly,
			CommandFlagFast,
		},
		firstKeyPos: 1,
		lastKeyPos:  1,
		stepCount:   1,
		handler:     lindex,
	},
}
