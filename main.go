package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

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

func main() {
	var pool Pool
	ln, _ := net.Listen("tcp", ":8888")
	log.Println("Server started...")

	for {
		conn, _ := ln.Accept()

		log.Printf("%s connected\n", conn.RemoteAddr())
		conn.Write([]byte("------------------------\n"))
		conn.Write([]byte(" Welcome to SeQL System\n"))
		conn.Write([]byte("------------------------\n"))
		conn.Write([]byte("Command :\n"))
		conn.Write([]byte("JOIN <domain>         join to target domain\n"))
		conn.Write([]byte("SEND <to> <message>   sending message to target domain\n"))
		conn.Write([]byte("EXIT                  exit from network\n"))
		go func() {
			defer conn.Close()
			for {
				data := readData(conn)
				if data != nil && data[0] != "" {
					switch strings.ToLower(data[0]) {
					case "list":
						conn.Write([]byte(pool.ToString()))
					case "join":
						err := pool.AddClient(Client{
							UID:        genUID(32),
							Domain:     data[1],
							IP:         conn.RemoteAddr().String(),
							Connection: conn,
						})
						if err != nil {
							log.Println(err)
							conn.Write([]byte(data[1] + " already Joined\n"))
						} else {
							conn.Write([]byte(data[1] + " joined\n"))
						}
					case "send":
						err := pool.SendMessage(data[1], data[2])
						if err != nil {
							conn.Write([]byte(data[1] + " not found\n"))
						}
					case "exit":
						log.Printf("%s disconnected\n", conn.RemoteAddr())
						pool.DeleteClient(conn.RemoteAddr().String())
						conn.Close()
					default:
						fmt.Println("miaw")
					}
				}
			}
		}()
	}

}
