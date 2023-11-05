package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	directory := "/home/amieo/Documents/Assignment/files"

	files, err := listFiles(directory)
	if err != nil {
		panic(err)
	}

	modifiedFunctions := []string{}
	for _, file := range files {
		code, err := ReadCode(file)
		if err != nil {
			fmt.Printf("Error reading code from file %s: %v\n", file, err)
			continue
		}

		if !strings.Contains(code, "context.Background()") && len(modifiedFunctions) == 0 {
			continue
		}

		log.Printf("asking gpt to check %s with modified functions %v \n", file, modifiedFunctions)
		prompt1 := ctxBackgroundPrompt(code)
		response, err := sendAndGetResponse(prompt1)
		if err != nil {
			fmt.Printf("Error sending code to OpenAI: %v\n", err)
			panic(err)
		}

		reWrite(file, response.Code)
		modifiedFunctions = append(modifiedFunctions, response.ModifiedFunctions...)

		if len(modifiedFunctions) == 0 {
			continue
		}

		for _, checkFile := range files {
			log.Printf("**asking gpt to check %s with modified functions %v \n", checkFile, modifiedFunctions)
			code, err := ReadCode(checkFile)
			if err != nil {
				fmt.Printf("Error reading code from file %s: %v\n", checkFile, err)
				panic(err)
			}

			prompt2 := modifiedFunctionsPrompt(modifiedFunctions, code)
			response, err := sendAndGetResponse(prompt2)
			reWrite(checkFile, response.Code)
			log.Printf("changing modified functions current %v, appending %v", modifiedFunctions, response.ModifiedFunctions)
			modifiedFunctions = append(modifiedFunctions, response.ModifiedFunctions...)
		}
	}
	for _, checkFile := range files {
		log.Printf("**asking gpt to check %s with modified functions %v \n", checkFile, modifiedFunctions)
		code, err := ReadCode(checkFile)
		if err != nil {
			fmt.Printf("Error reading code from file %s: %v\n", checkFile, err)
			panic(err)
		}

		prompt2 := modifiedFunctionsPrompt(modifiedFunctions, code)
		response, err := sendAndGetResponse(prompt2)
		reWrite(checkFile, response.Code)
		log.Printf("changing modified functions current %v, appending %v", modifiedFunctions, response.ModifiedFunctions)
		modifiedFunctions = append(modifiedFunctions, response.ModifiedFunctions...)
	}

	fmt.Println("Done")
}

func ctxBackgroundPrompt(code string) string {
	escapedCode := strconv.Quote(code)
	prompt := `
	take a deep breath and solve step by step
	You are a golang expert.
	
	{
		"prompt": "Given the following Go code, ensure that all instances of 'context.Background()' are replaced with the appropriate
				   propagated context from parent functions, except within the 'main.main' function. Functions where the signature is changed to 
				   include 'context.Context' as a parameter should be listed under 'modifiedFunctions' with full name of function/methods with signature, these function's signature were modified for passing context from parent.
				   Return a JSON object with 'modifiedFunctions' and 'code' keys, 
				   where 'code' contains the modified or unmodified code, with escaped quotes and unescaped new lines and tabs. Don't omit the code always send complete code",
		"code": %s
	}

	Don't do the following....
	2. Don't make any other changes in logic of code.
	3. Don't give me expalination at all.
	4. don't make change in single line or word other than what i say to do.
	5. Don't escape new line and tab character in response code. Escape quotes.

	`

	return fmt.Sprintf(prompt, escapedCode)
}

func modifiedFunctionsPrompt(modifiedFunctions []string, code string) string {
	escapedCode := strconv.Quote(code)
	prompt := `
	Take a deep breath...
	You are a golang expert.

	{
		"prompt": "Given the list of functions that have had their signature changed to include 'context.Context', and the following Go code, modify the code to pass 'context.Context' to any 
		           function call that is now expecting it. List any new functions modified under 'modifiedFunctions'. Return a JSON object with 'modifiedFunctions' and 'code' keys, where 'code' 
				   contains the modified or unmodified code, with escaped quotes and unescaped new lines and tabs. Don't omit the code always send complete code",
		"modified_functions": ["%v"],
		"code": %s
	}

		Don't do the following....
		1. Don't make any other changes in logic of code.
		2. Don't pass context to any other functions other than list of modified functions.
		3. Don't give me expalination at all.
		4. if modified funtions are not called in code so don't change the code at all.
		5. Don't escape new line and tab character in response code. escape quote

	`
	return fmt.Sprintf(prompt, modifiedFunctions, escapedCode)
}
