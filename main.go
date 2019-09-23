package main

import (
	"fmt"
	"uspscan/uspscan"
)

func main() {
	err := uspscan.Scan()
	if err != nil {
		fmt.Println(err)
	}
}
