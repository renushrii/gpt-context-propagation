package main

import (
	"strings"
	"testing"
)

func TestPrompt(t *testing.T) {
	file := "./originals/pkg/elasticsearch/elasticsearch.go"

	code, err := ReadCode(file)
	if err != nil {
		t.Error(err)
	}

	ctxPrompt := ctxBackgroundPrompt(code)

	resp, err := sendAndGetResponse(ctxPrompt)
	if err != nil {
		t.Error(err)
	}

	if strings.Contains(resp.Code, "context.Background()") {
		t.Error("context.Background() not replaced")
	}

	t.Log(file)
	t.Log(resp.Code)

	file = "./originals/api/controller/ImController.go"
	code, err = ReadCode(file)
	if err != nil {
		t.Error(err)
	}

	modifiedFuncsPrompt := modifiedFunctionsPrompt(resp.ModifiedFunctions, code)
	resp, err = sendAndGetResponse(modifiedFuncsPrompt)
	if err != nil {
		t.Error(err)
	}

	if len(resp.ModifiedFunctions) != 0 {
		t.Error("modified functions not equal")
	}

	t.Log(file)
	t.Log(resp.Code)
}
