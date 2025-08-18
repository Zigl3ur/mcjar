package loader

import (
	"fmt"
	"sync"
	"time"
)

var loaderGlyphs = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
var (
	finish    chan struct{}
	mu        sync.Mutex
	loaderMsg string
	running   bool
)

func Start(message string) {
	mu.Lock()
	defer mu.Unlock()

	if running {
		return
	}

	finish = make(chan struct{})
	running = true
	loaderMsg = message
	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for {
			for _, g := range loaderGlyphs {
				select {
				case <-finish:
					return
				case <-ticker.C:
					fmt.Printf("\r%s %s", g, loaderMsg)
				}
			}
		}
	}()
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if !running || finish == nil {
		return
	}

	close(finish)
	finish = nil
	running = false
	fmt.Print("\r\033[K") // clear line
}

func UpdateMessage(message string) {
	mu.Lock()
	defer mu.Unlock()
	loaderMsg = message
}
