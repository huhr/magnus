package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	/*file, _ := os.Create("/home/huhaoran/way/goofme/exec/app_done")
	for true {
		msg, _, _ := reader.ReadLine()
		//file.Write(msg)
		if msg == nil {
			return
		}
		fmt.Printf("%s", msg)
	}*/
	for true {
		reader := bufio.NewReader(os.Stdin)
		msg, _, err := reader.ReadLine()
		if err == nil {
			fmt.Printf("%s\n", msg)
		} else {
			fmt.Printf(err.Error())
			return
		}
	}
	print("hhee")
}
