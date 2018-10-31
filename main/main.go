package main

import (
	"time"

	"github.com/johnsudaar/atem"
)

func main() {
	_, err := atem.New("192.168.1.50:9910", "9123")
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(1 * time.Second)
	}

}
