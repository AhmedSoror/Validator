package main

import (
	"reflect"
	"testing"
)

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
	filepath := "./data/valid/operations.json"
	validateProgramTestCase(t, filepath, expectedResult)
}
func TestValidateProgramRec_FunctionCall(t *testing.T) {
	expectedResult := true
	filepath := "./data/valid/function_call.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_FunctionCallAsOperand(t *testing.T) {
	expectedResult := true
	filepath := "./data/valid/function_call_operand.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_OperationAsOperand(t *testing.T) {
	expectedResult := true
	filepath := "./data/valid/operation_operand.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UndeclaredFunction_fixed(t *testing.T) {
	expectedResult := true
	filepath := "./data/valid/undeclared_func_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UndeclaredVariable_fixed(t *testing.T) {
	expectedResult := true
	filepath := "./data/valid/undeclared_var_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UnassignedVariableInOperation_fixed(t *testing.T) {
	expectedResult := true
	filepath := "./data/valid/operation_with_unassigned_variable_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UnassignedVariableInFunctionCall_fixed(t *testing.T) {
	expectedResult := true
	filepath := "./data/valid/function_call_with_unassigned_variable_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_nestedFunctionCallUndeclaredFunction_fixed(t *testing.T) {
	expectedResult := true
	filepath := "./data/valid/nested_function_call_undeclared_function_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}
func TestValidateProgramRec_nestedFunctionCallUndeclaredVariable_fixed(t *testing.T) {
	expectedResult := true
	filepath := "./data/valid/nested_function_call_undeclared_variable_fixed.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

// --------------------------
// Test invalid programs
// --------------------------
func TestValidateProgramRec_UndeclaredFunction(t *testing.T) {
	expectedResult := false
	filepath := "./data/invalid/undeclared_func.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UndeclaredVariable(t *testing.T) {
	expectedResult := false
	filepath := "./data/invalid/undeclared_var.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UnassignedVariableInOperation(t *testing.T) {
	expectedResult := false
	filepath := "./data/invalid/operation_with_unassigned_variable.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_UnassignedVariableInFunctionCall(t *testing.T) {
	expectedResult := false
	filepath := "./data/invalid/function_call_with_unassigned_variable.json"
	validateProgramTestCase(t, filepath, expectedResult)
}

func TestValidateProgramRec_nestedFunctionCallUndeclaredFunction(t *testing.T) {
	expectedResult := false
	filepath := "./data/invalid/nested_function_call_undeclared_function.json"
	validateProgramTestCase(t, filepath, expectedResult)
}
func TestValidateProgramRec_nestedFunctionCallUndeclaredVariable(t *testing.T) {
	expectedResult := false
	filepath := "./data/invalid/nested_function_call_undeclared_variable.json"
	validateProgramTestCase(t, filepath, expectedResult)
}
