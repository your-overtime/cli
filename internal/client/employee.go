package client

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"text/tabwriter"

	"github.com/your-overtime/api/pkg"
)

func (c *Client) AddEmployee(name string, surname string, login string, pw string, wwt uint, nwwd uint, adminToken string) error {
	e, err := c.ots.SaveEmployee(pkg.Employee{
		User: &pkg.User{
			Name:     name,
			Surname:  surname,
			Login:    login,
			Password: pw,
		},
		WeekWorkingTimeInMinutes: wwt,
		NumWorkingDays:           nwwd,
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
	fmt.Fprintf(w, "NumWorkingDays\t: %d\n", e.NumWorkingDays)
	w.Flush()

	return nil
}

func (c *Client) DeleteEmployee(login string, adminToken string) error {
	err := c.ots.DeleteEmployee(login, adminToken)

	if err != nil {
		log.Debug(err)
		return err
	}

	fmt.Println("Employee deleted")

	return nil
}
