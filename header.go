package atem

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type header struct {
	BitMask     uint16
	PayloadSize uint16
	UID         uint16
	AckID       uint16
	PackageID   uint16
}

func (h *header) UnmarshalBinary(data []byte) error {
	if uint16(len(data)) < HeaderSize {
		return fmt.Errorf("Invalid header size: %v shout be >= %v", len(data), HeaderSize)
	}

	h.BitMask = uint16(data[0] >> 3)
	h.PayloadSize = binary.BigEndian.Uint16(data[0:2]) & 0x07FF
	h.UID = binary.BigEndian.Uint16(data[2:4])
	h.AckID = binary.BigEndian.Uint16(data[4:6])
	h.PackageID = binary.BigEndian.Uint16(data[10:12])

	return nil
}

func (h *header) MarshalBinary() (data []byte, err error) {
	val := uint16(h.BitMask << 11)
	val |= (h.PayloadSize + HeaderSize)

	buff := new(bytes.Buffer)

	binary.Write(buff, binary.BigEndian, val)
	binary.Write(buff, binary.BigEndian, h.UID)
	binary.Write(buff, binary.BigEndian, h.AckID)
	binary.Write(buff, binary.BigEndian, int32(0))
	binary.Write(buff, binary.BigEndian, h.PackageID)

	return buff.Bytes(), nil
}

func (c *AtemClient) commandHeader(bitmask, payloadSize, ackID uint16) []byte {
	packageID := uint16(0)

	if bitmask&(PacketTypeHello|PacketTypeAck) != 0 {
		c.packetCounter++
		packageID = c.packetCounter
	}

	h := &header{
		BitMask:     bitmask,
		PayloadSize: payloadSize,
		UID:         c.currentUid,
		AckID:       ackID,
		PackageID:   packageID,
	}

	value, err := h.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return value
}
