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

type cmd interface {
	// return the command with the arguments and the CRLF in the end
	GetCMD() string
}

type CMD struct {
	keyword   string
	arguments []string
}

func NewCMD(keyw string, args ...string) cmd {
	return &CMD{
		keyword:   keyw,
		arguments: args,
	}
}

func (cmd *CMD) GetCMD() string {
	strBuilder := &strings.Builder{}
	strBuilder.WriteString(cmd.keyword)
	for _, arg := range cmd.arguments {
		strBuilder.WriteString(space + arg)
	}
	strBuilder.WriteString(crlf)
	return strBuilder.String()
}

var quitCMD CMD = CMD{
	keyword:   commandQuit,
	arguments: nil,
}
