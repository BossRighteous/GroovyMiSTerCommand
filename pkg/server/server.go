package server

import (
	"fmt"
	"net"

	"github.com/BossRighteous/GroovyMiSTerCommand/pkg/command"
)

type UdpClient struct {
	host string
	conn net.PacketConn
	addr *net.UDPAddr
}

func (client *UdpClient) SendBeacon() {
	fmt.Println("Sending Beacon")
	_, err := client.conn.WriteTo([]byte{0}, client.addr)
	if err != nil {
		fmt.Println("UDP BEACON ERROR", err)
	}
}

func (client *UdpClient) Listen(cmdChan chan command.GroovyMiSTerCommand) {
	go func() {
		for {
			buf := make([]byte, 1024)
			rlen, _, err := client.conn.ReadFrom(buf)
			if err != nil {
				fmt.Println("UDP READ ERROR", err)
				continue
			}
			if rlen == 0 {
				fmt.Println("UDP READ EMPTY")
				continue
			}
			cmd, err := command.ParseGMC(buf[:rlen])
			if err != nil {
				fmt.Println("Command Parse Error")
				continue
			}
			cmdChan <- cmd
		}
	}()
}

func StartUdpClient(host string, cmdChan chan command.GroovyMiSTerCommand) *UdpClient {
	client := &UdpClient{}
	client.host = host
	conn, err := net.ListenPacket("udp4", ":32105")
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:32105", host))
	if err != nil {
		panic(err)
	}

	client.conn = conn
	client.addr = addr

	client.Listen(cmdChan)
	return client
}
