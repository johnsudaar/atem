package atem

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"

	"github.com/Scalingo/go-utils/logger"
	"github.com/pkg/errors"
)

func (c *AtemClient) send(buffer io.Reader) error {
	_, err := io.Copy(c.conn, buffer)
	if err != nil {
		return errors.Wrap(err, "fail to send buffer to ATEM")
	}
	return nil
}

func (c *AtemClient) connectToSwitcher() error {
	buff := bytes.NewBuffer(c.commandHeader(PacketTypeHello, 8, 0, false))

	binary.Write(buff, binary.BigEndian, uint32(0x01000000))
	binary.Write(buff, binary.BigEndian, uint32(0))

	err := c.send(buff)
	if err != nil {
		return errors.Wrap(err, "fail to send Hello Package")
	}

	err = c.listenSocket()
	if err != nil {
		return errors.Wrap(err, "fail to listen for first packet")
	}

	return nil
}

func (c *AtemClient) listenSocketLoop(ctx context.Context) {
	log := logger.Get(ctx)
	for {
		// Stopping mechanism
		c.stopMutex.Lock()
		stopping := c.stopping
		c.stopMutex.Unlock()
		if stopping != nil {
			stopping <- true
			return
		}

		err := c.listenSocket()
		if err != nil {
			log.WithError(err).Error("fail to listen on udp socket")
		}
	}
}

func (c *AtemClient) listenSocket() error {
	// TODO: Better error handling
	buffer := make([]byte, 1024)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return errors.Wrap(err, "fail to read socket")
	}
	packet := buffer[0:n]
	header := new(header)
	err = header.UnmarshalBinary(packet)
	if err != nil {
		return errors.Wrap(err, "fail to unmarshal packet")
	}

	c.currentUid = header.UID

	if (header.BitMask & PacketTypeHello) != 0 {
		ackBuffer := bytes.NewBuffer(c.commandHeader(PacketTypeAck, 0, 0x0, true))
		err := c.send(ackBuffer)
		if err != nil {
			return errors.Wrap(err, "fail to reply to hello packet")
		}
	} else if (header.BitMask & (PacketTypeAckRequest | PacketTypeResend)) != 0 {
		c.remotePacketCounter = header.PackageID
		ackBuffer := bytes.NewBuffer(c.commandHeader(PacketTypeAck, 0, header.PackageID, true))
		err := c.send(ackBuffer)
		if err != nil {
			return errors.Wrap(err, "fail to respond to ACK request")
		}
	}

	if uint16(len(packet)) > (HeaderSize+2) && (header.BitMask&(PacketTypeHello|PacketTypeResend)) == 0 {
		c.parsePayload(packet)
	}
	return nil
}
