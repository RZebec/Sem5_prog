package inputOutput

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
	Interface for the input and output.
 */
type InputOutput interface {
	ReadEntry() string
	Print(text string)
}

/*
	Struct for the default input and output implementation.
 */
type DefaultInputOutput struct {
}

/*
	Capsulates reading from Console.
*/
func (d *DefaultInputOutput) ReadEntry() string {
	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n')
	value = strings.TrimRight(value, "\n")
	return value
}

/*
	Capuslates writing to Console.
*/
func (d *DefaultInputOutput) Print(text string) {
	fmt.Println(text)
}
