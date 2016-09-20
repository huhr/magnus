package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		fmt.Println(sig)
	}()
	for true {
		msg, _, err := reader.ReadLine()
		if err == nil {
			fmt.Printf("%s\n", msg)
		} else {
			fmt.Printf("%s\n", msg)
			fmt.Println(err.Error())
			break
		}
	}
	fmt.Println("Process App Start ShutDown")
	time.Sleep(5 * time.Second)
	fmt.Println("Process App ShutDown")
}
