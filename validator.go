package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Program struct {
	Functions []Function `json:"functions"`
}

type Function struct {
	Name      string   `json:"name"`
	Arguments []string `json:"arguments"`
	Body      Block    `json:"body"`
	Calls     []string `json:"calls"`
	Declared  []string `json:"declared"`
}

type Block struct {
	Statements []Statement `json:"statements"`
}

type Statement struct {
	Type          string   `json:"type"`
	Variables     []string `json:"variables,omitempty"`
	Uses          []string `json:"uses,omitempty"`
	Block         Block    `json:"block,omitempty"`
	Calls         []string `json:"calls,omitempty"`
	Parameters    []string `json:"arguments,omitempty"`
	OperationType string   `json:"operation_type,omitempty"`
	AssignTo      string   `json:"assign_to,omitempty"`
}

func main() {
	// Read the JSON file
	jsonData, err := ioutil.ReadFile("program.json")
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
	fmt.Printf("%+v\n", program)
}
