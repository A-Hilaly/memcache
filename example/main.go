package main

import (
	"fmt"
	"time"

	"github.com/a-hilaly/memcache"
)

func test() {
	c := memcache.Debug()
	c.Start()
here:
	c.Put("a", "hey")
	c.Put("b", "hijack")
	c.Put("c", "kii")
	c.Immortalize("a")
	time.Sleep(500 * time.Millisecond)
	fmt.Println(c.ListKeys())
	time.Sleep(500 * time.Millisecond)
	fmt.Println(c.ListKeys())
	time.Sleep(500 * time.Millisecond)
	//c.Stop()
	//c.Put("hey", "xd")
	//fmt.Println(c.ListKeys())
	//time.Sleep(2 * time.Second)
	//fmt.Println(c.ListKeys())
	goto here
}

func main() {
	test()
}
