package atem

import (
	"encoding/binary"
)

const (
	TallyByIndexCommand = "TlIn"
	TallyChannelConfig  = "_TlC"
)

func (c *AtemClient) parsePayload(packet []byte) error {
	offset := HeaderSize
	size := binary.BigEndian.Uint16(packet[offset : offset+2])

	for (offset + size) <= uint16(len(packet)) {
		if size == 0 {
			break
		}
		startOffset := offset + 2
		endOffset := startOffset + size - 2
		payload := packet[startOffset:endOffset]

		cmd := string(payload[2:6])
		offset += size
		size = binary.BigEndian.Uint16(packet[offset : offset+2])

		switch cmd {
		case TallyByIndexCommand:
			c.parseTallyByIndex(payload)
		case TallyChannelConfig:
			c.parseTallyChannelConfig(payload)
		}
	}

	return nil
}
