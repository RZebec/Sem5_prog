package inputOutput

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type InputOutput interface {
	ReadEntry() string
	Print(text string)
}

type DefaultInputOutput struct {
}

/*
capsulate reading from Console is better to test
*/
func (d *DefaultInputOutput) ReadEntry() string {
	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n')
	value = strings.TrimRight(value, "\n")
	return value
}

/*
capuslate writing to Console is better to test
*/
func (d *DefaultInputOutput) Print(text string) {
	fmt.Println(text)
}
