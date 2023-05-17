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
	// A function declaration contains a block specifying the function body.
	/*
		Extra added info:
			- function identifier "name"
			- function arguments
			- function calls
			- declared variables in the function
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
	// AssignTo       string      `json:"assign_to,omitempty"`       // Variable to assign the result of an operation
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
func isValidFunctionCall(functionName string, arguments []Statement, declaredFunctionsMap map[string]bool, assignedVarMap map[string]bool) bool {
	// to validate a function:
	// - function is already declared
	if !declaredFunctionsMap[functionName] {
		fmt.Println("Invalid function call due to calling undefined function: ", functionName)
		return false
	}
	// all arguments are valid operands. Note: this will recursively call this function
	// in case one of the operands is a function call as well
	for _, arg := range arguments {
		// note: all argumets should be not only declared but also assigned in the assignedVarMap.
		// Thus second param to function call is set to false "not an assigned var"
		if !isValidOperand(arg, false, declaredFunctionsMap, assignedVarMap) {
			return false
		}
	}
	return true
}

func isValidOperand(operand Statement, isAssignedVar bool, declaredFunctionsMap map[string]bool, assignedVarMap map[string]bool) bool {
	// early handle special common case: in assignment, first operand must be a var
	if isAssignedVar && operand.Type != "variable" {
		fmt.Println("Invalid operand. Expected left hand side of assignment to be variable but recieved: ", operand.Value)
		return false
	}

	switch operand.Type {
	case "numerical":
		if _, err := strconv.Atoi(operand.Value); err != nil {
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
		if !isValidStatement(operand, declaredFunctionsMap, assignedVarMap) {
			fmt.Println("Invalid operand, due to invalid statement:", operand)
			return false
		}
	default:
		fmt.Println("Invalid operation operand type: ", operand.Type)
		return false
	}

	return true
}

func isValidStatement(statement Statement, declaredFunctionsMap map[string]bool, assignedVarMap map[string]bool) bool {
	/*
		 Conditions for validity:
			a. A function call must call a function that is declared in the same file.
			b. A variable can only be used in operations if it has been declared in a previous statement of the same block,
			or in case it has been declared in one of the previous statements of a surrounding block.
	*/
	switch statement.Type {
	case "block":
		if !verifyBlock(statement.Block, declaredFunctionsMap, assignedVarMap) {
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
		for i, operand := range statement.Operands {
			// in assignment operation, the assigned variable is the first
			isAssignedVar := (i == 0 && statement.OperationType == "assignment")
			if !isValidOperand(operand, isAssignedVar, declaredFunctionsMap, assignedVarMap) {
				return false
			}
		}
		// if the operation is an assignment operation
		if statement.OperationType == "assignment" {
			assignedVarMap[statement.Operands[0].Variable] = true
		}
		// TODO: set assigned variables in a map to be used for verification

		// check if there is a var to assign to and that variable is already declared
		// if statement.AssignTo != "" {
		// 	if !declaredVarMap[statement.AssignTo] {
		// 		fmt.Println("Invalid assigned, var not declared: ", statement.AssignTo)
		// 		return false
		// 	}
		// 	declaredVarMap[statement.AssignTo] = true
		// }
	case "function_call":
		if !isValidFunctionCall(statement.CalledFunction, statement.Arguments, declaredFunctionsMap, assignedVarMap) {
			return false
		}
	}
	return true
}

func verifyBlock(block Block, declaredFunctionsMap map[string]bool, assignedVarMap map[string]bool) bool {
	/*
		INFO: function to recursively validate a block.
		:params
			- block Block: a block from the program
			- functionMap map[string]bool: map contains defined functions, acts as set
			- varMap map[string]bool: map contains variables from bigger scope passed to this block
		:returns
			- boolean to indicate whether the block is valide or not as defined

	*/
	for _, statement := range block.Statements {
		if !isValidStatement(statement, declaredFunctionsMap, assignedVarMap) {
			return false
		}
	}
	return true
}

func VerifyProgramRec(program Program) bool {
	//--------------------
	// A function call must call a function that is declared in the same file.
	//--------------------

	// Create a map to store function declarations
	functionMap := make(map[string]bool)

	// Iterate over function declarations and populate the function map
	for _, function := range program.Functions {
		functionMap[function.Name] = true
	}

	// for each function, init var map and calidate the function body
	for _, function := range program.Functions {
		// initialize a map to help with checking declared variables (if exists in the map) and assigned variables (exits and value is true)
		assignedVarMap := make(map[string]bool)
		for _, arg := range function.Parameters {
			assignedVarMap[arg] = true
		}
		if !verifyBlock(function.Body, functionMap, assignedVarMap) {
			return false
		}
	}

	return true
}

// -----------------------------------------
// list declared but unused variables
// -----------------------------------------

func findUnusedVariables(program Program) []string {
	// TODO: to save space, we can use only one map where false is decalred and true is used variable
	declaredVariables := make(map[string]bool)
	usedVariables := make(map[string]bool)

	// Traverse each function in the program to get all used variables
	for _, function := range program.Functions {
		populateUsedVariablesInBlock(function.Body, declaredVariables, usedVariables)
	}

	unusedVariables := []string{}

	// Check for unused variables
	for variable := range declaredVariables {
		if !usedVariables[variable] {
			unusedVariables = append(unusedVariables, variable)
		}
	}

	return unusedVariables
}

func populateUsedVariablesInBlock(block Block, declaredVariables map[string]bool, usedVariables map[string]bool) {
	// Traverse each statement in the block
	for _, statement := range block.Statements {
		switch statement.Type {
		case "variable_declaration":
			declaredVariables[statement.Variable] = true
		case "operation":
			for i, operand := range statement.Operands {
				// declare all variables in an operation to be used except the assigned variable
				isAssignedVar := (i == 0 && statement.OperationType == "assignment")
				if !isAssignedVar && operand.Type == "variable" {
					usedVariables[operand.Variable] = true
				}
			}
		case "function_call":
			for _, arg := range statement.Arguments {
				if arg.Type == "variable" {
					usedVariables[arg.Variable] = true
				}
			}
			// if statement.AssignTo != "" {
			// 	usedVariables[statement.AssignTo] = true
			// }
		case "block":
			populateUsedVariablesInBlock(statement.Block, declaredVariables, usedVariables)
		}
	}
}

// -----------------------------------------
// list functions dependancies
// -----------------------------------------
func findFunctionCalls(program Program) map[string]set {
	/*
		INFO: function to populate a dictionary with 1st level function dependencies where the key: function name, val: list of called functions
	*/
	functionCalls := make(map[string][]string)
	// Iterate over each function in the program and populate functionCalls map
	for _, function := range program.Functions {
		getFunctionCallsRecursively(function.Body.Statements, function.Name, functionCalls)
	}

	rolled_out_dependancies := rollOutDependencies(functionCalls)

	return rolled_out_dependancies
}

func getFunctionCallsRecursively(statements []Statement, currentFunction string, functionCalls map[string][]string) {
	/*
		INFO: given a function name, populate a dictionary with 1st level function dependencies where the key: function name, val: list of called functions
	*/
	// first add the current function to list of functions we have
	emptyList := []string{}
	functionCalls[currentFunction] = append(functionCalls[currentFunction], emptyList...)
	// Traverse each statement in the function body
	for _, statement := range statements {
		switch statement.Type {

		case "function_call":
			{
				functionCalls[currentFunction] = append(functionCalls[currentFunction], statement.CalledFunction)
				getFunctionCallsRecursively(statement.Arguments, currentFunction, functionCalls)
			}
		case "block":
			{
				getFunctionCallsRecursively(statement.Block.Statements, currentFunction, functionCalls)
			}
		case "operation":
			getFunctionCallsRecursively(statement.Operands, currentFunction, functionCalls)
		}
	}
}

// ---------------------------------

func rollOutDependencies(dependencies map[string][]string) map[string]set {
	/*
		INFO: given a map of key: str, value: []string which are keys as well,
			the function rolls out the dependencies by adding the value of each key to the value list while eliminating duplicates.
			e.g: {A:[B], B:[C], C:[D]} -> {A:[B, C, D], B:[C, D], C:[D]}
	*/
	rolledOut := make(map[string]set)

	for key := range dependencies {
		visited := make(map[string]bool)
		rollOutHelper(dependencies, key, visited, rolledOut)
	}

	return rolledOut
}

func rollOutHelper(dependencies map[string][]string, key string, visited map[string]bool, rolledOut map[string]set) {
	visited[key] = true

	// first ensure that the key is initialized in the rollout map
	// important to include functions with no dependencies
	_, exists := rolledOut[key]
	if !exists {
		rolledOut[key] = set{}
	}
	// recursively add functions dependencies to current key
	for _, dep := range dependencies[key] {
		if !visited[dep] {
			rollOutHelper(dependencies, dep, visited, rolledOut)
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
		isValid := VerifyProgramRec(program)
		fmt.Println("Is program valid?", isValid)
	case "unused_variables":
		unusedVariables := findUnusedVariables(program)
		fmt.Println("unusedVariables: ", unusedVariables)
	case "functions_dependancies":
		functions_dependancies := findFunctionCalls(program)
		fmt.Println("functions_dependancies: ", functions_dependancies)
	default:
		fmt.Println("Please enter a valid mode")
	}

}

/*
TODO:
3- include unit tests for implemented functions
4- Consistent code documentation
5- Consistent conventions
6- code optimization / enhancements
*/
