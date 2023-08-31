package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructToMap(t *testing.T) {
	type TestStruct struct {
		Name  string
		Value float64
	}

	input := TestStruct{
		Name:  "John",
		Value: 42,
	}

	expected := map[string]interface{}{
		"Name":  "John",
		"Value": float64(42),
	}

	result := StructToMap(input)
	assert.Equal(t, expected, result)
}

func TestStructToFlatDotMap(t *testing.T) {
	type TestNestedStruct struct {
		NestedValue string
	}

	type TestStruct struct {
		Name   string
		Value  float64
		Nested TestNestedStruct
	}

	input := TestStruct{
		Name:  "John",
		Value: 42,
		Nested: TestNestedStruct{
			NestedValue: "nested",
		},
	}

	expected := map[string]interface{}{
		"Name":               "John",
		"Value":              float64(42),
		"Nested.NestedValue": "nested",
	}

	result := StructToFlatDotMap(input)

	assert.Equal(t, expected, result)
}

func TestCopyStruct(t *testing.T) {
	src := map[string]interface{}{
		"Name":  "John",
		"Value": int64(42),
	}

	var dest map[string]interface{}

	err := CopyStruct(&src, &dest)
	if err != nil {
		t.Errorf("Error copying struct: %v", err)
		return
	}

	assert.Equal(t, src, dest)
}
