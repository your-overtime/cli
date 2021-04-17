package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ReadTextFromStdin(text string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(text)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Debug(err)
	}
	return strings.ReplaceAll(input, "\n", "")
}
