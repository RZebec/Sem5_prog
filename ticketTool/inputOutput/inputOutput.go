package inputOutput

import (
	"bufio"
	"os"
	"strings"
)

func ReadEntry() string {
	reader := bufio.NewReader(os.Stdin)
	value, _ := reader.ReadString('\n')
	value = strings.TrimRight(value, "\n")
	return value
}
