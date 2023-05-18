package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func ReadTestCaseFromJSON(filePath string) Program {

	// Read the JSON file
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		os.Exit(1)
	}

	// Parse the JSON into the AST structure
	var program Program
	err = json.Unmarshal(jsonData, &program)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return Program{}
	}
	return program
}

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

func TestFindFunctionCalls_NoDependancies(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/functions/no_dependancies.json")

	// Define the expected output
	expectedResult := map[string]set{
		"myFunction": {},
	}

	// Call the function
	result := FindFunctionCalls(program)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}

func TestFindFunctionCalls_OneLevelDependancies(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/functions/one_level.json")

	// Define the expected output for the program sample
	expectedResult := map[string]set{
		"main":       {"myFunction": {}},
		"myFunction": {},
	}

	// Call the function and compare the result with the expected output
	result := FindFunctionCalls(program)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result for program sample. Got %v, want %v", result, expectedResult)
	}
}

func TestFindFunctionCalls_MutliLevelDependancies(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/functions/multilevel.json")

	// Define the expected output for the program sample
	expectedResult := map[string]set{
		"calculateProduct": {},
		"calculateSum":     {},
		"display":          {},
		"main":             {"calculateProduct": {}, "calculateSum": {}, "display": {}, "printNumber": {}},
		"printNumber":      {"display": {}},
	}

	// Call the function and compare the result with the expected output
	result := FindFunctionCalls(program)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result for program sample. Got %v, want %v", result, expectedResult)
	}
}
