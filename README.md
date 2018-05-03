# memcache


[![CircleCI](https://circleci.com/gh/A-Hilaly/memcache/tree/master.svg?style=svg&circle-token=8ae9aff37a33b81224f4bdb43b5d5621ac766f7b)](https://circleci.com/gh/A-Hilaly/memcache/tree/master)


 ## Example

 ```golang
package main

import (
	"fmt"
	"github.com/a-hilaly/memcache"
)

func main() {
    store := memcache.Default()
    store.Put("my-key", "myvalue")
    v, _ := store.Get("my-key")
    fmt.Println(v)

    store.Patch("my-key", struct{flag, size int}{50, 50})
    v, _ := store.Get("my-key")
    fmt.Println(v)

}
 ```
# Benchmarks
```shell
go test -v -bench .
```
see [benchmark file](BENCHMARKS)


 # TODO

 - ~~package~~
 - ~~cache audit~~
 - ~~tests~~
 - ~~benchmarks~~

