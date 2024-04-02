package syncx

import (
	"log"
	"runtime"
)

func GoSafe(fn func()) {
	go runSafe(fn)
}

func runSafe(fn func()) {
	defer func() {
		if p := recover(); p != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			log.Printf("==> %s\n\n", string(buf[:n]))
		}
	}()
	fn()
}
