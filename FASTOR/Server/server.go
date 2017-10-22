package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var connections []Connections
var connDetail string
var mutex = &sync.Mutex{}
var flag string

//Connections to store details of connection
type Connections struct {
	conn  net.Conn
	ip    string
	port  string
	relay string
}

func delaySecond(n time.Duration) {
	time.Sleep(n * time.Second) // <------------ here
}

func printConnection(conn Connections) {
	fmt.Println("Connection's IP -> " + conn.ip)
	fmt.Println("port at connection is listening -> " + conn.port)
	var str string
	if conn.relay == "0" {
		str = "Entry Relay"
	} else if conn.relay == "1" {
		str = "Middle Relay"
	} else {
		str = "Exit Relay"
	}
	fmt.Println("Connection's Relay-> " + str)

}

func printConnections() {
	fmt.Println("printConnections() prints alive connections after every 15 seconds")
	for {
		fmt.Println("Alive Connections-> ", len(connections))
		for _, v := range connections {
			printConnection(v)
			fmt.Println()
		}
		fmt.Println()
		delaySecond(15)
	}
}

func handleConnection(conn net.Conn, check string) {
	conn.Write([]byte(check + ":" + flag))
	msg := make([]byte, 4096)
	n, _ := conn.Read(msg)

	temp := msg[0:n]

	ind := strings.Index(string(temp), ":")
	relay := temp[0:ind]
	port := temp[ind+1 : n]

	ip := conn.RemoteAddr().String()
	ind = strings.Index(ip, ":")
	IP := ip[0:ind]
	if IP == "[" {
		IP = "127.0.0.1"
	}

	if flag == "1" {
		fmt.Println("new Connected Device's IP-> " + IP)
	}
	connection := Connections{
		conn:  conn,
		relay: string(relay),
		port:  string(port),
		ip:    IP,
	}

	mutex.Lock()
	connections = append(connections, connection)
	mutex.Unlock()

	for {
		delaySecond(10)
		alive := make([]byte, 50)
		_, err := conn.Read(alive)
		if err != nil {
			mutex.Lock()
			for i, v := range connections {
				if v.conn == conn {
					connections = append(connections[:i], connections[i+1:]...)
					if flag == "1" {
						fmt.Println("Connection Closed->")
						printConnection(v)
					}
					break
				}
			}
			mutex.Unlock()
			conn.Close()
		}
	}
}

func sendClients() {
	if flag == "1" {
		fmt.Println("sendClients() list of alive clients to every client after every  10 seconds")
	}
	for {
		delaySecond(10)
		if len(connections) != 0 {
			connDetail = ""
			for i := 0; i < len(connections); i++ {
				if i != 0 {
					connDetail += "-"
				}
				connDetail = connDetail + connections[i].ip + ":" + connections[i].port + ":" + connections[i].relay
			}
			for i := 0; i < len(connections); i++ {
				connections[i].conn.Write([]byte(connDetail))
			}
		}
		if flag == "1" {
			fmt.Println("Sent Clients")
		}
	}
}

func main() {
	fmt.Println("Enter 0 to run on localhost, else enter 1")
	var check string
	fmt.Scan(&check)

	fmt.Println("Do you want to print all the data i.e Enter 1 to ON the flag bit else enter 0")
	fmt.Scan(&flag)

	ln, err := net.Listen("tcp", ":1805")

	if err != nil {
		log.Fatal(err)
	}

	go sendClients()
	if flag == "1" {
		go printConnections()
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, check)
	}
}
