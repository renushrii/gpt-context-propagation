# README.md for Go Code Contextual Modification Program

## Overview
This Go program is designed to refactor code, specifically focusing on replacing instances of context.Background() with the appropriate context passed from parent functions. This process avoids direct changes in logic, ensuring that only the intended modifications are made.

## Prerequisites
- Go language setup
- Access to OpenAI API for processing the code with GPT models

## How to Use
- git clone https://github.com/DDDecade0715/simple-chat ./files
- Clone the repository and navigate to the directory.
- Ensure that your Go environment is properly set up.
- Modify the directory variable to point to your Go files directory.
- Run the program with 'go run .'

## Explaination
- First it  fillters all .go extansion files then find those files which has context.Background() string.
- Changing the code by passing the context to parent function for further functions who need it.
- In response we get jason format of modified functions and modified code if exist.
- then we repeate again with those modified functions to pass context.
- Maintains the original logic of the code, only modifying context-related parts.
- Outputs a JSON object with 'modifiedFunctions' and 'code' keys.

## Important Notes
- Always ensure the context.Background() is used appropriately, and refactor only when necessary.
- You must configure the prompts or make small changes according to the model of GPT you are using for the code analysis. This prompt is valid for gpt4
- Do not change any single line or word other than what is required for the context propagation.
- Escaped quotes are required in the prompt, but do not escape new lines and tab characters.
