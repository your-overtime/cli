package client

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/your-overtime/api/pkg"
)

func (c *Client) AddActivity(desc string, start *time.Time, end *time.Time) (*pkg.Activity, error) {
	return c.ots.AddActivity(pkg.Activity{
		InputActivity: pkg.InputActivity{
			Start:       start,
			End:         end,
			Description: desc,
		},
	})
}

func (c *Client) UpdateActivity(desc string, start *time.Time, end *time.Time, id uint) (*pkg.Activity, error) {
	ca, err := c.ots.GetActivity(id)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	if len(desc) > 0 {
		ca.Description = desc
	}
	if start != nil {
		ca.Start = start
	}
	if end != nil {
		ca.End = end
	}
	return c.ots.UpdateActivity(*ca)
}

func (c *Client) DeleteActivity(id uint) error {
	return c.ots.DelActivity(id)
}

func (c *Client) ImportKimai(filePath string) error {
	log.Debug(filePath)
	fr, err := os.Open(filePath)
	if err != nil {
		return err
	}
	r := csv.NewReader(fr)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		date := record[0]
		if len(date) != 10 {
			continue
		}

		start, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", date, record[1]))
		if err != nil {
			return nil
		}
		end, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", date, record[2]))
		if err != nil {
			return nil
		}
		desc := record[10]
		// TODO return activities to print them again
		_, err = c.AddActivity(desc, &start, &end)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (c *Client) StartActivity(desc string) (*pkg.Activity, error) {
	start := time.Now()
	return c.ots.AddActivity(pkg.Activity{
		InputActivity: pkg.InputActivity{
			Start:       &start,
			Description: desc,
			End:         nil,
		},
	})
}

func (c *Client) StopActivity() (*pkg.Activity, error) {
	a, err := c.ots.StopRunningActivity()
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	if a.Start == nil {
		return nil, errNoActiviyRunning
	}
	return a, nil
}

func (c *Client) GetActivities(start time.Time, end time.Time) ([]pkg.Activity, error) {
	return c.ots.GetActivities(start, end)
}
