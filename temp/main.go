package main

import (
	"fmt"
	"time"
)


func foo() {
	panic("errorroro")
}

func main() {
	go func() {
		foo()
	}()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover from panic")
		}
	}()
	time.Sleep(time.Second * 5)
}