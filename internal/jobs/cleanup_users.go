package jobs

import (
	"log"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/service"
)

func StartCleanupUserWorker(authService *service.AuthService) {
	ticker := time.NewTicker(15 * time.Minute)

	for range ticker.C {
		err := authService.DeleteExpiredUnverifiedUsers()
		if err != nil {
			log.Println("cleanup error:", err)
		}
	}
}
