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

	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/cli/internal/utils"
)

type Client struct {
	ots pkg.OvertimeService
}

func Init(host string, authHeader string) Client {
	return Client{
		ots: pkg.InitOvertimeClient(host, authHeader),
	}
}

func (c *Client) AddEmployee(name string, surname string, login string, pw string, wwt uint, adminToken string) error {
	e, err := c.ots.SaveEmployee(pkg.Employee{
		User: &pkg.User{
			Name:     name,
			Surname:  surname,
			Login:    login,
			Password: pw,
		},
		WeekWorkingTimeInMinutes: wwt,
	}, adminToken)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	fmt.Fprintf(w, "ID\t: %d\n", e.ID)
	fmt.Fprintf(w, "Login\t: %s\n", e.Login)
	fmt.Fprintf(w, "Name\t: %s\n", e.Name)
	fmt.Fprintf(w, "Surname\t: %s\n", e.Surname)
	fmt.Fprintf(w, "WeekWorkingTimeInMinutes\t: %d\n", e.WeekWorkingTimeInMinutes)
	w.Flush()

	return nil
}

func (c *Client) DeleteEmployee(login string, adminToken string) error {
	err := c.ots.DeleteEmployee(login, adminToken)

	if err != nil {
		log.Debug(err)
		return err
	}

	println("Employee deleted")

	return nil
}

func (c *Client) CreateTokenViaCli(name string) error {
	t, err := c.CreateToken(name)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	fmt.Fprintf(w, "ID\t: %d\n", t.ID)
	fmt.Fprintf(w, "Token\t: %s\n", t.Token)
	fmt.Fprintf(w, "Name\t: %s\n\n", t.Name)
	w.Flush()

	return nil
}

func (c *Client) CreateToken(name string) (*pkg.Token, error) {
	t, err := c.ots.CreateToken(pkg.InputToken{Name: name}, pkg.Employee{})

	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return t, nil
}

func (c *Client) DeleteToken(id uint) error {
	err := c.ots.DeleteToken(id, pkg.Employee{})

	if err != nil {
		log.Debug(err)
		return err
	}

	println("Token deleted")

	return nil
}

func (c *Client) GetTokens() error {
	ts, err := c.ots.GetTokens(pkg.Employee{})

	if err != nil {
		log.Debug(err)
		return err
	}

	println()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	for _, t := range ts {
		fmt.Fprintf(w, "ID\t: %d\n", t.ID)
		fmt.Fprintf(w, "Token\t: %s\n", t.Token)
		fmt.Fprintf(w, "Name\t: %s\n\n", t.Name)
	}

	w.Flush()

	return nil
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

func printHoliday(w *tabwriter.Writer, a *pkg.Holiday) {
	fmt.Fprintf(w, "ID\t: %d\n", a.ID)
	fmt.Fprintf(w, "Description\t: %s\n", a.Description)
	fmt.Fprintf(w, "Start\t: %s\n", utils.FormatTime(a.Start))
	fmt.Fprintf(w, "End\t: %s\n", utils.FormatTime(a.End))
	diff := a.End.Sub(a.Start)
	fmt.Fprintf(w, "Duration\t: %s\n", formatMinutes(int64(diff.Minutes())))
}

func formatMinutes(t int64) string {
	ds, hs1 := math.Modf(float64(t) / (24 * 60))
	hs2, mf := math.Modf(hs1 * 24)
	if ds == 0 {
		return fmt.Sprintf("%02dh:%02dm", int(hs2), int(mf*60))
	}
	return fmt.Sprintf("%02dd:%02dh:%02dm", int(ds), int(hs2), int(mf*60))
}

func formatMinutesToHoursAndMinutes(t int64) string {
	hs2, mf := math.Modf(float64(t) / 60)
	return fmt.Sprintf("%02dh:%02dm", int(hs2), int(mf*60))
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

func (c *Client) AddHoliday(desc string, start time.Time, end time.Time, legalHoliday bool) error {
	h, err := c.ots.AddHoliday(pkg.Holiday{
		Start:        start,
		End:          end,
		Description:  desc,
		LegalHoliday: legalHoliday,
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

func (c *Client) UpdateHoliday(desc string, start *time.Time, end *time.Time, id uint, legalHoliday *bool) error {
	ch, err := c.ots.GetHoliday(id, pkg.Employee{})
	if err != nil {
		log.Debug(err)
		return err
	}
	if legalHoliday != nil {
		ch.LegalHoliday = *legalHoliday
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

func (c *Client) Export(since *time.Time, output string) error {
	fmt.Println("Export started")
	now := time.Now()
	if since == nil {
		yStart := time.Date(now.Year(), 0, 0, 0, 0, 0, 0, now.Location())
		since = &yStart
	}

	exportData := ExportData{}

	acs, err := c.ots.GetActivities(*since, now, pkg.Employee{})
	if err != nil {
		return err
	}
	exportData.Activities = acs

	hds, err := c.ots.GetHolidays(*since, now, pkg.Employee{})
	if err != nil {
		return err
	}
	exportData.Holidays = hds

	wds, err := c.ots.GetWorkDays(*since, now, pkg.Employee{})
	if err != nil {
		return err
	}
	exportData.Holidays = wds

	bytes, err := json.MarshalIndent(&exportData, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(output, bytes, os.ModeAppend)
	if err != nil {
		return err
	}
	fmt.Println("Export finished")
	return nil
}

func (c *Client) Import(input string) error {
	fmt.Println("Import started")
	bytes, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	exportData := ExportData{}
	err = json.Unmarshal(bytes, &exportData)
	if err != nil {
		return err
	}

	for _, a := range exportData.Activities {
		_, err := c.ots.AddActivity(a, pkg.Employee{})
		if err != nil {
			return err
		}
	}

	for _, h := range exportData.Holidays {
		_, err := c.ots.AddHoliday(h, pkg.Employee{})
		if err != nil {
			return err
		}
	}

	for _, h := range exportData.WorkDays {
		_, err := c.ots.AddWorkDay(h, pkg.Employee{})
		if err != nil {
			return err
		}
	}

	fmt.Println("Import finished")
	return nil
}
