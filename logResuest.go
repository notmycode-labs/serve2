package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip := strings.Split(r.RemoteAddr, ":")[0]
		method := r.Method
		url := r.URL.Path
		userAgent := r.UserAgent()

		log.Printf("[%s] %s %s from %s (User-Agent: %s)\n", time.Now().Format(time.RFC3339), method, url, ip, userAgent)

		next.ServeHTTP(w, r)

		elapsed := time.Since(start)
		log.Printf("[%s] %s %s from %s completed in %s\n", time.Now().Format(time.RFC3339), method, url, ip, elapsed)
	}
}
