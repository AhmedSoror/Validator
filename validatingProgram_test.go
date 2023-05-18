package main

import (
	"reflect"
	"testing"
)

func TestValidateProgramRec_OperationsOnly(t *testing.T) {
	// Define the input program
	program := ReadTestCaseFromJSON("./data/valid/operations.json")

	// Define the expected output
	expectedResult := true

	// Call the function
	result := ValidateProgramRec(program)

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
	result := ValidateProgramRec(program)

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
	result := ValidateProgramRec(program)

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
	result := ValidateProgramRec(program)

	// Compare the result with the expected output
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Unexpected result. Got %v, want %v", result, expectedResult)
	}
}
