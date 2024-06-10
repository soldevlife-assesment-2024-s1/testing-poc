package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-json"
	"github.com/gorules/zen-go"
)

func main() {
	engine := zen.NewEngine(zen.EngineConfig{})

	graph, err := os.ReadFile("./ticket_price.json")
	if err != nil {
		panic(err)
	}

	decision, err := engine.CreateDecision(graph)
	if err != nil {
		panic(err)
	}

	response, err := decision.Evaluate(map[string]any{"level": "wood", "region": "asean"})
	if err != nil {
		panic(err)
	}

	// decode json response
	var result map[string]any

	byteRes, err := response.Result.MarshalJSON()
	if err != nil {
		panic(err)
	}

	response.Result.UnmarshalJSON(byteRes)

	err = json.Unmarshal(byteRes, &result)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
