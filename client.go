package atem

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

const (
	HeaderSize uint16 = 0x0c

	PacketTypeNoCommand  = 0x00
	PacketTypeAckRequest = 0x01
	PacketTypeHello      = 0x02
	PacketTypeResend     = 0x04
	PacketTypeUndefined  = 0x08
	PacketTypeAck        = 0x10
)

type AtemClient struct {
	packetCounter uint16
	conn          *net.UDPConn
	atemAddr      string
	localAddr     string
	currentUid    uint16
}

func New(addr string, localPort string) (*AtemClient, error) {
	atemAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%s", localPort))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", localAddr, atemAddr)
	if err != nil {
		return nil, err
	}
	client := &AtemClient{
		packetCounter: 0,
		atemAddr:      addr,
		localAddr:     fmt.Sprintf("0.0.0.0:%s", localPort),
		conn:          conn,
		currentUid:    0x4242,
	}

	err = client.connectToSwitcher()
	if err != nil {
		return nil, errors.Wrap(err, "fail to send HELLO packet to switcher")
	}

	go client.listenSocket()
	return client, nil
}
