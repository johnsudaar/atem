package atem

import (
	"encoding/binary"
	"fmt"
)

const (
	TallyByIndexCommand = "TlIn"
	TallyChannelConfig  = "_TlC"
)

func (c *AtemClient) parsePayload(packet []byte) error {
	offset := HeaderSize
	size := binary.BigEndian.Uint16(packet[offset : offset+2])

	for (offset + size) < uint16(len(packet)) {
		startOffset := offset + 2
		endOffset := startOffset + size - 2
		payload := packet[startOffset:endOffset]

		cmd := string(payload[2:6])
		offset += size
		size = binary.BigEndian.Uint16(packet[offset : offset+2])

		fmt.Println(cmd)

		switch cmd {
		case TallyByIndexCommand:
			c.parseTallyByIndex(payload)
		case TallyChannelConfig:
			c.parseTallyChannelConfig(payload)
		}
	}

	return nil
}
