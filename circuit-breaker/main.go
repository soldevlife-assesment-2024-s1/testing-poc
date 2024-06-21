package main

import (
	"fmt"
	"time"

	circuit "github.com/rubyist/circuitbreaker"
)

func main() {
	// Initialize circuit breaker nginx
	cbNginx := InitCircuitBreaker("consecutive")
	// Initialize circuit breaker lighttpd
	cbLighttpd := InitCircuitBreaker("consecutive")

	// Initialize http client nginx with circuit breaker
	clientNginx := InitHttpClient(cbNginx)
	// Initialize http client lighttpd with circuit breaker
	clientLighttpd := InitHttpClient(cbLighttpd)

	for {
		if !cbNginx.Ready() {
			fmt.Println("Circuit breaker is open nginx")
			time.Sleep(2 * time.Second)
		}
		// Make a request to the server
		respNginx, err := clientNginx.Get("http://localhost:1111")
		if err != nil {
			// Handle error
			fmt.Println(err)
		}

		// Print the response
		fmt.Println("nginx", respNginx)
		time.Sleep(2 * time.Second)

		if !cbLighttpd.Ready() {
			fmt.Println("Circuit breaker is open Lighttpd")
			time.Sleep(2 * time.Second)
		}
		// Make a request to the server
		respLighthttpd, err := clientLighttpd.Get("http://localhost:2222")
		if err != nil {
			// Handle error
			fmt.Println(err)
		}

		// Print the response
		fmt.Println("Lighttpd", respLighthttpd)
		time.Sleep(2 * time.Second)
	}
}

// Init initializes the circuit breaker based on the configuration and breaker type: consecutive, error_rate, threshold
func InitCircuitBreaker(breakerType string) (cb *circuit.Breaker) {
	switch breakerType {
	case "consecutive":
		cb = circuit.NewConsecutiveBreaker(
			int64(3),
		)
	case "error_rate":
		cb = circuit.NewRateBreaker(
			95, 100,
		)
	default:
		cb = circuit.NewThresholdBreaker(
			int64(10),
		)
	}
	return cb
}

// InitHttpClient initializes the http client based on the configuration and circuit breaker that has been initialized before
func InitHttpClient(cb *circuit.Breaker) *circuit.HTTPClient {
	timeout := time.Duration(60) * time.Second
	client := circuit.NewHTTPClientWithBreaker(
		cb,
		timeout,
		nil,
	)
	return client
}
