package artransfer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var HttpOutputEchoCtx echo.Context // TODO 只支持单线执行，同时请求两个地方会出问题
var HttpOutput func(continueRun bool, text string)
var Assumeyes bool = false

func logError(a ...interface{}) {
	if HttpOutput != nil {
		HttpOutput(true, "[ERROR] "+fmt.Sprint(a...))
		return
	}

	logrus.Error(a...)
}

func logFatal(a ...interface{}) {
	if HttpOutput != nil {
		HttpOutput(false, "[FATAL] "+fmt.Sprint(a...))
		return
	}

	logrus.Fatal(a...)
}

func logWarn(a ...interface{}) {
	if HttpOutput != nil {
		HttpOutput(true, "[WARN] "+fmt.Sprint(a...))
		return
	}

	logrus.Warn(a...)
}

func logInfo(a ...interface{}) {
	if HttpOutput != nil {
		HttpOutput(true, "[INFO] "+fmt.Sprint(a...))
		return
	}

	logrus.Info(a...)
}

func print(a ...interface{}) {
	if HttpOutput != nil {
		HttpOutput(true, fmt.Sprint(a...))
		return
	}

	fmt.Print(a...)
}

func printf(format string, a ...interface{}) {
	print(fmt.Sprintf(format, a...))
}

func println(a ...interface{}) {
	print(fmt.Sprintln(a...))
}

func PrintTable(rows [][]interface{}) {
	if HttpOutput != nil {
		println("-------------------------")
		for _, row := range rows {
			l := len(row)
			print(" + ")
			for i, col := range row {
				print(col)
				if i < l-1 {
					print(": ")
				}
			}
			println()
		}
		println("-------------------------")
		return
	}

	t := table.NewWriter()

	for _, r := range rows {
		t.AppendRow(r)
	}

	tStyle := table.StyleLight
	tStyle.Options.SeparateRows = true
	t.SetStyle(tStyle)

	println(t.Render())
}

func PrintEncodeData(dataType string, val interface{}) {
	print(SprintEncodeData(dataType, val))
}

func SprintEncodeData(dataType string, val interface{}) string {
	return fmt.Sprintf("[%s]\n\n   %#v\n\n", dataType, val)
}

func Confirm(s string) bool {
	if Assumeyes {
		printf("%s [y/n]: y", s)
		return true
	}

	r := bufio.NewReader(os.Stdin)

	for {
		printf("%s [y/n]: ", s)

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
