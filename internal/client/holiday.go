package client

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/your-overtime/api/pkg"
)

func (c *Client) AddHoliday(desc string, start time.Time, end time.Time, legalHoliday bool, sick bool) error {
	hType := pkg.HolidayTypeFree
	if legalHoliday {
		hType = pkg.HolidayTypeLegalHoliday
	} else if sick {
		hType = pkg.HolidayTypeSick
	}

	h, err := c.ots.AddHoliday(pkg.Holiday{
		Start:       start,
		End:         end,
		Description: desc,
		Type:        hType,
	}, pkg.Employee{})

	if err != nil {
		log.Debug(err)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	printHoliday(w, h)
	w.Flush()

	return nil
}

func (c *Client) GetHolidays(start time.Time, end time.Time, asJSON bool) error {
	hs, err := c.ots.GetHolidays(start, end, pkg.Employee{})

	if asJSON {
		jsonData, err := json.MarshalIndent(hs, "", " ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
		return nil
	}
	if err != nil {
		log.Debug(err)
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, '.', tabwriter.TabIndent)
	mins := 0
	for i, h := range hs {
		fmt.Fprintf(w, "No\t: %d\n", i+1)
		printHoliday(w, &h)
		fmt.Fprintln(w)
		mins += int(h.End.Sub(h.Start).Minutes())
	}
	fmt.Fprintf(w, "Duration\t: %s\n", formatMinutes(int64(mins)))
	w.Flush()
	return nil
}

func (c *Client) UpdateHoliday(desc string, start *time.Time, end *time.Time, id uint, legalHoliday bool, sick bool) error {
	ch, err := c.ots.GetHoliday(id, pkg.Employee{})
	if err != nil {
		log.Debug(err)
		return err
	}

	if legalHoliday {
		ch.Type = pkg.HolidayTypeLegalHoliday
	} else if sick {
		ch.Type = pkg.HolidayTypeSick
	} else {
		ch.Type = pkg.HolidayTypeFree
	}

	if len(desc) > 0 {
		ch.Description = desc
	}
	if start != nil {
		ch.Start = *start
	}
	if end != nil {
		ch.End = *end
	}
	h, err := c.ots.UpdateHoliday(*ch, pkg.Employee{})

	if err != nil {
		log.Debug(err)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	printHoliday(w, h)
	w.Flush()

	return nil
}

func (c *Client) DeleteHoliday(id uint) error {
	err := c.ots.DelHoliday(id, pkg.Employee{})

	if err != nil {
		log.Debug(err)
		return err
	}

	println("Holiday deleted")

	return nil
}
