package main

import (
	"fmt"
	"time"

	"github.com/johnsudaar/atem"
)

type writer struct{}

func (writer) WriteTally(st atem.TallyStatuses) {
	fmt.Printf("%+v\n", st)
}

func main() {
	client, err := atem.New("192.168.1.50:9910",
		atem.WithTallyWriter(writer{}),
	)
	if err != nil {
		panic(err)
	}

	for {
		for i := 1; i <= 4; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println("SEND!!!")

			err = client.SetProgram(atem.MESource0, atem.VideoSource(i))
			if err != nil {
				panic(err)
			}
		}
	}
}
