package main

import (
	"net/http"

	"github.com/joho/godotenv"
	rateLimiter "github.com/joseasousa/rate_limiter/rate_limiter"
)

func main() {
	godotenv.Load(".env")

	rateLimiter := rateLimiter.NewRateLimiter()

	r := http.NewServeMux()

	handle := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	handleFunc := rateLimiter(http.HandlerFunc(handle))

	r.Handle("/",
		handleFunc,
	)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
