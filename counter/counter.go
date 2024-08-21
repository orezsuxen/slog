package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	var nums int = 10
	var delay int = 1000

	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err == nil && n != 0 {
			nums = n
		}
	}

	if len(os.Args) > 2 {
		d, err := strconv.Atoi(os.Args[2])
		if err == nil && d != 0 {
			delay = d
		}
	}

	// fmt.Println(nums)
	// fmt.Println(delay)

	for i := 0; i < nums; i++ {
		start := time.Now()
		for time.Since(start) < (time.Duration(delay) * time.Millisecond) {
			time.Sleep(1 * time.Millisecond)
		}
		fmt.Println(i)

	}
}
