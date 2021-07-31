package client

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/cli/internal/utils"
)

func (c *Client) AddActivity(desc string, start *time.Time, end *time.Time) error {
	a, err := c.ots.AddActivity(pkg.Activity{
		Start:       start,
		End:         end,
		Description: desc,
	}, pkg.Employee{})

	if err != nil {
		log.Debug(err)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	printActivity(w, a)
	w.Flush()

	return nil
}

func (c *Client) UpdateActivity(desc string, start *time.Time, end *time.Time, id uint) error {
	ca, err := c.ots.GetActivity(id, pkg.Employee{})
	if err != nil {
		log.Debug(err)
		return err
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
	a, err := c.ots.UpdateActivity(*ca, pkg.Employee{})

	if err != nil {
		log.Debug(err)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	printActivity(w, a)
	w.Flush()

	return nil
}

func (c *Client) DeleteActivity(id uint) error {
	err := c.ots.DelActivity(id, pkg.Employee{})

	if err != nil {
		log.Debug(err)
		return err
	}

	println("Activity deleted")

	return nil
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
		err = c.AddActivity(desc, &start, &end)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (c *Client) StartActivity(desc string) error {
	a, err := c.ots.StartActivity(desc, pkg.Employee{})
	if err != nil {
		log.Debug(err)
		fmt.Println("\nA activity is currently running")
		o, err := c.ots.CalcOverview(pkg.Employee{}, time.Now())
		if err != nil {
			log.Debug(err)
			return err
		}
		if o.ActiveActivity == nil {
			panic(o)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
		printActivity(w, o.ActiveActivity)
		w.Flush()
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	fmt.Fprintf(w, "ID:\t%d\n", a.ID)
	fmt.Fprintf(w, "Description:\t%s\n", a.Description)
	fmt.Fprintf(w, "Start:\t%s\n", utils.FormatTime(*a.Start))
	w.Flush()

	return nil
}

func (c *Client) StopActivity() error {
	a, err := c.ots.StopRunningActivity(pkg.Employee{})
	if err != nil {
		log.Debug(err)
		return err
	}
	if a.Start == nil {
		fmt.Println("\nNo activity is running")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	printActivity(w, a)
	w.Flush()
	return nil
}

func (c *Client) GetActivities(start time.Time, end time.Time, asJSON bool) error {
	as, err := c.ots.GetActivities(start, end, pkg.Employee{})

	if asJSON {
		jsonData, err := json.MarshalIndent(as, "", " ")
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
	now := time.Now()
	for i, a := range as {
		fmt.Fprintf(w, "No\t: %d\n", i+1)
		printActivity(w, &a)
		fmt.Fprintln(w)
		if a.End != nil {
			mins += int(a.End.Sub(*a.Start).Minutes())
		} else {
			mins += int(now.Sub(*a.Start).Minutes())
		}
	}
	fmt.Fprintf(w, "Duration\t: %s\n", formatMinutes(int64(mins)))
	w.Flush()
	return nil
}
