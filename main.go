/*
*
POP3 client

example:
conn, err := pop3.NewConn("pop.gmail.com:995")

	if err != nil {
	    log.Fatal(err)
	}

response := conn.List()
conn.Close()
*/
package main

import (
	"log"
)

func main() {
	conn := NewConn("pop.gmail.com:995", nil)
	resp, err := conn.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for _, r := range resp.info {
		log.Print(string(r))
	}

	resp, err = conn.User("")
	for _, r := range resp.info {
		log.Print(string(r))
	}
	resp, err = conn.Pass("")
	for _, r := range resp.info {
		log.Print(string(r))
	}

	resp, err = conn.Stat()
	for _, r := range resp.info {
		log.Print(string(r))
	}

	resp, err = conn.List("")
	for _, r := range resp.info {
		log.Print(string(r))
	}

	log.Println("Listing 1")
	resp, err = conn.List("1")
	for _, r := range resp.info {
		log.Print(string(r))
	}

	quitResponse, quitErr := conn.Quit()
	if quitErr != nil {
		log.Fatal(quitErr)
	}
	for _, r := range quitResponse.info {
		log.Print(string(r))
	}
	conn.Close()
}
