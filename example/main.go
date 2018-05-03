package main

import (
	"fmt"
	"time"

	"github.com/a-hilaly/memcache"
)

func test() {
	c := memcache.Default()
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
	goto here
}

func main() {
	test()
}
