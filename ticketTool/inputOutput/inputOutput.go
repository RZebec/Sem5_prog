package inputOutput

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type InputOutput interface{
	ReadEntry() string
	Print(text string)
}

type DefaultInputOutput struct {

}

func (d *DefaultInputOutput) ReadEntry() string {
	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n')
	value = strings.TrimRight(value, "\n")
	return value
}

func (d *DefaultInputOutput) Print(text string)  {
	fmt.Println(text)
}
