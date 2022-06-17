package irc

import (
	"strings"
)

type Message struct {
	Prefix   string
	Command  string
	Params   string
	Trailing string
	String   string
}

func (c *Message) Parse(line string) {
	line = strings.TrimSuffix(line, "\r")
	line = strings.TrimSuffix(line, "\r\n")
	orig := line
	c.String = orig
	// Prefix
	if line[0] == ':' {
		i := strings.Index(line, " ")
		c.Prefix = line[1:i]
		line = line[i+1:]
	}
	// Command
	i := strings.Index(line, " ")
	if -1 == i {
		i = len(line)
	}
	c.Command = line[0:i]
	line = line[i:]
	// Params
	i = strings.Index(line, " :")
	if -1 == i {
		i = len(line)
	}
	if i != 0 {
		c.Params = line[1:i]
	}
	// Trailing
	if len(line)-i > 2 {
		c.Trailing = line[i+2:]
	}
}
