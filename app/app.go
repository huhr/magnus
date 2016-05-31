package main

import (
	"fmt"
	"os"
	"bufio"
)

func main() {
    //file, _ := os.Create("/home/huhaoran/way/goofme/exec/app_done")
	reader := bufio.NewReader(os.Stdin)
	for true {
		msg, _, _ := reader.ReadLine()
		//file.Write(msg)
		if msg == nil {
			return
		}
		fmt.Printf("%s", msg)
	}
}
