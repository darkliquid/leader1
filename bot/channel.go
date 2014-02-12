package bot

type ChannelModes struct {
	// MODE +p, +s, +t, +n, +m
	Private, Secret, ProtectedTopic, NoExternalMsg, Moderated bool

	// MODE +i, +O, +z
	InviteOnly, OperOnly, SSLOnly bool

	// MODE +r, +Z
	Registered, AllSSL bool

	// MODE +k
	Key string

	// MODE +l
	Limit int
}

type ChannelPrivileges struct {
	// MODE +q, +a, +o, +h, +v
	Owner, Admin, Op, HalfOp, Voice bool	
}

type Channel struct {
	Name, Topic string
	Modes       *ChannelModes
	Nicks       map[string]*ChannelPrivileges
}