package main

import (
	"fmt"
	"net"
	"strings"
)

func listCommand(connection net.Conn, data []string) {
	if len(data) > 1 {
		switch data[0] {
		case "send":
			fmt.Printf("%s sending to %s\n", connection.LocalAddr(), data[1])

		case "connect":
			fmt.Println("connecting ")
		case "exit":
			fmt.Println(data)
			connection.Close()
		default:
			fmt.Println("default :")
		}
	}
}

func readData(connection net.Conn) []string {
	data := make([]byte, (1024 * 4))
	n, _ := connection.Read(data)
	str := strings.Split(string(data[0:n]), " ")
	var ret []string
	for _, val := range str {
		ret = append(ret, strings.TrimSpace(val))
	}
	return ret
}

type Pool struct {
	Clients []Client
}

type Client struct {
	IP         string
	Connection net.Conn
}

func (p *Pool) SendMessage(ip string, message string) {
	for _, val := range p.Clients {
		if val.IP == ip {
			val.Connection.Write([]byte(message))
		}
	}
}

func (p *Pool) AddClient(client Client) {
	p.Clients = append(p.Clients, client)
}

func main() {
	var pool Pool
	ln, _ := net.Listen("tcp", ":8888")
	fmt.Println("Server started...")

	for {
		conn, _ := ln.Accept()

		pool.AddClient(Client{
			IP:         conn.RemoteAddr().String(),
			Connection: conn,
		})
		// clients = append(clients, client{
		// 	IP:         conn.RemoteAddr().String(),
		// 	Connection: &conn,
		// })

		fmt.Printf("%+v\n", pool)

		fmt.Printf("Incoming connection from %s\n", conn.RemoteAddr())
		conn.Write([]byte("------------------------\n"))
		conn.Write([]byte(" Welcome to SeQL System\n"))
		conn.Write([]byte("------------------------\n"))
		conn.Write([]byte("Command :\n"))
		conn.Write([]byte("SEND <to> <message>\n"))
		go func() {
			defer conn.Close()
			for {
				data := readData(conn)
				// fmt.Println(len(data))
				listCommand(conn, data)
			}
		}()
	}

}
