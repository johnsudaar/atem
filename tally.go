//go:generate stringer -type=VideoSource
package atem

import (
	"encoding/binary"
	"fmt"
)

type VideoSource uint16

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
	for i := uint16(0); i < size; i++ {
		tally := payload[8+i]

		if tally != 0 {

			fmt.Println(VideoSource(i).String())
		}

		if tally%2 != 0 {
			fmt.Println("PROGRAM")
		}

		if (tally>>1)%2 != 0 {
			fmt.Println("PREVIEW")
		}
	}
}
