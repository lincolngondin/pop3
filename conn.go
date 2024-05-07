// Package pop3 implements an Post Ofice Protocol Version 3 (POP3) client as defined in RFC 1939.
// The package returns all server response as raw bytes, you must parse yourself in the way you like
package pop3

import (
	"bufio"
	"crypto/tls"
)

type Conn struct {
	addr      string
	conn      *tls.Conn
	tlsConfig *tls.Config
}

// Start the TCP connection with the server, the server will send an greeting in success.
func (conn *Conn) Start() (*response, error) {
	connection, err := tls.Dial("tcp", conn.addr, conn.tlsConfig)
	if err != nil {
		return nil, err
	}
    conn.conn = connection
    response, responseErr := conn.readresponse()
    if responseErr != nil {
        return nil, responseErr
    }
	return response, nil
}

// End the session with the POP3 server, this doesn't close the TCP connection, you must call Close() for that.
func (conn *Conn) Quit() (*response, error) {
    conn.sendCommand(NewCMD(commandQuit))
    response, responseErr := conn.readresponse()
    if responseErr != nil {
        return nil, responseErr
    }
    return response, nil
}

// Close the TCP connection with the server, remember that you MUST send an quit command calling conn.Quit()
// before to close the session.
func (conn *Conn) Close() error {
    return conn.conn.Close()
}

func splitresponseLines(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
    for n := 0; n < len(data)-1; n++{
        if data[n] == byte('\r') && data[n+1] == byte('\n') {
            if n+2 == len(data){
                return n+2, data[0:n+2], bufio.ErrFinalToken
            } else {
                return n+2, data[0:n+2], nil
            }
        }
    }
    return 0, nil, nil
    
}

func (conn *Conn) readresponse() (*response, error) {
    response := &response{
        Data: make([][]byte, 0, 10),
    }
    scanner := bufio.NewScanner(conn.conn)
    scanner.Split(splitresponseLines)
    for scanner.Scan() == true {
        data := scanner.Bytes()
        respLine := make([]byte, len(data))
        copy(respLine, data)
        response.Data = append(response.Data, respLine)
    }

    return response, nil
}

func (conn *Conn) sendCommand(cmd command) error {
    _, err := conn.conn.Write([]byte(cmd.GetCMD()))
    return err
}

// User command, name is an string identifying a mailbox.
func (conn *Conn) User(name string) (*response, error){
    conn.sendCommand(NewCMD(commandUser, name))
    response, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return response, nil
}

// Pass command
func (conn *Conn) Pass(password string) (*response, error){
    conn.sendCommand(NewCMD(commandPass, password))
    resp, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// APOP command, name is an string identifying a mailbox and a MD5 digest string
func (conn *Conn) Apop(name, digest string) (*response, error){
    conn.sendCommand(NewCMD(commandPass, name, digest))
    resp, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// The server returns na positive response with an line containing information for the maildrop.
func (conn *Conn) Stat() (*response, error) {
    conn.sendCommand(NewCMD(commandStat))
    response, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return response, nil
}

 
// If argument is given the server returns an line containing info of the message
// If no argument the response is multiline giving info of each message
func (conn *Conn) List(msg string) (*response, error) {
    if msg != "" {
        conn.sendCommand(NewCMD(commandList, msg))
    } else {
        conn.sendCommand(NewCMD(commandList))
    }
    response, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return response, nil
}

// The pop3 server replies with message correspoding to the given msg number, or -ERR if no such message
func (conn *Conn) Retr(msg string) (*response, error) {
    conn.sendCommand(NewCMD(commandRetr, msg))
    resp, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// The pop3 server replies with an positive response
func (conn *Conn) Noop() (*response, error) {
    conn.sendCommand(NewCMD(commandNoop))
    resp, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// The pop3 server marks the msg as deleted
func (conn *Conn) Dele(msg string) (*response, error) {
    conn.sendCommand(NewCMD(commandDele, msg))
    resp, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return resp, nil
}


// The pop3 server unmark all message marked as deleted
func (conn *Conn) Rset() (*response, error) {
    conn.sendCommand(NewCMD(commandRset))
    resp, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// The pop3 server send the headers of the message msg the line separating from the body and the 
// number of lines n of the body of the indicated message
func (conn *Conn) Top(msg, n string) (*response, error) {
    conn.sendCommand(NewCMD(commandTop, msg, n))
    resp, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// if argument was given the server issue an line containing info about that message
// if no argument, issue an multiline response containing info of each message in maildrop
func (conn *Conn) Uidl(msg string) (*response, error) {
    if msg != "" {
        conn.sendCommand(NewCMD(commandTop, msg))
    } else {
        conn.sendCommand(NewCMD(commandTop))
    }
    resp, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// Execute arbitrary command
func (conn *Conn) Exec(cmd command) (*response, error) {
    conn.sendCommand(cmd)
    resp, err := conn.readresponse()
    if err != nil {
        return nil, err
    }
    return resp, nil

}

// Create new connection object, you must call Start() to init the connection
func NewConn(addr string, config *tls.Config) *Conn {
	return &Conn{
		addr:   addr,
		tlsConfig: config,
	}
}
