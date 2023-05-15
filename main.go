package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

// Define the structure for the AST representation

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

func main() {
	// Read the JSON file
	jsonData, err := ioutil.ReadFile("./data/valid/sample_3.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
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

	// Verify the program
	isValid := VerifyProgramRec(program)
	fmt.Println("Program is valid:", isValid)
}

/*
TODO:
1- function caalls with a var: var must be assigned first
2- operands in operations doesn't need to be strings, could be function as well

*/
