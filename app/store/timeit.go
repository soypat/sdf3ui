package store

import (
	"fmt"
	"time"
)

// Just a temporary place for this TimeIt function
func TimeIt(name string) (stopAndPrint func()) {
	start := time.Now()
	return func() {
		fmt.Printf("%s:%s\n", name, time.Since(start))
	}
}
