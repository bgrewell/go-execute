package utilities

import (
	"reflect"
	"testing"
)

func TestFieldsSplitsStringIntoFieldsRespectingQuotes(t *testing.T) {
	input := `"Hello, World!" 'This is a test'`
	expected := []string{`"Hello, World!"`, `'This is a test'`}
	result, err := Fields(input)
	if err != nil {
		t.Fatalf("Error splitting string: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unexpected result: %v", result)
	}
}

func TestFieldsHandlesEscapedQuotes(t *testing.T) {
	input := `"Hello, \"World!\""`
	expected := []string{`"Hello, \"World!\""`}
	result, err := Fields(input)
	if err != nil {
		t.Fatalf("Error splitting string: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unexpected result: %v", result)
	}
}

func TestFieldsReturnsErrorForUnclosedQuote(t *testing.T) {
	input := `"Hello, World!`
	_, err := Fields(input)
	if err == nil {
		t.Fatalf("Expected error for unclosed quote, got nil")
	}
}

func TestFieldsHandlesEmptyString(t *testing.T) {
	input := ""
	expected := []string{}
	result, err := Fields(input)
	if err != nil {
		t.Fatalf("Error splitting string: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Unexpected result: %v", result)
	}
}
