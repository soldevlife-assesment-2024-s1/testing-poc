package main

import "github.com/roylee0704/gron"

func main() {
	c := gron.New()
	c.AddFunc(, func() {
		println("Happy new year!")
	})
	c.Start()
}
