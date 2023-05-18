# Overview
The tool deals with a simplified representation for an imperative programming language that allows static program analysis to be applied. 
This representation is defined as follows:
- A program contains one or more function declarations.
- A function declaration contains a block specifying the function body.
- There are four kinds of statements: blocks, variable declarations, operations and function calls:
    - A block contains zero or more statements.
    - A variable declaration declares a variable that may be used in an operation.
    - An operation has zero or more variable uses.
    - A function call references a function declaration.


The tool that expects a file specifying a program for the previously defined representation. The tool should offer three different operations:
1. Verify that a program is valid. The conditions for a program to be valid are:
    - A function call must call a function that is declared in the same file.
    - A variable can only be used in operations if it has been declared in a previous  statement of the same block, or in case it has been declared in
one of the previous statements of a surrounding block.
2. List variables that are declared but not used.
3. For each function, list which other functions they depend on. A function depends on another function if it is directly or indirectly called. For example, if function A calls Function B and Function B calls function C, then function A depends on both
functions B and C.

---
# Assumptions:

- variables are dynamically typed
- operands in an operation or arguments in a function call can be:
    - numerical
    - variable
    - function call
    - operation
- variables can't be declared twice
- function's parameters are considered as declaration for variable and they are already assigned
- in Assignment operation, the assigned variable is the first variable in the operations list
- in Assignment operation, the parser correctly parse the input into exactly two operands where the first one is the assigned variable. Otherwise throws error
- When operating in function dependancies or unused variables modes, it is assumed that the program is already valid:
    - for example the same variable is not declared twice


---
# JSON formate: 

Below, you can find a sample representation for the AST of the program as json file.
Note: sampe data are provided in the data folder

```json
{
    "functions": [
        {
            "name": "main",
            "parameters": [],
            "body": {
                "statements": [
                    {
                        "type": "block",
                        "block": {
                            "statements": [
                                {
                                    "type": "variable_declaration",
                                    "variable": "x"
                                },
                                {
                                    "type": "operation",
                                    "operation_type": "assignment",
                                    "operands": [
                                        {
                                            "type": "variable",
                                            "variable": "x"
                                        },
                                        {
                                            "type": "numerical",
                                            "value": "10"
                                        }
                                    ]
                                },
                                {
                                    "type": "function_call",
                                    "called_function": "printNumber",
                                    "arguments": [
                                        {
                                            "type": "function_call",
                                            "called_function": "calculateSum",
                                            "arguments": [
                                                {
                                                    "type": "numerical",
                                                    "value": "2"
                                                },
                                                {
                                                    "type": "numerical",
                                                    "value": "10"
                                                }
                                            ]
                                        }
                                    ]
                                }
                            ]
                        }
                    }
                ]
            }
        },
        {
            "name": "printNumber",
            "parameters": [
                "num"
            ],
            "body": {
                "statements": [
                    {
                        "type": "function_call",
                        "called_function": "display",
                        "arguments": [
                            {
                                "type": "variable",
                                "variable": "num"
                            }
                        ]
                    }
                ]
            }
        }
    ]
}

```

---
# How to run
- To run the tool use the following command line:
>`go run main.go -file <json_file_path> -mode <operation_mode>`

Replace the json file path, with the actual path, and the operation mode is one of the following:
- `verify`
- `unused_variables`
- `functions_dependancies`

ex: 
>`go run main.go -file './data/valid/operations.json' -mode 'verify'`


To run tests:
> `go test -v ./validator/`
---