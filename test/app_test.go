package test

import "testing"

func TestNewTestApp(t *testing.T) {
	app, err := NewTestApp()
	if app == nil || err != nil {
		t.Fatal(err)
	}

	if err := app.Cleanup(); err != nil {
		t.Fatal(err)
	}
}
