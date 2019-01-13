package atem

import (
	"fmt"
	"net"
	"sync"

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

	tallyWriter TallyWriter
	stopMutex   sync.Mutex
	stopping    chan bool
}

type ClientOpt func(*AtemClient)

func WithTallyWriter(writer TallyWriter) ClientOpt {
	return func(c *AtemClient) {
		c.tallyWriter = writer
	}
}

func New(addr string, localPort string, opts ...ClientOpt) (*AtemClient, error) {
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

	//conn.SetReadDeadline(1 * time.Second)

	client := &AtemClient{
		packetCounter: 0,
		atemAddr:      addr,
		localAddr:     fmt.Sprintf("0.0.0.0:%s", localPort),
		conn:          conn,
		currentUid:    0x4242,
		stopping:      nil,
	}

	for _, opt := range opts {
		opt(client)
	}

	err = client.connectToSwitcher()
	if err != nil {
		client.conn.Close()
		return nil, errors.Wrap(err, "fail to send HELLO packet to switcher")
	}

	go client.listenSocket()
	return client, nil
}

func (c *AtemClient) Close() error {
	c.stopMutex.Lock()
	c.stopping = make(chan bool)
	c.stopMutex.Unlock()
	<-c.stopping
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
