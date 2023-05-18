/*
----------------------
Language description
----------------------
This tool deals with a simplified representation for an imperative programming language that allows static program analysis to be applied.
The  representation is defined as follows:
	● A program contains one or more function declarations.
	● A function declaration contains a block specifying the function body.
	● There are four kinds of statements: blocks, variable declarations, operations and function calls:
			○ A block contains zero or more statements.
			○ A variable declaration declares a variable that may be used in an operation.
			○ An operation has zero or more variable uses.
			○ A function call references a function declaration.

----------------------
Tool features implementd in this script:
----------------------
The tool expects a file specifying a program for the previously defined representation, and offers three different operations:
	1. Verify that a program is valid. The conditions for a program to be valid are:
		a. A function call must call a function that is declared in the same file.
		b. A variable can only be used in operations if it has been declared in a previous statement of the same block, or in case it has been declared in
		one of the previous statements of a surrounding block.

		extra added conditions:
		----------------
			- Variable must be assigned before being used
			- Function arguments are valid operands:
				- Numerical value
				- declared and assigned variable
				- valid function call

	2. List variables that are declared but not used.
	3. For each function, list which other functions they depend on. A function depends on another function if it is directly or indirectly called.
		For example, if function A calls Function B and Function B calls function C, then function A depends on both functions B and C.


----------------------
Assumptions:
----------------------
	- variables are dynamically typed
	- operands in an operation or arguments in a function call can be:
		- numerical
		- variable
		- function call
		- operation
	- variables can't be declared twice
	- function's parameters are considered as declaration for variable and they are already assigned
	- in Assignment operation, the assigned variable is the first variable in the operations list
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// -----------------------------------------
// Define the structure for the AST representation
// -----------------------------------------

// Program represents the top-level structure of the program.
type Program struct {
	// A program contains one or more function declarations.
	Functions []Function `json:"functions"` // List of function declarations
}

// Function represents a function declaration.
type Function struct {
	/*
		A function declaration contains a block specifying the function body.
		Extra added info:
			- function identifier "name"
			- function Parameters
	*/
	Name       string   `json:"name"`       // Name of the function
	Parameters []string `json:"parameters"` // List of function arguments
	Body       Block    `json:"body"`       // Function body
}

// Block represents a block of statements.
type Block struct {
	// A block contains zero or more statements.
	Statements []Statement `json:"statements"` // List of statements in the block
}

// Statement represents an individual statement.
type Statement struct {
	Type           string      `json:"type"`                      // Type of statement (block, variable_declaration, operation, function_call)
	Value          string      `json:"value,omitempty"`           // declared variable
	Variable       string      `json:"variable,omitempty"`        // declared variable
	Block          Block       `json:"block,omitempty"`           // Nested block
	OperationType  string      `json:"operation_type,omitempty"`  // Type of operation (e.g., addition, multiplication)
	Operands       []Statement `json:"Operands,omitempty"`        // List of variable used as Operands
	CalledFunction string      `json:"called_function,omitempty"` // function call
	Arguments      []Statement `json:"arguments,omitempty"`       // List of function call arguments
}

// --------------------------
// Define set structure
// --------------------------
// The set type is a type alias of `map[string]struct{}`
type set map[string]struct{}

// Adds a string to the set
func (s set) add(key string) {
	s[key] = struct{}{}
}

// Removes an animal from the set
func (s set) remove(key string) {
	delete(s, key)
}

// Returns a boolean value describing if the animal exists in the set
func (s set) has(key string) bool {
	_, ok := s[key]
	return ok
}

// Appends two sets together
func (s set) append(other set) {
	for key := range other {
		s.add(key)
	}
}

// -----------------------------------------
// Validate a program
// -----------------------------------------
/*
	The conditions for a program to be valid are:
		a. A function call must call a function that is declared in the same file.
		b. A variable can only be used in operations if it has been declared in a previous statement of the same block,
		or in case it has been declared in one of the previous statements of a surrounding block.

	Added conditions:
		-
*/

// IsValidFunctionCall validates a function call by checking the following conditions:
// - The function is already declared.
// - All arguments are valid operands
// - All variable arguments are both declared and assigned.
func IsValidFunctionCall(functionName string, arguments []Statement, declaredFunctionsMap map[string]bool, assignedVarMap map[string]bool) bool {
	// ensure function is already declared
	if !declaredFunctionsMap[functionName] {
		fmt.Println("Invalid function call due to calling undefined function: ", functionName)
		return false
	}
	// ensure that all arguments are valid operands.
	// Note: this will recursively call this function in case one of the operands is a function call as well
	for _, arg := range arguments {
		// note: all argumets should be both declared and assigned in the assignedVarMap.
		// Thus second param to function call is set to false "not an assigned var"
		if !IsValidOperand(arg, false, declaredFunctionsMap, assignedVarMap) {
			return false
		}
	}
	return true
}

// isValidOperand validates an operand by checking the following conditions:
// - If isAssignedVar is true, the operand must be of type "variable" for assignment.
// - For numerical operands, it checks if the value can be converted to a float, integers are accepted as well.
// - For variable operands, it checks if the variable is declared and assigned (unless isAssignedVar is true).
// - For function call operands and operation operands, it recursively checks the validity of the statement using isValidStatement.
//
func IsValidOperand(operand Statement, isAssignedVar bool, declaredFunctionsMap map[string]bool, assignedVarMap map[string]bool) bool {
	// early handle special common case: in assignment, first operand must be a var
	if isAssignedVar && operand.Type != "variable" {
		fmt.Println("Invalid operand. Expected left hand side of assignment to be variable but recieved: ", operand.Value)
		return false
	}

	switch operand.Type {
	case "numerical":
		if _, err := strconv.ParseFloat(operand.Value, 64); err != nil {
			fmt.Printf("Invalid operand, expected numerical type and value %v couldn't be converted\n", operand.Value)
			return false
		}
	case "variable":
		// a variable should be declared to be used in an operation.
		// The assignment is checked in the calling function to check on the opertaion type, if it is an assignment operation or sth else
		assigned, declared := assignedVarMap[operand.Variable]
		if !declared {
			fmt.Printf("Invalid operand, variable: %v is not declared\n", operand.Variable)
			return false
		}
		if !isAssignedVar && !assigned {
			fmt.Printf("Invalid operand, variable: %v is used without assignment\n", operand.Variable)
			return false
		}
	case "function_call":
	case "operation":
		// function calls and operations are statements
		if !IsValidStatement(operand, declaredFunctionsMap, assignedVarMap) {
			fmt.Println("Invalid operand, due to invalid statement:", operand)
			return false
		}
	default:
		fmt.Println("Invalid operation operand type: ", operand.Type)
		return false
	}

	return true
}

// IsValidStatement checks the validity of a statement by calling the corresponding validating function
// based on the statement type
func IsValidStatement(statement Statement, declaredFunctionsMap map[string]bool, assignedVarMap map[string]bool) bool {
	switch statement.Type {
	case "block":
		if !ValidateBlock(statement.Block, declaredFunctionsMap, assignedVarMap) {
			fmt.Println("Invalid block: ", statement.Block)
			return false
		}
	case "variable_declaration":
		if _, declared := assignedVarMap[statement.Variable]; declared {
			fmt.Println("Invalid variable declaration: ", statement.Variable, " variable already declared")
			return false
		}
		// add variable to the assignment map as false, since it now exists in the map it means it is already declared
		assignedVarMap[statement.Variable] = false
	case "operation":
		// TODO: in case of assignment, ensure that the operands list has length 2
		// we can create a separate validation function for each operation type
		for i, operand := range statement.Operands {
			// in assignment operation, the assigned variable is the first
			isAssignedVar := (i == 0 && statement.OperationType == "assignment")
			if !IsValidOperand(operand, isAssignedVar, declaredFunctionsMap, assignedVarMap) {
				return false
			}
		}
		// if the operation is an assignment operation
		if statement.OperationType == "assignment" {
			assignedVarMap[statement.Operands[0].Variable] = true
		}
	case "function_call":
		if !IsValidFunctionCall(statement.CalledFunction, statement.Arguments, declaredFunctionsMap, assignedVarMap) {
			return false
		}
	}
	return true
}

func ValidateBlock(block Block, declaredFunctionsMap map[string]bool, assignedVarMap map[string]bool) bool {
	for _, statement := range block.Statements {
		if !IsValidStatement(statement, declaredFunctionsMap, assignedVarMap) {
			return false
		}
	}
	return true
}

func ValidateProgramRec(program Program) bool {
	// Create a map to store all function declarations
	functionMap := make(map[string]bool)

	for _, function := range program.Functions {
		functionMap[function.Name] = true
	}

	for _, function := range program.Functions {
		// initialize a map to help with checking declared variables and assigned variables
		// declared variables are set to false and assigned varaiables are set to true
		// The map is initialized with the arguments passed to the function
		assignedVarMap := make(map[string]bool)
		for _, arg := range function.Parameters {
			assignedVarMap[arg] = true
		}
		if !ValidateBlock(function.Body, functionMap, assignedVarMap) {
			return false
		}
	}

	return true
}

// -----------------------------------------
// list declared but unused variables
// -----------------------------------------

func UnusedVariables(program Program) []string {
	// define a map to be populated with declared and used variables as follows:
	// if variable is declared, add it to the map with value: false
	// if it is used, set the value to true
	usedVariables := make(map[string]bool)

	// Traverse each function in the program to get all used variables
	for _, function := range program.Functions {
		PopulateUsedVariablesInBlock(function.Body, usedVariables)
	}

	unusedVariables := []string{}

	// Check for unused variables
	for variable := range usedVariables {
		if !usedVariables[variable] {
			unusedVariables = append(unusedVariables, variable)
		}
	}

	return unusedVariables
}

// PopulateUsedVariablesInStatement populates a given map as following:
// 	 if variable is declared, add it to the map with value: false
// 	 if it is used, set the value to true. A variable is used if:
// 		- used in an operation other than the left hand side of the assignment opertaion
// 		- used in the argument to a function call
func PopulateUsedVariablesInStatement(statement Statement, usedVariables map[string]bool) {
	switch statement.Type {
	case "variable_declaration":
		usedVariables[statement.Variable] = false
	case "operation":
		for i, operand := range statement.Operands {
			switch operand.Type {
			case "variable":
				// declare all variables in an operation to be used except the assigned variable
				isAssignedVar := (i == 0 && statement.OperationType == "assignment")
				if !isAssignedVar {
					usedVariables[operand.Variable] = true
				}
			case "function_call":
			case "operation":
				PopulateUsedVariablesInStatement(operand, usedVariables)

			}
		}
	case "function_call":
		for _, arg := range statement.Arguments {
			switch arg.Type {
			case "variable":
				usedVariables[arg.Variable] = true
			case "function_call":
			case "operation":
				PopulateUsedVariablesInStatement(arg, usedVariables)

			}
		}
	case "block":
		PopulateUsedVariablesInBlock(statement.Block, usedVariables)
	}

}

func PopulateUsedVariablesInBlock(block Block, usedVariables map[string]bool) {
	// Traverse each statement in the block
	for _, statement := range block.Statements {
		PopulateUsedVariablesInStatement(statement, usedVariables)
	}
}

// -----------------------------------------
// list functions dependancies
// -----------------------------------------
// FindFunctionCalls populates a dictionary with 1st level function dependencies
// 	populated map has key: function name, val: list of direct called functions
func FindFunctionCalls(program Program) map[string]set {
	functionCalls := make(map[string][]string)
	// Iterate over each function in the program and populate functionCalls map
	for _, function := range program.Functions {
		GetFunctionCallsRecursively(function.Body.Statements, function.Name, functionCalls)
	}
	// unfold all dependancies
	rolled_out_dependancies := RollOutDependencies(functionCalls)

	return rolled_out_dependancies
}

func GetFunctionCallsRecursively(statements []Statement, currentFunction string, functionCalls map[string][]string) {
	// first add the current function to list of functions we have
	emptyList := []string{}
	functionCalls[currentFunction] = append(functionCalls[currentFunction], emptyList...)
	// Traverse each statement in the function body
	for _, statement := range statements {
		switch statement.Type {

		case "function_call":
			{
				functionCalls[currentFunction] = append(functionCalls[currentFunction], statement.CalledFunction)
				GetFunctionCallsRecursively(statement.Arguments, currentFunction, functionCalls)
			}
		case "block":
			{
				GetFunctionCallsRecursively(statement.Block.Statements, currentFunction, functionCalls)
			}
		case "operation":
			GetFunctionCallsRecursively(statement.Operands, currentFunction, functionCalls)
		}
	}
}

// RollOutDependencies: given a map of key: str, value: []string which are keys as well,
// the function rolls out the dependencies by adding the value of each key to the value list while eliminating duplicates.
// e.g: {A:[B], B:[C], C:[D]} -> {A:[B, C, D], B:[C, D], C:[D]}
func RollOutDependencies(dependencies map[string][]string) map[string]set {
	rolledOut := make(map[string]set)

	for key := range dependencies {
		visited := make(map[string]bool)
		RollOutHelper(dependencies, key, visited, rolledOut)
	}

	return rolledOut
}

func RollOutHelper(dependencies map[string][]string, key string, visited map[string]bool, rolledOut map[string]set) {
	visited[key] = true

	// ensure that the key is initialized in the rollout map
	// to include functions with no dependencies
	_, exists := rolledOut[key]
	if !exists {
		rolledOut[key] = set{}
	}
	// recursively add functions dependencies to current key
	for _, dep := range dependencies[key] {
		if !visited[dep] {
			RollOutHelper(dependencies, dep, visited, rolledOut)
		}
		rolledOut[key].add(dep)
		rolledOut[key].append(rolledOut[dep])
	}
}

// -----------------------------------------
// main function
// -----------------------------------------
func main() {
	// read arguments from command line
	filePath := flag.String("file", "", "Path to the JSON file")
	mode := flag.String("mode", "", "Mode of operation")
	flag.Parse()

	// Validate command line arguments
	if *filePath == "" {
		fmt.Println("File path is required.")
		os.Exit(1)
	}
	if *mode == "" {
		fmt.Println("Mode is required.")
		os.Exit(1)
	}

	// Read the JSON file
	jsonData, err := ioutil.ReadFile(*filePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		os.Exit(1)
	}

	// Parse the JSON into the AST structure
	var program Program
	err = json.Unmarshal(jsonData, &program)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Print the parsed AST
	// fmt.Printf("%v\n", program)

	switch *mode {
	case "verify":
		// Verify the program
		isValid := ValidateProgramRec(program)
		fmt.Println("Is program valid?", isValid)
	case "unused_variables":
		unusedVariables := UnusedVariables(program)
		fmt.Println("unusedVariables: ", unusedVariables)
	case "functions_dependancies":
		functions_dependancies := FindFunctionCalls(program)
		fmt.Println("functions_dependancies: ", functions_dependancies)
	default:
		fmt.Println("Please enter a valid mode")
	}

}

/*
TODO:
3- include unit tests for implemented functions
6- code optimization / enhancements
*/
