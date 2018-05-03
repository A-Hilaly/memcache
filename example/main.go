package main

import (
	"fmt"

	"github.com/a-hilaly/memcache"
)

func main() {
	tain()
}

func tain() {
	store := memcache.Default()
	store.Put("my-key", "myvalue")
	v, _ := store.Get("my-key")
	fmt.Println(v)

	store.Patch("my-key", struct{ flag, size int }{50, 50})
	v, _ = store.Get("my-key")
	fmt.Println(v)

}
