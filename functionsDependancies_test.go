package main

import (
	"reflect"
	"testing"
)

func TestRollOutDependencies(t *testing.T) {
	// Define the input dependencies
	dependencies := map[string][]string{
		"A": {"B", "C"},
		"B": {"D"},
		"C": {"E"},
		"D": {"C"},
		"E": {},
	}

	// Define the expected rolledOut dependencies
	expected := map[string]set{
		"A": {"B": {}, "C": {}, "D": {}, "E": {}},
		"B": {"C": {}, "D": {}, "E": {}},
		"C": {"E": {}},
		"D": {"C": {}, "E": {}},
		"E": {},
	}

	// Call the RollOutDependencies function
	result := RollOutDependencies(dependencies)

	// Compare the actual result with the expected result
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RollOutDependencies() = %v, want %v", result, expected)
	}
}
