//go:generate stringer -type=VideoSource
package atem

import (
	"encoding/binary"
)

type VideoSource uint16

type TallyStatus struct {
	Source  VideoSource
	Program bool
	Preview bool
}

func (t TallyStatus) String() string {
	str := t.Source.String() + "\t"
	if t.Program {
		str += "PGM"
	} else if t.Preview {
		str += "PVW"
	} else {
		str += "OFF"
	}
	return str
}

type TallyStatuses []TallyStatus

func (t TallyStatuses) String() string {
	str := ""

	for _, s := range t {
		str += s.String() + "\n"
	}
	return str
}

type TallyWriter interface {
	WriteTally(TallyStatuses)
}

const (
	Input_1 VideoSource = iota
	Input_2
	Input_3
	Input_4
	Input_5
	Input_6
	Input_7
	Input_8
	Input_9
	Input_10
	Input_11
	Input_12
	Input_13
	Input_14
	Input_15
	Input_16
	Input_17
	Input_18
	Input_19
	Input_20
)

func (c *AtemClient) parseTallyByIndex(payload []byte) {
	size := binary.BigEndian.Uint16(payload[6:8])

	tallyStatuses := TallyStatuses{}
	for i := uint16(0); i < size; i++ {
		tally := payload[8+i]
		tallyStatuses = append(tallyStatuses, TallyStatus{
			Source:  VideoSource(i),
			Program: tally%2 != 0,
			Preview: (tally>>1)%2 != 0,
		})
	}
	if c.tallyWriter != nil {
		c.tallyWriter.WriteTally(tallyStatuses)
	}
}

func (c *AtemClient) parseTallyChannelConfig(payload []byte) {
	count := payload[4]
	c.configLock.Lock()
	defer c.configLock.Unlock()
	c.atemConfig.TallyChannels = int(count)
}
