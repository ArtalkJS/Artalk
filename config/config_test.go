package config

import "testing"

func TestInit(t *testing.T) {
	Init("../artalk-go.example.yaml")
	t.Log(Instance)
}
