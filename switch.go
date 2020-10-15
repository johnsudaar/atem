package atem

import (
	"bytes"
	"encoding/binary"

	"github.com/juju/errgo"
)

type MESource byte

const (
	MESource0 = 0
	MESource1 = 1
)

func (c *AtemClient) SetProgram(me MESource, source VideoSource) error {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, source)
	binary.Write(buff, binary.BigEndian, uint8(me))
	binary.Write(buff, binary.BigEndian, uint8(0x80))

	cmd := c.commandBuffer([]byte("CPgI"), buff.Bytes())

	err := c.send(cmd)
	if err != nil {
		return errgo.Notef(err, "fail to send command")
	}

	return nil
}
