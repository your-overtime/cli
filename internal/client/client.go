package client

import (
	"fmt"
	"os"
	"text/tabwriter"

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

func (c *Client) StartActivity(desc string) error {
	a, err := c.ots.StartActivity(desc, pkg.Employee{})
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent|tabwriter.Debug)
	fmt.Println(w, fmt.Sprintf("Description:\t%s", a.Description))
	fmt.Println(w, fmt.Sprintf("Start:\t%s", a.Start))
	fmt.Println(w, fmt.Sprintf("End:\t%s", a.End))
	fmt.Println(w, fmt.Sprintf("ID:\t%d", a.ID))

	return nil
}

func (c *Client) CalcCurrentOverview() error {
	o, err := c.ots.CalcCurrentOverview(pkg.Employee{})
	fmt.Println(o, err)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent|tabwriter.Debug)

	fmt.Println(w, fmt.Sprintf("ActiveTime:\t%d", o.ActiveTimeThisWeek))
	fmt.Println(w, fmt.Sprintf("Overtime:\t%d", o.OvertimeInMinutes))

	if o.ActiveActivity != nil {
		a := o.ActiveActivity
		fmt.Println(w, fmt.Sprintf("Description:\t%s", a.Description))
		fmt.Println(w, fmt.Sprintf("Start:\t%s", a.Start))
		fmt.Println(w, fmt.Sprintf("End:\t%s", a.End))
		fmt.Println(w, fmt.Sprintf("ID:\t%d", a.ID))
	}

	return nil
}
