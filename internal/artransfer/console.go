package artransfer

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Console struct {
	outputFunc func(string)
}

func NewConsole() *Console {
	return &Console{}
}

func (c *Console) SetOutputFunc(f func(string)) {
	c.outputFunc = f
}

func (c *Console) IsOutputFuncSet() bool {
	return c.outputFunc != nil
}

func (c *Console) Error(a ...any) {
	if c.outputFunc != nil {
		c.outputFunc("[E] " + fmt.Sprintln(a...))
	} else {
		log.Error(a...)
	}
}

func (c *Console) Fatal(a ...any) {
	if c.outputFunc != nil {
		c.outputFunc("[F] " + fmt.Sprintln(a...))
	} else {
		log.Fatal(a...)
	}
}

func (c *Console) Warn(a ...any) {
	if c.outputFunc != nil {
		c.outputFunc("[W] " + fmt.Sprintln(a...))
	} else {
		log.Warn(a...)
	}
}

func (c *Console) Info(a ...any) {
	if c.outputFunc != nil {
		c.outputFunc("[I] " + fmt.Sprintln(a...))
	} else {
		log.Info(a...)
	}
}

func (c *Console) Print(a ...any) {
	if c.outputFunc != nil {
		c.outputFunc(fmt.Sprint(a...))
	} else {
		fmt.Print(a...)
	}
}

func (c *Console) Printf(format string, a ...any) {
	c.Print(fmt.Sprintf(format, a...))
}

func (c *Console) Println(a ...any) {
	c.Print(fmt.Sprintln(a...))
}

func (c *Console) PrintTable(rows [][]any) {
	if c.outputFunc != nil {
		c.Println("-------------------------")
		for _, row := range rows {
			l := len(row)
			c.Print(" + ")
			for i, col := range row {
				c.Print(col)
				if i < l-1 {
					c.Print(": ")
				}
			}
			c.Println()
		}
		c.Println("-------------------------")
		return
	}

	t := table.NewWriter()

	for _, r := range rows {
		t.AppendRow(r)
	}

	tStyle := table.StyleLight
	tStyle.Options.SeparateRows = true
	t.SetStyle(tStyle)

	c.Println(t.Render())
}

func (c *Console) PrintEncodeData(dataType string, val any) {
	c.Print(fmt.Sprintf("[%s]\n\n   %#v\n\n", dataType, val))
}

func (c *Console) Confirm(s string) bool {
	if c.outputFunc != nil {
		c.Printf("%s [y/n]: y", s)
		return true
	}

	r := bufio.NewReader(os.Stdin)

	for {
		c.Printf("%s [y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		resp := strings.ToLower(strings.TrimSpace(res))
		if resp == "y" || resp == "yes" {
			return true
		} else if resp == "n" || resp == "no" {
			return false
		}
	}
}
