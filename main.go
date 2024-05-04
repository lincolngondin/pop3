package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
)

type CMD struct {
	Keyword   string
	Arguments []string
}


func main() {
	conn, err := NewConn("pop.gmail.com:995")
	//conn.SetReadDeadline(time.Now().Add(time.Second*30))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

    c := make(chan bool)
	//buffer := make([]byte, 1024)
    var txt []byte
    read := bufio.NewScanner(conn)
	go func() {
		for {
            if read.Scan() {
                txt = read.Bytes()
            } else {
                break
            }
            
			//log.Printf("Read %d bytes!\n", n)
            //log.Println(string(txt))
            fmt.Println("S: ", string(txt))
			//log.Println(string(buffer[:n]))
		}
        c <- true
	}()

    var buf bufio.Reader = *bufio.NewReader(os.Stdin)
    go func(){
        for {
            cmd, err := buf.ReadString('\n')
            cmd = strings.TrimRight(cmd, "\n")
            if err != nil {
                log.Println(err)
            }
            //fmt.Println("Enviando comando ", []byte(cmd+"\r\n"))
            conn.Write([]byte(cmd+"\r\n"))
        }
    }()

    <-c
    log.Println("Saindo...!")

}

type Conn struct {
    conn *tls.Conn
}

func NewConn(addr string) (*tls.Conn, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
