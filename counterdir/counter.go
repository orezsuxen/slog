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
		fmt.Println("counter got args! ", os.Args)
		n, err := strconv.Atoi(os.Args[1])
		if err == nil && n != 0 {
			// fmt.Println("counter args 1", os.Args[1])
			nums = n
		} else {
			fmt.Println("counter n,err", n, err)
			return
		}
	}

	if len(os.Args) > 2 {
		d, err := strconv.Atoi(os.Args[2])
		if err == nil && d != 0 {
			// fmt.Println("counter args 2", os.Args[2])
			delay = d
		} else {
			fmt.Println("counter d,err", d, err)
			return
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
