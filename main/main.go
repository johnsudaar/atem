package main

import (
	"fmt"
	"time"

	"github.com/johnsudaar/atem"
)

func main() {
	_, err := atem.New("192.168.1.50:9910")
	if err != nil {
		panic(err)
	}
	fmt.Println("AAA")

	for {
		time.Sleep(1 * time.Second)
	}

}
