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

func FormatDay(t time.Time) string {
	return t.Format("2006-01-02")
}

func Today() time.Time {
	today, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		panic(err)
	}
	return today
}

func PreviousWorkday(now time.Time) time.Time {
	prev := now.AddDate(0, 0, -1)
	// i'm assume that there where no activities on weekends
	if prev.Weekday() == time.Sunday {
		prev = prev.AddDate(0, 0, -1)
	}
	if prev.Weekday() == time.Saturday {
		prev = prev.AddDate(0, 0, -1)
	}
	return prev
}

func UniqueStrings(length int, cb func(i int) string) []string {
	resultMap := map[string]bool{}
	for i := 0; i < length; i++ {
		val := cb(i)
		if _, exists := resultMap[val]; !exists {
			resultMap[val] = true
		}
	}
	result := make([]string, len(resultMap))
	i := 0
	for val := range resultMap {
		result[i] = val
		i++
	}
	return result
}
