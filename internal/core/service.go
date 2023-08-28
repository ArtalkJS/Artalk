package core

import "fmt"

type Service interface {
	Init() error
	Dispose() error
}

func genServiceName[T any]() string {
	var t T

	// struct
	name := fmt.Sprintf("%T", t)
	if name != "<nil>" {
		return name
	}

	// interface
	return fmt.Sprintf("%T", new(T))
}
