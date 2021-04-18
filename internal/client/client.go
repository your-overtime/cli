package client

import (
	"fmt"
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

func (c *Client) StartActivity(desc string) error {
	a, err := c.ots.StartActivity(desc, pkg.Employee{})
	if err != nil {
		log.Debug(err)
		fmt.Println("\nA activity is currently running")
		o, err := c.ots.CalcCurrentOverview(pkg.Employee{})
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

func (c *Client) GetActivities(start time.Time, end time.Time) error {
	as, err := c.ots.GetActivities(start, end, pkg.Employee{})
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
		hs, mf := math.Modf(diff.Hours())
		fmt.Fprintf(w, "Duration\t: %d:%d\n", int(hs), int(mf*60))
	}
}

func (c *Client) CalcCurrentOverview() error {
	o, err := c.ots.CalcCurrentOverview(pkg.Employee{})
	if err != nil {
		log.Debug(err)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)

	fmt.Fprintln(w, "\nOverview")
	fmt.Fprintf(w, "Current time\t: %s\n", utils.FormatTime(o.Date))
	fmt.Fprintf(w, "WeeK number\t: %d\n", o.WeekNumber)
	fmt.Fprintf(w, "ActiveTime\t: %d\n", o.ActiveTimeThisWeek)
	fmt.Fprintf(w, "Overtime\t: %d\n", o.OvertimeInMinutes)

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
