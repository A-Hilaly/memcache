# memcache
[![CircleCI](https://circleci.com/gh/A-Hilaly/memcache/tree/master.svg?style=svg&circle-token=8ae9aff37a33b81224f4bdb43b5d5621ac766f7b)](https://circleci.com/gh/A-Hilaly/memcache/tree/master) [![codecov](https://codecov.io/gh/A-Hilaly/memcache/branch/master/graph/badge.svg)](https://codecov.io/gh/A-Hilaly/memcache)

Zero dependencies memcache key|value store library. It supports multi threaded programs and offer an easy to use auditing properties.

### Features

- Zero dependencies 
- Thread safe cache & store objects
- Audit utilities to delete expired items

### Usage

```shell
go get -u github.com/a-hilaly/memcache
```

### Example

 ```go
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
	fmt.Println(mc.Get("key2")) // still exists
}

}
 ```
### Benchmarks

see [BENCHMARKS.md](BENCHMARKS.md)