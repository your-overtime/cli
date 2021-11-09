package client

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/cli/internal/utils"
)

func (c *Client) CalcCurrentOverview() error {
	o, err := c.ots.CalcOverview(pkg.Employee{}, time.Now())
	if err != nil {
		log.Debug(err)
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 16, 4, 3, '.', tabwriter.TabIndent)

	fmt.Fprintln(w, "\nOverview")
	fmt.Fprintf(w, "Current time\t: %s\n", utils.FormatTime(o.Date))
	fmt.Fprintf(w, "Week number\t: %d\n", o.WeekNumber)
	fmt.Fprintf(w, "Used holidays\t: %d\n", o.UsedHolidays)
	fmt.Fprintf(w, "Available holidays\t: %d\n", o.HolidaysStillAvailable)
	fmt.Fprintf(w, "Duration\t: Day\t Week\t Month \t Year\n")
	fmt.Fprintf(w, "ActiveTime\t: %s\t %s\t %s\t %s\n",
		formatMinutesToHoursAndMinutes(o.ActiveTimeThisDayInMinutes), formatMinutesToHoursAndMinutes(o.ActiveTimeThisWeekInMinutes),
		formatMinutes(o.ActiveTimeThisMonthInMinutes), formatMinutes(o.ActiveTimeThisYearInMinutes))
	fmt.Fprintf(w, "Overtime\t: %s\t %s\t %s\t %s\n",
		formatMinutesToHoursAndMinutes(o.OvertimeThisDayInMinutes), formatMinutesToHoursAndMinutes(o.OvertimeThisWeekInMinutes),
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
