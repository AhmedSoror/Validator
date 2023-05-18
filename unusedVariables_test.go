package main

import (
	"testing"
)

// helperFprTestCase a helper function to validate a test case
func helperFprTestCase(t *testing.T, filepath string, expectedResult []string) {
	// Define the input program
	program := ReadTestCaseFromJSON(filepath)

	// Call the function
	result := UnusedVariables(program)

	// Compare the result with the expected output
	// if !reflect.DeepEqual(result, expectedResult) {
	if !sameStringSlice(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}

// sameStringSlice function to compare 2 lists contains the same elements regardless of the order
func sameStringSlice(x []string, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}

// --------------------------
// Test valid programs
// --------------------------
func TestUnusedVariables_TestCase1(t *testing.T) {
	expectedResult := []string{"myFunction_result"}
	filepath := "./data/unused_variables/one_unused_variable.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_VariableUsedInFunctionCall(t *testing.T) {
	expectedResult := []string{}
	filepath := "./data/unused_variables/variable_used_in_function_call.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_VariablesFromFunctionParameters(t *testing.T) {
	expectedResult := []string{"myFunction_param_1", "myFunction_result"}
	filepath := "./data/unused_variables/unused_variables_from_function_parameters.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_NestedExampleWithFunctionCallAndPperations(t *testing.T) {
	expectedResult := []string{}
	filepath := "./data/unused_variables/nested_example_with_function_call_and_operations.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_UnusedVariablesInNestedBlocks(t *testing.T) {
	expectedResult := []string{"main_y", "main_z"}
	filepath := "./data/unused_variables/unused_variables_in_nested_blocks.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_UnusedVariableDeclaredInTwoFunctions(t *testing.T) {
	expectedResult := []string{"main_result"}
	filepath := "./data/unused_variables/same_var_declared_in_two_places.json"
	helperFprTestCase(t, filepath, expectedResult)
}
