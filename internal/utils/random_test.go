package utils

import "testing"

func TestRandomString(t *testing.T) {
	length := 10

	randomStr := RandomString(length)
	if len(randomStr) != length {
		t.Errorf("Expected random string length %d, but got %d", length, len(randomStr))
	}
}

func TestRandomStringWithAlphabet(t *testing.T) {
	length := 15
	alphabet := "ABC123"

	randomStr := RandomStringWithAlphabet(length, alphabet)
	if len(randomStr) != length {
		t.Errorf("Expected random string length %d, but got %d", length, len(randomStr))
	}
}

func TestPseudorandomString(t *testing.T) {
	length := 20

	pseudoRandomStr := PseudorandomString(length)
	if len(pseudoRandomStr) != length {
		t.Errorf("Expected pseudorandom string length %d, but got %d", length, len(pseudoRandomStr))
	}
}

func TestPseudorandomStringWithAlphabet(t *testing.T) {
	length := 25
	alphabet := "xyz789"

	pseudoRandomStr := PseudorandomStringWithAlphabet(length, alphabet)
	if len(pseudoRandomStr) != length {
		t.Errorf("Expected pseudorandom string length %d, but got %d", length, len(pseudoRandomStr))
	}
}
