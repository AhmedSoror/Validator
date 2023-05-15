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
	// Calls     []string `json:"calls"`     // List of function calls
	// Declared  []string `json:"declared"`  // List of declared variables

}

// Block represents a block of statements.
type Block struct {
	// A block contains zero or more statements.
	Statements []Statement `json:"statements"` // List of statements in the block
}

// Statement represents an individual statement.
type Statement struct {
	Type             string   `json:"type"`                     // Type of statement (block, variable_declaration, operation, function_call)
	DeclaredVariable string   `json:"variable,omitempty"`       // declared variable
	Block            Block    `json:"block,omitempty"`          // Nested block
	OperationType    string   `json:"operation_type,omitempty"` // Type of operation (e.g., addition, multiplication)
	Operands         []string `json:"Operands,omitempty"`       // List of variable used as Operands
	Calls            string   `json:"calls,omitempty"`          // function call
	Parameters       []string `json:"arguments,omitempty"`      // List of function call parameters
	AssignTo         string   `json:"assign_to,omitempty"`      // Variable to assign the result of an operation
}

func validateBlock(block Block, functionMap map[string]bool, varMap map[string]bool) bool {
	for _, statement := range block.Statements {
		switch statement.Type {
		case "block":
			if !validateBlock(statement.Block, functionMap, varMap) {
				return false
			}
		case "variable_declaration":
			if _, ok := varMap[statement.DeclaredVariable]; ok {
				return false
			}
			varMap[statement.DeclaredVariable] = true
		case "operation":
			for _, operand := range statement.Operands {
				// if _, ok := varMap[operand]; !ok
				if _, err := strconv.Atoi(operand); err != nil && !varMap[operand] {
					return false
				}
			}
			if !varMap[statement.AssignTo] {
				return false
			}
			if statement.AssignTo != "" {
				varMap[statement.AssignTo] = true
			}
		case "function_call":
			if _, ok := functionMap[statement.Calls]; !ok {
				return false
			}
			for _, param := range statement.Parameters {
				if _, ok := varMap[param]; !ok {
					return false
				}
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
	fmt.Printf("%v\n", program)

	// Verify the program
	isValid := VerifyProgramRec(program)
	fmt.Println("Program is valid:", isValid)
}
