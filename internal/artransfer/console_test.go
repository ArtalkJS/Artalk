package artransfer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsole(t *testing.T) {
	// Test NewConsole
	t.Run("NewConsole", func(t *testing.T) {
		c := NewConsole()
		assert.NotNil(t, c)
		assert.False(t, c.IsOutputFuncSet())
	})

	// Test SetOutputFunc and IsOutputFuncSet
	t.Run("SetOutputFunc", func(t *testing.T) {
		c := NewConsole()
		outputFunc := func(string) {}
		c.SetOutputFunc(outputFunc)
		assert.True(t, c.IsOutputFuncSet())
	})

	// Test Error method
	t.Run("Error", func(t *testing.T) {
		var output string
		outputFunc := func(msg string) {
			output = msg
		}
		c := NewConsole()
		c.SetOutputFunc(outputFunc)
		c.Error("error message")
		assert.Equal(t, "[E] error message\n", output)
	})

	// Test Fatal method
	t.Run("Fatal", func(t *testing.T) {
		var output string
		outputFunc := func(msg string) {
			output = msg
		}
		c := NewConsole()
		c.SetOutputFunc(outputFunc)
		c.Fatal("fatal message")
		assert.Equal(t, "[F] fatal message\n", output)
	})

	// Test Warn method
	t.Run("Warn", func(t *testing.T) {
		var output string
		outputFunc := func(msg string) {
			output = msg
		}
		c := NewConsole()
		c.SetOutputFunc(outputFunc)
		c.Warn("warn message")
		assert.Equal(t, "[W] warn message\n", output)
	})

	// Test Info method
	t.Run("Info", func(t *testing.T) {
		var output string
		outputFunc := func(msg string) {
			output = msg
		}
		c := NewConsole()
		c.SetOutputFunc(outputFunc)
		c.Info("info message")
		assert.Equal(t, "[I] info message\n", output)
	})

	// Test Print method
	t.Run("Print", func(t *testing.T) {
		var output string
		outputFunc := func(msg string) {
			output = msg
		}
		c := NewConsole()
		c.SetOutputFunc(outputFunc)
		c.Print("print message")
		assert.Equal(t, "print message", output)
	})

	// Test Printf method
	t.Run("Printf", func(t *testing.T) {
		var output string
		outputFunc := func(msg string) {
			output = msg
		}
		c := NewConsole()
		c.SetOutputFunc(outputFunc)
		c.Printf("formatted %s", "message")
		assert.Equal(t, "formatted message", output)
	})

	// Test Println method
	t.Run("Println", func(t *testing.T) {
		var output string
		outputFunc := func(msg string) {
			output = msg
		}
		c := NewConsole()
		c.SetOutputFunc(outputFunc)
		c.Println("line message")
		assert.Equal(t, "line message\n", output)
	})

	// Test PrintTable method
	t.Run("PrintTable", func(t *testing.T) {
		var output string
		outputFunc := func(msg string) {
			output += msg
		}
		c := NewConsole()
		c.SetOutputFunc(outputFunc)
		rows := [][]any{{"col1", "col2"}, {"val1", "val2"}}
		c.PrintTable(rows)
		expectedOutput := "-------------------------\n + col1: col2\n + val1: val2\n-------------------------\n"
		assert.Equal(t, expectedOutput, output)
	})

	// Test PrintEncodeData method
	t.Run("PrintEncodeData", func(t *testing.T) {
		var output string
		outputFunc := func(msg string) {
			output = msg
		}
		c := NewConsole()
		c.SetOutputFunc(outputFunc)
		data := map[string]int{"key": 123}
		c.PrintEncodeData("map", data)
		expectedOutput := "[map]\n\n   map[string]int{\"key\":123}\n\n"
		assert.Equal(t, expectedOutput, output)
	})

	// Test Confirm method
	t.Run("Confirm", func(t *testing.T) {
		// Simulate user input
		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()
		r, w, _ := os.Pipe()
		os.Stdin = r

		go func() {
			fmt.Fprintln(w, "y")
			w.Close()
		}()

		c := NewConsole()
		assert.True(t, c.Confirm("Do you confirm?"))
	})
}
