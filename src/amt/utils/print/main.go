package print

import (
	"os"
	"fmt"
	"bufio"
)

import (
	"amt/utils/print/colors"
)

func Eprintln(err string) {
	fmt.Fprintf(os.Stderr, "%s%s%s\n", colors.YELLOW, err, colors.RESET)
}

func BufferedPrint(items []string) {
	writer := bufio.NewWriter(os.Stdout)

	writer.WriteString(colors.GREEN)

	for _, item := range items {
		writer.WriteString(item + "\n")
	}

	writer.WriteString(colors.RESET)

	writer.Flush()
}

func Cprintln(msg string, color string) {
	fmt.Printf("%s%s%s\n", color, msg, colors.RESET)
}