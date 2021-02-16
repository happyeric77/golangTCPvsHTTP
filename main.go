package main

import (
	"strings"
	"fmt"
	"log"
	"net"
	"bufio"
)

func main(){
	li, err := net.Listen("tcp", ":7160")
	if err != nil {
		log.Fatal(err)
	}
	defer li.Close()
	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	readRequest(conn)
}

func readRequest(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	lineCount := 0
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if lineCount == 0 {
			mux(conn, ln)
		}
		if ln == "" {
			break
		}
		lineCount ++
	}	
}

func mux(conn net.Conn, ln string) {
	m := strings.Fields(ln)[0] //method
	u := strings.Fields(ln)[1] //url
	fmt.Println("***Method: ", m)
	fmt.Println("***URL: ", u)

	//multiplexer
	if m == "GET" {
		switch u {
		case "/":
			writeResponse(conn, "Homepage")
		case "/about":
			writeResponse(conn, "AboutPage")
		}
	}
}

func writeResponse(conn net.Conn, page string) {
	body := fmt.Sprintf("<h1>%s</h1>", page)
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, body)
}