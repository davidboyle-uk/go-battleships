package tcp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rockwell-uk/go-logger/logger"
)

const (
	PARSE_ERROR = "Parse Error"
	PROTO_DELIM = "|"
)

// Parse inbound packet.
func ParseMessage(m string) (Proto, error) {
	var action, body string
	var player int

	if m == "" {
		return Proto{}, tcpError{
			s: PARSE_ERROR,
			v: "unable to parse packet",
		}
	}

	payload := strings.SplitN(strings.TrimRight(m, "\n"), PROTO_DELIM, 3)

	l := len(payload)
	action = payload[0]
	if l > 1 {
		i, err := strconv.Atoi(payload[1])
		if err != nil {
			return Proto{}, tcpError{
				s: PARSE_ERROR,
				v: "unable to parse packet (player)",
			}
		}
		player = i
	}
	if l > 2 {
		body = payload[2]
	}

	parsed := Proto{
		Action: action,
		Player: player,
		Body:   body,
	}

	logger.Log(
		logger.LVL_INTERNAL,
		fmt.Sprintf("parsed: %v", parsed),
	)

	return parsed, nil
}
