//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"go-reading-log-api-next/test/testdata"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_expected.go <project_id>")
		os.Exit(1)
	}

	projectID := os.Args[1]

	// Generate expected values based on project ID
	switch projectID {
	case "450":
		generateProject450()
	default:
		fmt.Printf("Unknown project ID: %s\n", projectID)
		os.Exit(1)
	}
}

func generateProject450() {
	fmt.Println("Generating expected values for Project 450...")

	// Read the source JSON files
	goData, err := os.ReadFile("../data/project-450-go.json")
	if err != nil {
		fmt.Printf("Error reading Go JSON: %v\n", err)
		os.Exit(1)
	}

	railsData, err := os.ReadFile("../data/project-450-rails.json")
	if err != nil {
		fmt.Printf("Error reading Rails JSON: %v\n", err)
		os.Exit(1)
	}

	// Parse and process the data
	// For now, we'll generate the expected values using the CalculateExpectedValues function

	// Create a simple project response for testing
	tzLocation := time.FixedZone("BRT", -3*60*60)
	ctx := context.WithValue(context.Background(), "timezone", tzLocation)

	expected := testdata.GenerateExpectedValues(
		450,
		"História da Igreja VIII.1",
		691,
		691,
		"2026-02-19T00:00:00Z",
	)

	// Marshal to JSON for verification
	jsonData, err := json.MarshalIndent(expected, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))

	// Write to a file for reference
	err = os.WriteFile("../data/project-450-expected.json", jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing expected JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nExpected values generated successfully!")
	fmt.Println("Output written to: ../data/project-450-expected.json")
}
