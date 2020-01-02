package atem

import (
	"bytes"
	"encoding/binary"
	"io"
)

func (c *AtemClient) commandBuffer(command, payload []byte) io.Reader {
	size := uint16(len(command) + len(payload) + 4)
	buff := bytes.NewBuffer(c.commandHeader(PacketTypeAckRequest, size, 0, false))
	binary.Write(buff, binary.BigEndian, size)
	binary.Write(buff, binary.BigEndian, uint16(0))

	buff.Write(command)
	buff.Write(payload)

	return buff
}
