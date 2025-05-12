package helper

import (
	"context"
	"log"
	"sync"
	"time"
)

var (
	AttemptMu     sync.Mutex
	LoginAttempt  = make(map[string]int)
	TokenAttempt  = make(map[string]int)
	CleanupTicker = time.NewTicker(1 * time.Minute)
)

func CleanUp(ctx context.Context) {
	for {
		select {
		case <-CleanupTicker.C:
			AttemptMu.Lock()
			LoginAttempt = make(map[string]int)
			TokenAttempt = make(map[string]int)
			AttemptMu.Unlock()
		case <-ctx.Done():
			CleanupTicker.Stop()
			log.Println("CleanUp: received shutdown signal, stopping...")
			return
		}
	}
}
