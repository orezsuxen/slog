package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("starting mirror")
	msg := make([]byte, 5)

	for {

		_, err := os.Stdin.Read(msg)
		if err != nil {
			panic(err)
		}
		if msg[0] == 'q' {
			break
		} else {
			fmt.Println(msg)
		}

	}
	fmt.Println("ending mirror")
}
