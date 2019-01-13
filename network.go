package atem

import (
	"bytes"
	"encoding/binary"
	"io"

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
	buff := bytes.NewBuffer(c.commandHeader(PacketTypeHello, 8, 0))

	binary.Write(buff, binary.BigEndian, uint32(0x01000000))
	binary.Write(buff, binary.BigEndian, uint32(0))

	err := c.send(buff)
	if err != nil {
		return errors.Wrap(err, "fail to send Hello Package")
	}

	return nil
}

func (c *AtemClient) listenSocket() {
	// TODO: Better error handling
	buffer := make([]byte, 1024)
	for {

		// Stopping mechanism
		c.stopMutex.Lock()
		stopping := c.stopping
		c.stopMutex.Unlock()
		if stopping != nil {
			stopping <- true
			return
		}

		n, err := c.conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		packet := buffer[0:n]
		header := new(header)
		err = header.UnmarshalBinary(packet)
		if err != nil {
			panic(err)
		}

		c.currentUid = header.UID

		if (header.BitMask & PacketTypeHello) != 0 {
			ackBuffer := bytes.NewBuffer(c.commandHeader(PacketTypeAck, 0, 0x0))
			err := c.send(ackBuffer)
			if err != nil {
				panic(err)
			}
		} else if (header.BitMask & PacketTypeAckRequest) != 0 {
			ackBuffer := bytes.NewBuffer(c.commandHeader(PacketTypeAck, 0, header.PackageID))
			err := c.send(ackBuffer)
			if err != nil {
				panic(err)
			}
		}

		if uint16(len(packet)) > (HeaderSize+2) && (header.BitMask&(PacketTypeHello|PacketTypeResend)) == 0 {
			c.parsePayload(packet)
		}
	}
}
