package client

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func (c *Client) Import(input string) error {
	fmt.Println("Import started")
	bytes, err := os.ReadFile(input)
	if err != nil {
		log.Debug(err)
		return err
	}

	exportData := ExportData{}
	err = json.Unmarshal(bytes, &exportData)
	if err != nil {
		log.Debug(err)
		return err
	}

	for _, a := range exportData.Activities {
		_, err := c.ots.AddActivity(a)
		if err != nil {
			log.Debug(err)
			return err
		}
	}

	for _, h := range exportData.Holidays {
		_, err := c.ots.AddHoliday(h)
		if err != nil {
			log.Debug(err)
			return err
		}
	}

	for _, h := range exportData.WorkDays {
		_, err := c.ots.AddWorkDay(h)
		if err != nil {
			log.Debug(err)
			return err
		}
	}

	fmt.Println("Import finished")
	return nil
}
