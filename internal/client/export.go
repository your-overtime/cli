package client

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func (c *Client) Export(since *time.Time, output string) error {
	fmt.Println("Export started")
	now := time.Now()
	if since == nil {
		yStart := time.Date(now.Year(), 0, 0, 0, 0, 0, 0, now.Location())
		since = &yStart
	}

	exportData := ExportData{}

	acs, err := c.ots.GetActivities(*since, now)
	if err != nil {
		log.Debug(err)
		return err
	}
	exportData.Activities = acs

	hds, err := c.ots.GetHolidays(*since, now)
	if err != nil {
		log.Debug(err)
		return err
	}
	exportData.Holidays = hds

	wds, err := c.ots.GetWorkDays(*since, now)
	if err != nil {
		log.Debug(err)
		return err
	}
	exportData.WorkDays = wds

	bytes, err := json.MarshalIndent(&exportData, "", " ")
	if err != nil {
		log.Debug(err)
		return err
	}
	err = os.WriteFile(output, bytes, fs.FileMode(0775))
	if err != nil {
		log.Debug(err)
		return err
	}
	fmt.Println("Export finished")
	return nil
}
