package main

import (
	"fmt"
	"time"

	circuit "github.com/rubyist/circuitbreaker"
)

func main() {
	// Initialize circuit breaker
	cb := InitCircuitBreaker("consecutive")

	// Initialize http client with circuit breaker
	client := InitHttpClient(cb)

	for {
		if !cb.Ready() {
			fmt.Println("Circuit breaker is open")
			time.Sleep(2 * time.Second)
			continue
		}
		// Make a request to the server
		resp, err := client.Get("http://localhost:1111")
		if err != nil {
			// Handle error
			fmt.Println(err)
		}

		// Print the response
		fmt.Println(resp)
		time.Sleep(2 * time.Second)
	}
}

// Init initializes the circuit breaker based on the configuration and breaker type: consecutive, error_rate, threshold
func InitCircuitBreaker(breakerType string) (cb *circuit.Breaker) {
	switch breakerType {
	case "consecutive":
		cb = circuit.NewConsecutiveBreaker(
			int64(10),
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
