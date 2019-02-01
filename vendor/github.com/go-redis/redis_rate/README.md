# Rate limiting for go-redis

[![Build Status](https://travis-ci.org/go-redis/redis_rate.svg?branch=master)](https://travis-ci.org/go-redis/redis_rate)
[![GoDoc](https://godoc.org/github.com/go-redis/redis_rate?status.svg)](https://godoc.org/github.com/go-redis/redis_rate)

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"

    "golang.org/x/time/rate"
    "github.com/go-redis/redis_rate"
    "github.com/go-redis/redis"
)

func handler(w http.ResponseWriter, req *http.Request, rateLimiter *redis_rate.Limiter) {
    userID := "user-12345"
    limit := int64(5)

    rate, delay, allowed := rateLimiter.AllowMinute(userID, limit)
    if !allowed {
        h := w.Header()
        h.Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
        h.Set("X-RateLimit-Remaining", strconv.FormatInt(limit-rate, 10))
        delaySec := int64(delay/time.Second)
        h.Set("X-RateLimit-Delay", strconv.FormatInt(delaySec, 10))
        http.Error(w, "API rate limit exceeded.", 429)
        return
    }

    fmt.Fprintf(w, "Hello world!\n")
    fmt.Fprint(w, "Rate limit remaining: ", strconv.FormatInt(limit-rate, 10))
}

func statusHandler(w http.ResponseWriter, req *http.Request, rateLimiter *redis_rate.Limiter) {
    userID := "user-12345"
    limit := int64(5)

    // With n=0 we just retrieve the current limit.
    rate, delay, allowed := rateLimiter.AllowN(userID, limit, time.Minute, 0)
    fmt.Fprintf(w, "Current rate: %d", rate)
    fmt.Fprintf(w, "Delay: %s", delay)
    fmt.Fprintf(w, "Allowed: %v", allowed)
}

func main() {
    ring := redis.NewRing(&redis.RingOptions{
        Addrs: map[string]string{
            "server1": "localhost:6379",
        },
    })
    limiter := redis_rate.NewLimiter(ring)
    // Optional.
    limiter.Fallback = rate.NewLimiter(rate.Every(time.Second), 100)
    
    http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        handler(w, req, limiter)
    })
    
    http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
        statusHandler(w, req, limiter)
    })
    
    http.HandleFunc("/favicon.ico", http.NotFound)
    log.Println("listening on localhost:8888...")
    log.Println(http.ListenAndServe("localhost:8888", nil))
}
```
