package main

import (
	"fmt"
	"time"

	"github.com/a-hilaly/memcache"
)

func main() {
	// create new cache
	mc := memcache.New(10, 10, 2*time.Second)

	// create audit fn
	auditFn := memcache.NewAuditorFunc(time.Millisecond*100, time.Millisecond*500)

	// start auditing
	done, stop := auditFn(mc)

	// put some key & value
	mc.Put("key1", "value")

	time.Sleep(time.Second * 1)

	// put key2
	mc.Put("key2", "value2")

	time.Sleep(time.Second * 1)

	// send stop signal to the goroutine auditing the cache
	stop <- struct{}{}
	// wait for it to finish
	<-done

	// check keys values
	fmt.Println(mc.Get("key"))  // expired
	fmt.Println(mc.Get("key2")) // still existing
}
