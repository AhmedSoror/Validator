package validator

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

// validateProgramTestCase a helper function to validate a test case
func validateProgramTestCase(t *testing.T, filepath string, expectedResult bool) {
	// Define the input program
	program := ReadTestCaseFromJSON(filepath)

	// Call the function
	result := ValidateProgramRec(program, false)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}

// --------------------------
// Test valid programs
// --------------------------
func TestValidateProgramRec_OperationsOnly(t *testing.T) {
	expectedResult := true
	filepath := "../data/valid/operations.json"
	validateProgramTestCase(t, filepath, expectedResult)
}
func TestValidateProgramRec_FunctionCall(t *testing.T) {
	expectedResult := true
	filepath := "../data/valid/function_call.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_FunctionCallAsOperand(t *testing.T) {
	expectedResult := true
	filepath := "../data/valid/function_call_operand.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_OperationAsOperand(t *testing.T) {
	expectedResult := true
	filepath := "../data/valid/operation_operand.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UndeclaredFunction_fixed(t *testing.T) {
	expectedResult := true
	filepath := "../data/valid/undeclared_func_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UndeclaredVariable_fixed(t *testing.T) {
	expectedResult := true
	filepath := "../data/valid/undeclared_var_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UnassignedVariableInOperation_fixed(t *testing.T) {
	expectedResult := true
	filepath := "../data/valid/operation_with_unassigned_variable_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UnassignedVariableInFunctionCall_fixed(t *testing.T) {
	expectedResult := true
	filepath := "../data/valid/function_call_with_unassigned_variable_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_nestedFunctionCallUndeclaredFunction_fixed(t *testing.T) {
	expectedResult := true
	filepath := "../data/valid/nested_function_call_undeclared_function_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}
func TestValidateProgramRec_nestedFunctionCallUndeclaredVariable_fixed(t *testing.T) {
	expectedResult := true
	filepath := "../data/valid/nested_function_call_undeclared_variable_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

// --------------------------
// Test invalid programs
// --------------------------
func TestValidateProgramRec_UndeclaredFunction(t *testing.T) {
	expectedResult := false
	filepath := "../data/invalid/undeclared_func.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UndeclaredVariable(t *testing.T) {
	expectedResult := false
	filepath := "../data/invalid/undeclared_var.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UnassignedVariableInOperation(t *testing.T) {
	expectedResult := false
	filepath := "../data/invalid/operation_with_unassigned_variable.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UnassignedVariableInFunctionCall(t *testing.T) {
	expectedResult := false
	filepath := "../data/invalid/function_call_with_unassigned_variable.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_nestedFunctionCallUndeclaredFunction(t *testing.T) {
	expectedResult := false
	filepath := "../data/invalid/nested_function_call_undeclared_function.json"
	validateProgramTestCase(t, filepath, expectedResult)
}
func TestValidateProgramRec_nestedFunctionCallUndeclaredVariable(t *testing.T) {
	expectedResult := false
	filepath := "../data/invalid/nested_function_call_undeclared_variable.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_functionCallWithInvalidArity(t *testing.T) {
	expectedResult := false
	filepath := "../data/invalid/function_call_with_wrong_arity.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

// -------------------------------------
// Test functions dependancies
// -------------------------------------

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
	program := ReadTestCaseFromJSON("../data/functions/no_dependancies.json")

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
	program := ReadTestCaseFromJSON("../data/functions/one_level.json")

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
	program := ReadTestCaseFromJSON("../data/functions/multilevel.json")

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

// ---------------------------------------

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
// Test unused variables
// --------------------------
func TestUnusedVariables_TestCase1(t *testing.T) {
	expectedResult := []string{"myFunction_result"}
	filepath := "../data/unused_variables/one_unused_variable.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_VariableUsedInFunctionCall(t *testing.T) {
	expectedResult := []string{}
	filepath := "../data/unused_variables/variable_used_in_function_call.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_VariablesFromFunctionParameters(t *testing.T) {
	expectedResult := []string{"myFunction_param_1", "myFunction_result"}
	filepath := "../data/unused_variables/unused_variables_from_function_parameters.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_NestedExampleWithFunctionCallAndPperations(t *testing.T) {
	expectedResult := []string{}
	filepath := "../data/unused_variables/nested_example_with_function_call_and_operations.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_UnusedVariablesInNestedBlocks(t *testing.T) {
	expectedResult := []string{"main_y", "main_z"}
	filepath := "../data/unused_variables/unused_variables_in_nested_blocks.json"
	helperFprTestCase(t, filepath, expectedResult)
}

func TestUnusedVariables_UnusedVariableDeclaredInTwoFunctions(t *testing.T) {
	expectedResult := []string{"main_result"}
	filepath := "../data/unused_variables/same_var_declared_in_two_places.json"
	helperFprTestCase(t, filepath, expectedResult)
}
