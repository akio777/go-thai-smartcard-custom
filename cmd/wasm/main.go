package main

import "github.com/somprasongd/gothaismartcard/smc"

// var c chan bool

// func init() {
// 	c = make(chan bool)
// }

func main() {
	c := make(chan struct{})
	println("Web Assembly is ready")
	smc.Connect(nil)
	<-c
}
