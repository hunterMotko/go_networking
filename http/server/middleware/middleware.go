package middleware

import (
	"log"
	"net/http"
	"time"
)

func Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      if r.Method == http.MethodTrace {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
      }
      w.Header().Set("X-Content-Type-Options", "nosniff")
      start := time.Now()
      next.ServeHTTP(w, r)
      log.Printf("NExt handler duration %v", time.Now().Sub(start))
    },
  )
}
