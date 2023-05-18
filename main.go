package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	validator "validator/validator"
)

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
	var program validator.Program
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
		isValid := validator.ValidateProgramRec(program, true)
		fmt.Println("Is program valid?", isValid)
	case "unused_variables":
		unusedVariables := validator.UnusedVariables(program)
		fmt.Println("unusedVariables: ", unusedVariables)
	case "functions_dependancies":
		functions_dependancies := validator.FindFunctionCalls(program)
		fmt.Println("functions_dependancies: ", functions_dependancies)
	default:
		fmt.Println("Please enter a valid mode")
	}

}
