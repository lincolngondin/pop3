package pop3

import (
	"strings"
)

const (
	space = " "
	crlf  = "\r\n"
)

const (
	commandApop = "APOP"
	commandDele = "DELE"
	commandList = "LIST"
	commandNoop = "NOOP"
	commandPass = "PASS"
	commandQuit = "QUIT"
	commandRetr = "RETR"
	commandRset = "RSET"
	commandStat = "STAT"
	commandTop  = "TOP"
	commandUidl = "UIDL"
	commandUser = "USER"
)

type command interface {
	// return the command with the arguments and the CRLF in the end
	GetCMD() string
}

type cmd struct {
	keyword   string
	arguments []string
}

// Create an arbitrary command with keyword key and arguments args.
// Depends of the server if the implementation of your command will be supported.
func NewCMD(keyw string, args ...string) *cmd {
	return &cmd{
		keyword:   keyw,
		arguments: args,
	}
}

// Return the command in the way that must be send for the server
// in the format defined in RFC 1939 section 3
func (cmd *cmd) GetCMD() string {
	strBuilder := &strings.Builder{}
	strBuilder.WriteString(cmd.keyword)
	for _, arg := range cmd.arguments {
		strBuilder.WriteString(space + arg)
	}
	strBuilder.WriteString(crlf)
	return strBuilder.String()
}
