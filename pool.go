package main

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type Pool struct {
	Clients []Client
}

type Client struct {
	UID        string
	Domain     string
	IP         string
	Connection net.Conn
}

func (p *Pool) ToString() string {
	str := "List Connection:\n"
	for _, val := range p.Clients {
		str = str + fmt.Sprintf("%s - %s (%s)\n", val.UID, val.IP, val.Domain)
	}
	return str
}

func (p *Pool) SendMessage(domain, message string) error {
	for _, val := range p.Clients {
		if val.Domain == domain {
			log.Printf("Sending %s to %s\n", message, domain)
			val.Connection.Write([]byte(message + "\n"))
			return nil
		}
	}
	return errors.New(domain + " not found")
}

func (p *Pool) AddClient(client Client) error {
	for _, val := range p.Clients {
		if val.Domain == client.Domain {
			return errors.New(client.Domain + "already joined in list")
		}
	}
	p.Clients = append(p.Clients, client)
	log.Printf("%s joined as %s\n", client.IP, client.Domain)
	return nil
}

func (p *Pool) DeleteClient(ip string) {
	for i, val := range p.Clients {
		if val.IP == ip {
			p.Clients = append(p.Clients[:i], p.Clients[i+1:]...)
		}
	}
}
