package display

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

const (
	cmdHeaderClose     byte = 1
	cmdHeaderInit      byte = 2
	cmdHeaderSwitchRes byte = 3
	cmdHeaderBlit      byte = 6
)

type UdpDisplayClient struct {
	host         string
	conn         net.PacketConn
	addr         *net.UDPAddr
	frame        uint32
	mtuBlockSize int32
}

func (client *UdpDisplayClient) SendPacket(buffer []byte) {
	//fmt.Println("Sending Packet length", len(buffer))
	_, err := client.conn.WriteTo(buffer, client.addr)
	if err != nil {
		fmt.Println("Connection Error (Close):", err)
		client.conn.Close()
	}
}

func (client *UdpDisplayClient) SendMTU(buffer []byte) {
	bytesToSend := int32(len(buffer))
	chunkMaxSize := int32(client.mtuBlockSize)
	var chunkSize int32 = 0
	var offset int32 = 0
	for bytesToSend > 0 {
		chunkSize = chunkMaxSize
		if bytesToSend <= chunkMaxSize {
			chunkSize = bytesToSend
		}
		bytesToSend = bytesToSend - chunkSize
		client.SendPacket(buffer[offset : offset+chunkSize])
		offset += chunkSize
	}
}

func (client *UdpDisplayClient) CmdClose() {
	buffer := make([]byte, 1)
	buffer[0] = cmdHeaderClose
	client.SendPacket(buffer)
	client.Close()
}

func (client *UdpDisplayClient) CmdInit() {
	buffer := make([]byte, 5)
	buffer[0] = cmdHeaderInit
	buffer[1] = 0 // lz4 compression flag
	buffer[2] = 0 // sound rate
	buffer[3] = 0 // sound channel
	buffer[4] = 0 // rgb mode
	client.SendPacket(buffer)
}

func (client *UdpDisplayClient) CmdSwitchres() {
	buffer := make([]byte, 26)
	buffer[0] = cmdHeaderSwitchRes
	binary.LittleEndian.PutUint64(buffer[1:9], math.Float64bits(6.7))
	binary.LittleEndian.PutUint16(buffer[9:11], 320)
	binary.LittleEndian.PutUint16(buffer[11:13], 336)
	binary.LittleEndian.PutUint16(buffer[13:15], 367)
	binary.LittleEndian.PutUint16(buffer[15:17], 426)
	binary.LittleEndian.PutUint16(buffer[17:19], 240)
	binary.LittleEndian.PutUint16(buffer[19:21], 244)
	binary.LittleEndian.PutUint16(buffer[21:23], 247)
	binary.LittleEndian.PutUint16(buffer[23:25], 262)
	buffer[25] = 0
	client.SendPacket(buffer)
}

func (client *UdpDisplayClient) CmdBlit(frameBuffer []byte) {
	client.frame++
	buffer := make([]byte, 7)
	buffer[0] = cmdHeaderBlit
	binary.LittleEndian.PutUint32(buffer[1:5], client.frame)
	binary.LittleEndian.PutUint16(buffer[5:7], 0)
	client.SendPacket(buffer)
	client.SendMTU(frameBuffer)
}

func (client *UdpDisplayClient) Close() {
	client.conn.Close()
}

func (client *UdpDisplayClient) Open() {
	conn, err := net.ListenPacket("udp4", ":32100")
	if err != nil {
		fmt.Println("Connection Error (Open):", err)
	}
	client.conn = conn
}

func NewUdpClient(host string) *UdpDisplayClient {
	var client UdpDisplayClient
	client.host = host
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:32100", host))
	if err != nil {
		panic(err)
	}

	//client.conn = conn
	client.addr = addr
	client.mtuBlockSize = 1472
	return &client
}
