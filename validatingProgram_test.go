package main

import (
	"reflect"
	"testing"
)

// --------------------------
// Test valid programs
// --------------------------
func TestValidateProgramRec_OperationsOnly(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/valid/operations.json")

	// Define the expected output
	expectedResult := true

	// Call the function
	result := ValidateProgramRec(program, false)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}
func TestValidateProgramRec_FunctionCall(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/valid/function_call.json")

	// Define the expected output
	expectedResult := true

	// Call the function
	result := ValidateProgramRec(program, false)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}

func TestValidateProgramRec_FunctionCallAsOperand(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/valid/function_call_operand.json")

	// Define the expected output
	expectedResult := true

	// Call the function
	result := ValidateProgramRec(program, false)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}

func TestValidateProgramRec_OperationAsOperand(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/valid/operation_operand.json")

	// Define the expected output
	expectedResult := true

	// Call the function
	result := ValidateProgramRec(program, false)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}

// --------------------------
// Test invalid programs
// --------------------------
func TestValidateProgramRec_UndeclaredFunction(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/invalid/undeclared_func.json")

	// Define the expected output
	expectedResult := false

	// Call the function
	result := ValidateProgramRec(program, false)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}

func TestValidateProgramRec_UndeclaredVariable(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/invalid/undeclared_var.json")

	// Define the expected output
	expectedResult := false

	// Call the function
	result := ValidateProgramRec(program, false)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}

func TestValidateProgramRec_UnassignedVariableInOperation(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/invalid/operation_with_unassigned_variable.json")

	// Define the expected output
	expectedResult := false

	// Call the function
	result := ValidateProgramRec(program, false)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}

func TestValidateProgramRec_UnassignedVariableInFunctionCall(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/invalid/function_call_with_unassigned_variable.json")

	// Define the expected output
	expectedResult := false

	// Call the function
	result := ValidateProgramRec(program, false)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}
