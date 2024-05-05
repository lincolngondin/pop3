package main

import (
	"bufio"
	"crypto/tls"
	"log"
)


type Conn struct {
	addr      string
	conn      *tls.Conn
	tlsConfig *tls.Config
}

// Start the TCP connection with the server, the server will send an greeting in success
func (conn *Conn) Start() (*Response, error) {
	connection, err := tls.Dial("tcp", conn.addr, conn.tlsConfig)
	if err != nil {
		return nil, err
	}
    conn.conn = connection
    response, responseErr := conn.readResponse(false)
    if responseErr != nil {
        return nil, responseErr
    }
	return response, nil
}

func splitResponseLines(data []byte, atEOF bool) (int, []byte, error) {
    log.Print(atEOF, data, string(data))
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

func (conn *Conn) readResponse(multiline bool) (*Response, error) {
    response := &Response{
        info: make([][]byte, 0, 10),
    }
    scanner := bufio.NewScanner(conn.conn)
    scanner.Split(splitResponseLines)
    for scanner.Scan() == true {
        data := scanner.Bytes()
        respLine := make([]byte, len(data))
        copy(respLine, data)
        response.info = append(response.info, respLine)
    }

    return response, nil
}

func (conn *Conn) sendCommand(cmd cmd) error {
    _, err := conn.conn.Write([]byte(cmd.GetCMD()))
    return err
}

func (conn *Conn) User(name string) (*Response, error){
    conn.sendCommand(NewCMD(commandUser, name))
    response, err := conn.readResponse(false)
    if err != nil {
        return nil, err
    }
    return response, nil
}

func (conn *Conn) Pass(password string) (*Response, error){
    conn.sendCommand(NewCMD(commandPass, password))
    resp, err := conn.readResponse(false)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

func (conn *Conn) Quit() (*Response, error) {
    conn.sendCommand(NewCMD(commandQuit))
    response, responseErr := conn.readResponse(false)
    if responseErr != nil {
        return nil, responseErr
    }
    return response, nil
}

// The server returns na positive response with an line containing information for the maildrop
func (conn *Conn) Stat() (*Response, error) {
    conn.sendCommand(NewCMD(commandStat))
    response, err := conn.readResponse(false)
    if err != nil {
        return nil, err
    }
    return response, nil
}

/* 
If argument is given the server returns an line containing info of the message
If no argument the response is multiline giving info of each message
*/
func (conn *Conn) List(msg string) (*Response, error) {
    if msg != "" {
        conn.sendCommand(NewCMD(commandList, msg))
    } else {
        conn.sendCommand(NewCMD(commandList))
    }
    response, err := conn.readResponse(msg != "")
    if err != nil {
        return nil, err
    }
    return response, nil
}

func (conn *Conn) Retr(msg string) (*Response, error) {
    conn.sendCommand(NewCMD(commandRetr, msg))
    resp, err := conn.readResponse(true)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// Close the TCP connection with the server, remember that you MUST send an quit command calling conn.Quit()
// before to close the session.
func (conn *Conn) Close() error {
    return conn.conn.Close()
}


func NewConn(addr string, config *tls.Config) *Conn {
	return &Conn{
		addr:   addr,
		tlsConfig: config,
	}
}
