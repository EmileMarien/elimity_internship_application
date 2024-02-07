package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReadTokenFromFile(t *testing.T) {
	// Create a temporary token file
	tempTokenFile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempTokenFile.Name())

	// Write token to the temporary file
	token := "mytoken"
	_, err = tempTokenFile.WriteString(token)
	if err != nil {
		t.Fatalf("Failed to write token to temp file: %v", err)
	}
	tempTokenFile.Close()

	// Test reading the token from the file
	readToken, err := readTokenFromFile(tempTokenFile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if readToken != token {
		t.Errorf("Expected token %s, got %s", token, readToken)
	}
}
