package main

import (
	"context"
	"fmt"
	"time"

	"github.com/johnsudaar/atem"
)

type writer struct{}

func (writer) WriteTally(st atem.TallyStatuses) {
	fmt.Printf("%+v\n", st)
}

func main() {
	fmt.Println("Connect !")
	client, err := atem.New(context.Background(), "192.168.1.57:9910",
		atem.WithTallyWriter(writer{}),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected !")

	for {
		for i := 0; i <= 3; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println("SEND!!!")

			err = client.SetProgram(atem.MESource0, atem.VideoSource(i))
			if err != nil {
				panic(err)
			}
		}
	}
}
