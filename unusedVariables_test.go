package main

import (
	"reflect"
	"testing"
)

// helperFprTestCase a helper function to validate a test case
func helperFprTestCase(t *testing.T, filepath string, expectedResult []string) {
	// Define the input program
	program := ReadTestCaseFromJSON(filepath)

	// Call the function
	result := UnusedVariables(program)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}

// --------------------------
// Test valid programs
// --------------------------
func TestUnusedVariables_TestCase1(t *testing.T) {
	expectedResult := []string{"result"}
	filepath := "./data/unused_variables/one_unused_variable.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_VariableUsedInFunctionCall(t *testing.T) {
	expectedResult := []string{}
	filepath := "./data/unused_variables/variable_used_in_function_call.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_VariablesFromFunctionParameters(t *testing.T) {
	expectedResult := []string{"param_1", "result"}
	filepath := "./data/unused_variables/unused_variables_from_function_parameters.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_NestedExampleWithFunctionCallAndPperations(t *testing.T) {
	expectedResult := []string{}
	filepath := "./data/unused_variables/nested_example_with_function_call_and_operations.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_UnusedVariablesInNestedBlocks(t *testing.T) {
	expectedResult := []string{"y", "z"}
	filepath := "./data/unused_variables/unused_variables_in_nested_blocks.json"
	helperFprTestCase(t, filepath, expectedResult)
}
