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
	Name      string   `json:"name"`      // Name of the function
	Arguments []string `json:"arguments"` // List of function arguments
	Body      Block    `json:"body"`      // Function body
	// CalledFunction     []string `json:"called_function"`     // List of function calls
	// Declared  []string `json:"declared"`  // List of declared variables

}

// Block represents a block of statements.
type Block struct {
	// A block contains zero or more statements.
	Statements []Statement `json:"statements"` // List of statements in the block
}

// Statement represents an individual statement.
type Statement struct {
	Type             string   `json:"type"`                      // Type of statement (block, variable_declaration, operation, function_call)
	DeclaredVariable string   `json:"variable,omitempty"`        // declared variable
	Block            Block    `json:"block,omitempty"`           // Nested block
	OperationType    string   `json:"operation_type,omitempty"`  // Type of operation (e.g., addition, multiplication)
	Operands         []string `json:"Operands,omitempty"`        // List of variable used as Operands
	CalledFunction   string   `json:"called_function,omitempty"` // function call
	Parameters       []string `json:"arguments,omitempty"`       // List of function call parameters
	AssignTo         string   `json:"assign_to,omitempty"`       // Variable to assign the result of an operation
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

func isValidOperand(operand string, varMap map[string]bool) bool {
	if _, err := strconv.Atoi(operand); err != nil && !varMap[operand] {
		return false
	}
	return true
}

func validateBlock(block Block, functionMap map[string]bool, varMap map[string]bool) bool {
	/*
		INFO: function to recursively validate a block. Conditions for validity:
				a. A function call must call a function that is declared in the same file.
				b. A variable can only be used in operations if it has been declared in a previous statement of the same block,
				or in case it has been declared in one of the previous statements of a surrounding block.
		:params
			- block Block: a block from the program
			- functionMap map[string]bool: map contains defined functions, acts as set
			- varMap map[string]bool: map contains variables from bigger scope passed to this block
		:returns
			- boolean to indicate whether the block is valide or not as defined

	*/
	for _, statement := range block.Statements {
		switch statement.Type {
		case "block":
			if !validateBlock(statement.Block, functionMap, varMap) {
				fmt.Println("Invalid block: ", statement.Block)
				return false
			}
		case "variable_declaration":
			if _, ok := varMap[statement.DeclaredVariable]; ok {
				fmt.Println("Invalid variable declaration: ", statement.DeclaredVariable)
				return false
			}
			varMap[statement.DeclaredVariable] = true
		case "operation":
			for _, operand := range statement.Operands {
				if !isValidOperand(operand, varMap) {
					fmt.Println("Invalid operand", operand)
					return false
				}
			}
			if statement.AssignTo != "" && !varMap[statement.AssignTo] {
				// check if there is a var to assign to and that variable is already declared
				fmt.Println("Invalid assignee, var not declared: ", statement.AssignTo)
				return false
			}
			if statement.AssignTo != "" {
				// TODO: instead of setting the assignee as true in varMap, it should be done in the assignmentMap
				varMap[statement.AssignTo] = true
			}
		case "function_call":
			if _, ok := functionMap[statement.CalledFunction]; !ok {
				fmt.Println("Invalid function call due to calling undefined function: ", statement.CalledFunction)
				return false
			}
			for _, param := range statement.Parameters {
				if !isValidOperand(param, varMap) {
					fmt.Println("Invalid function call due to undeclared parameter: ", param)
					return false
				}
			}
			// TODO: revisit the next assignment part as it should be handled in operations and not here
			if statement.AssignTo != "" && !varMap[statement.AssignTo] {
				// check if there is a var to assign to and that variable is already declared
				fmt.Println("Invalid assignee after function call, var not declared: ", statement.AssignTo)
				return false
			}

			if statement.AssignTo != "" {
				varMap[statement.AssignTo] = true
			}
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
		varMap := make(map[string]bool)
		for _, arg := range function.Arguments {
			varMap[arg] = true
		}
		if !validateBlock(function.Body, functionMap, varMap) {
			return false
		}
	}

	return true
}

// -----------------------------------------
// list declared but unused variables
// -----------------------------------------

func findUnusedVariables(program Program) []string {
	declaredVariables := make(map[string]bool)
	usedVariables := make(map[string]bool)

	// Traverse each function in the program to get all used variables
	for _, function := range program.Functions {
		traverseBlock(function.Body, declaredVariables, usedVariables)
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

func traverseBlock(block Block, declaredVariables map[string]bool, usedVariables map[string]bool) {
	// Traverse each statement in the block
	for _, statement := range block.Statements {
		switch statement.Type {
		case "variable_declaration":
			declaredVariables[statement.DeclaredVariable] = true
		case "operation":
			for _, operand := range statement.Operands {
				// TODO: check on the operand type to ensure it's a var
				usedVariables[operand] = true
			}
			// if statement.AssignTo != "" {
			// 	usedVariables[statement.AssignTo] = true
			// }
		case "function_call":
			for _, param := range statement.Parameters {
				usedVariables[param] = true
			}
			// if statement.AssignTo != "" {
			// 	usedVariables[statement.AssignTo] = true
			// }
		case "block":
			traverseBlock(statement.Block, declaredVariables, usedVariables)
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
		if statement.Type == "function_call" {
			functionCalls[currentFunction] = append(functionCalls[currentFunction], statement.CalledFunction)
		} else if statement.Type == "block" {
			getFunctionCallsRecursively(statement.Block.Statements, currentFunction, functionCalls)
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
1- function caalls with a var: var must be assigned first
2- operands in operations doesn't need to be strings, could be function as well
3- include unit tests for implemented functions
4- Consistent code documentation
5- Consistent conventions
6- code optimization / enhancements
*/
