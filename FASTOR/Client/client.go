package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var relay string
var port string
var webPort int
var connections []string
var dataChan *chan string

var check string
var flag string

func getMyIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	var str []string

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				str = append(str, (ipnet.IP.String()))
			}
		}
	}
	if check == "0" {
		return "127.0.0.1"
	}
	return str[len(str)-1]
}

func returnValue(ip string, passPort string) bool {
	if check == "0" {
		if passPort != port {
			return true
		}
		return false
	}

	if check == "1" {
		if ip != getMyIP() {
			return true
		}
		return false
	}
	return false
}

func delaySecond(n time.Duration) {
	time.Sleep(n * time.Second) // <------------ here
}

func handler(w http.ResponseWriter, r *http.Request) {
	// url := r.URL.Path[1:]
	// if len(url) != 0 {
	// 	http.Redirect(w, r, "/", http.StatusFound)
	// }
	// t, _ := template.ParseFiles("welcome.html")
	// t.Execute(w, nil)
	fmt.Fprint(w, "Welcome To FASTOR")
}

func handleFastor(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/FASTOR/"):]
	url = "http://" + url
	if flag == "1" {
		fmt.Println("Recieved request-> " + url)
	}
	for {
		rand.Seed(time.Now().UnixNano())
		numb := rand.Intn(len(connections) - 1)
		detail := strings.Split(connections[numb], ":")
		if detail[2] == "0" && returnValue(detail[0], detail[1]) == true {
			if flag == "1" {
				fmt.Println("Dialing to entry relay-> " + detail[0] + ":" + detail[1])
			}
			newConn, err := net.Dial("tcp", detail[0]+":"+detail[1])
			if err != nil {
				log.Fatal(err)
			}
			if flag == "1" {
				fmt.Println("Sending url link -> " + url)
			}
			url = url + "|||1|||" + getMyIP() + "|||" + port
			newConn.Write([]byte(url))
			newConn.Close()
			break
		}
	}
	str := <-*dataChan
	//	fmt.Println(string(([]byte(str))[0:0]))

	if flag == "1" {
		fmt.Println("Data Recieved")
	}
	fmt.Fprint(w, str)
	// t, _ := template.ParseFiles("request.html")
	// t.Execute(w, nil)
}

func handle(conn net.Conn) {
	fmt.Println("Welcome to FASTOR")

	fmt.Println("Enter relay:\nEnter 1 for Middle Relay or 2 for Exit Relay, and 0 for not participating")
	fmt.Scan(&relay)

	h := make([]byte, 4)
	n, _ := conn.Read(h)
	k := string(h[0:n])

	kk := strings.Split(k, ":")
	check = kk[0]
	flag = kk[1]

	conn.Write([]byte(relay + ":" + port))

	fmt.Println("You can access the webserver using port " + strconv.Itoa(webPort))

	fmt.Println("Server Port-> " + port)

	for {
		delaySecond(10)
		conn.Write([]byte("I am Alive"))
	}
}

func recvClients(conn net.Conn) {
	for {
		delaySecond(10)
		Clientsbytes := make([]byte, 16284)
		n, err := conn.Read(Clientsbytes)
		if err != nil {
			log.Fatal(err)
			return
		}

		if n != 0 {
			connections = nil
			C := string(Clientsbytes[0:n])
			connections = strings.Split(C, "-")
		}
		if flag == "1" {
			fmt.Println("Recieved Clients->")
			for i := 0; i < len(connections); i++ {
				fmt.Println(connections[i])
			}
			fmt.Println()
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	web := make([]byte, 16777216)
	//accept approx upto 24KB

	if flag == "1" {
		fmt.Println(conn.RemoteAddr().String() + " connected to me " + getMyIP() + ":" + port)
	}
	n, _ := conn.Read(web)
	webAddr := web[0:n]

	// this is either url or html content
	str := strings.Split(string(webAddr), "|||")
	// fmt.Println("Data Recieved i.e Actual user's ip and port -> ")
	// fmt.Println(str)
	url := str[0]
	condition := str[1]
	IP := str[2]
	portNo := str[3]

	if condition == "1" {
		if flag == "1" {
			fmt.Println("URL recieved -> " + url)
		}
	}
	// 0 to send back
	// 1 to send forward

	if condition == "1" {
		if relay == "0" {
			if flag == "1" {
				fmt.Println("This is Entry Relay")
			}
			n := len(connections) - 1
			// fmt.Println("no of connections = " + strconv.Itoa(n))
			for {
				rand.Seed(time.Now().UnixNano())
				numb := rand.Intn(n)
				// fmt.Print("Random number generated-> ")
				// fmt.Println(numb)
				detail := strings.Split(connections[numb], ":")
				if detail[2] == "1" && returnValue(detail[0], detail[1]) == true {
					if flag == "1" {
						fmt.Println("Dialing to middle relay-> " + detail[0] + ":" + detail[1])
					}
					newConn, err := net.Dial("tcp", detail[0]+":"+detail[1])
					if err != nil {
						log.Fatal(err)
					}
					newConn.Write(webAddr)
					newConn.Close()
					break
				}
			}
			return
		}

		if relay == "1" {
			if flag == "1" {
				fmt.Println("This is Middlie Relay")
			}
			n := len(connections) - 1
			// fmt.Println("no of connections = " + strconv.Itoa(n))
			for {
				rand.Seed(time.Now().UnixNano())
				numb := rand.Intn(n)
				// fmt.Print("Random number generated-> ")
				// fmt.Println(numb)
				detail := strings.Split(connections[numb], ":")
				if detail[2] == "2" && returnValue(detail[0], detail[1]) == true {
					if flag == "1" {
						fmt.Println("Dialing to exit relay-> " + detail[0] + ":" + detail[1])
					}
					newConn, err := net.Dial("tcp", detail[0]+":"+detail[1])
					if err != nil {
						log.Fatal(err)
					}
					newConn.Write(webAddr)
					newConn.Close()
					break
				}
			}
			return
		}

		if relay == "2" {
			if flag == "1" {
				fmt.Println("This is exit relay")
				fmt.Println("Fetching the URL")
			}
			res, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			robots, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			//			fmt.Println(robots)
			n := len(connections) - 1
			// fmt.Println("no of connections = " + strconv.Itoa(n))
			for {
				rand.Seed(time.Now().UnixNano())
				numb := rand.Intn(n)
				// fmt.Print("Random number generated-> ")
				// fmt.Println(numb)
				detail := strings.Split(connections[numb], ":")
				if detail[2] == "1" && returnValue(detail[0], detail[1]) == true {
					if flag == "1" {
						fmt.Println("Dialing to middle relay-> " + detail[0] + ":" + detail[1])
					}
					newConn, err := net.Dial("tcp", detail[0]+":"+detail[1])
					if err != nil {
						log.Fatal(err)
					}
					newConn.Write([]byte(string(robots) + "|||0|||" + IP + "|||" + portNo))
					newConn.Close()
					break
				}
			}
			return
		}
	}

	if condition == "0" {
		if getMyIP() == IP && portNo == port {
			filename := "request.html"
			ioutil.WriteFile(filename, []byte(url), 0600)
			*dataChan <- url
			return
		}

		if relay == "1" {
			if flag == "1" {
				fmt.Println("This is Middle Relay")
			}
			n := len(connections) - 1
			// fmt.Println("no of connections = " + strconv.Itoa(n))
			for {
				rand.Seed(time.Now().UnixNano())
				numb := rand.Intn(n)
				// fmt.Print("Random number generated-> ")
				// fmt.Println(numb)
				detail := strings.Split(connections[numb], ":")
				if detail[2] == "0" && (portNo != detail[1] || IP != detail[0]) {
					if flag == "1" {
						fmt.Println("Dialing to entry relay-> " + detail[0] + ":" + detail[1])
						fmt.Println("Sending fetched html")
					}
					newConn, err := net.Dial("tcp", detail[0]+":"+detail[1])
					if err != nil {
						log.Fatal(err)
					}
					newConn.Write(webAddr)
					newConn.Close()
					break
				}
			}
			return
		}

		if relay == "0" {
			if flag == "1" {
				fmt.Println("This is Entry Relay")
				fmt.Println("Dialing to actual user-> " + IP + ":" + portNo)
			}
			newConn, err := net.Dial("tcp", IP+":"+portNo)
			if err != nil {
				log.Fatal(err)
			}
			newConn.Write(webAddr)
			newConn.Close()
			return
		}
	}
}

func main() {

	d := make(chan string)
	dataChan = &d

	http.HandleFunc("/", handler)
	http.HandleFunc("/FASTOR/", handleFastor)
	rand.Seed(time.Now().UnixNano())
	webPort = rand.Intn(63000)
	webPort = webPort + 2001
	Wport := ":" + strconv.Itoa(webPort)
	go http.ListenAndServe(Wport, nil)

	conn, err := net.Dial("tcp", "127.0.0.1:1805")
	if err != nil {
		log.Fatal(err)
	}

	//	fmt.Println(getMyIP())

	go handle(conn)

	go recvClients(conn)

	numb := webPort
	for numb == webPort {
		rand.Seed(time.Now().UnixNano())
		numb = rand.Intn(63000)
		numb += 2001
	}

	port = strconv.Itoa(numb)
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}
