package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func ReadTextFromStdin(text string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(text)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Debug(err)
	}
	return strings.TrimSpace(strings.ReplaceAll(input, "\n", ""))
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}
