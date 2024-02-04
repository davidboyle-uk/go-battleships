package tcp

import (
	"fmt"
)

type Proto struct {
	Action string
	Player int
	Body   string
}

func (p Proto) String() string {
	return fmt.Sprintf("%s%s%v%s%s\n",
		p.Action, PROTO_DELIM,
		p.Player, PROTO_DELIM,
		p.Body)
}

type tcpError struct {
	s string
	v string
}

func (n tcpError) Error() string {
	return fmt.Sprintf("%s: %s", n.s, n.v)
}
