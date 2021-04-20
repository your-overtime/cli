package client

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"

	"git.goasum.de/jasper/overtime-cli/internal/utils"
	"git.goasum.de/jasper/overtime/pkg"
)

type Client struct {
	ots pkg.OvertimeService
}

func Init(host string, token string) Client {
	return Client{
		ots: pkg.InitOvertimeClient(host, fmt.Sprintf("token %s", token)),
	}
}

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
		o, err := c.ots.CalcOverview(pkg.Employee{})
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
	for i, a := range as {
		fmt.Fprintf(w, "No\t: %d\n", i+1)
		printActivity(w, &a)
		fmt.Fprintln(w)
	}
	w.Flush()
	return nil
}

func printActivity(w *tabwriter.Writer, a *pkg.Activity) {
	fmt.Fprintf(w, "ID\t: %d\n", a.ID)
	fmt.Fprintf(w, "Description\t: %s\n", a.Description)
	fmt.Fprintf(w, "Start\t: %s\n", utils.FormatTime(*a.Start))
	if a.End != nil {
		fmt.Fprintf(w, "End\t: %s\n", utils.FormatTime(*a.End))
		diff := a.End.Sub(*a.Start)
		fmt.Fprintf(w, "Duration\t: %s\n", formatMinutes(int64(diff.Minutes())))
	}
}

func formatMinutes(t int64) string {
	ds, hs1 := math.Modf(float64(t) / (24 * 60))
	hs2, mf := math.Modf(hs1 * 24)
	if ds == 0 {
		return fmt.Sprintf("%02dh:%02dm", int(hs2), int(mf*60))
	}
	return fmt.Sprintf("%02dd:%02dh:%02dm", int(ds), int(hs2), int(mf*60))
}

func (c *Client) CalcCurrentOverview() error {
	o, err := c.ots.CalcOverview(pkg.Employee{})
	if err != nil {
		log.Debug(err)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 16, 4, 3, '.', tabwriter.TabIndent)

	fmt.Fprintln(w, "\nOverview")
	fmt.Fprintf(w, "Current time\t: %s\n", utils.FormatTime(o.Date))
	fmt.Fprintf(w, "Week number\t: %d\n", o.WeekNumber)
	fmt.Fprintf(w, "Duration\t: Day\t Week\t Month \t Year\n")
	fmt.Fprintf(w, "ActiveTime\t: %s\t %s\t %s\t %s\n",
		formatMinutes(o.ActiveTimeThisDayInMinutes), formatMinutes(o.ActiveTimeThisWeekInMinutes),
		formatMinutes(o.ActiveTimeThisMonthInMinutes), formatMinutes(o.ActiveTimeThisYearInMinutes))
	fmt.Fprintf(w, "Overtime\t: %s\t %s\t %s\t %s\n",
		formatMinutes(o.OvertimeThisDayInMinutes), formatMinutes(o.OvertimeThisWeekInMinutes),
		formatMinutes(o.OvertimeThisMonthInMinutes), formatMinutes(o.OvertimeThisYearInMinutes))

	if o.ActiveActivity != nil {
		a := o.ActiveActivity
		fmt.Fprintf(w, "\nRunning activity\n")
		fmt.Fprintf(w, "ID\t: %d\n", a.ID)
		fmt.Fprintf(w, "Description\t: %s\n", a.Description)
		fmt.Fprintf(w, "Start:\t: %s\n", utils.FormatTime(*a.Start))
	}

	w.Flush()

	return nil
}
