package client

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/AlecAivazis/survey/v2"
	"github.com/your-overtime/api/pkg"
)

func (c *Client) ChangeAccount(cn bool, cs bool, cl bool, cp bool, cwwt bool, cwwd bool) error {
	e, err := c.ots.GetAccount()
	if err != nil {
		return err
	}

	all := !(cn || cs || cl || cp || cwwt || cwwd)

	fields := map[string]interface{}{}

	if all || cn {
		updateStringValue(fields, fmt.Sprintf("Name: %s", e.Name), "Please type the new name", "Name")
	}

	if all || cs {
		updateStringValue(fields, fmt.Sprintf("Surname: %s", e.Surname), "Please type the new surname", "Surname")
	}

	if all || cl {
		updateStringValue(fields, fmt.Sprintf("Login: %s", e.Login), "Please type the new login", "Login")
	}

	if all || cp {
		changePassword(fields)
	}

	if all || cwwt {
		updateValue := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Week working time: %s change?", formatMinutesToHoursAndMinutes(int64(e.WeekWorkingTimeInMinutes))),
		}
		survey.AskOne(prompt, &updateValue)
		if updateValue {
			value := ""
			prompt := &survey.Input{
				Message: "Please type your working time per week [32h:30m]",
			}
			survey.AskOne(prompt, &value)
			v, err := ConvertWWTStrToMins(value)
			if err != nil {
				return err
			}
			fields["WeekWorkingTimeInMinutes"] = v
		}
	}

	if all || cwwd {
		updateValue := false
		prompt1 := &survey.Confirm{
			Message: fmt.Sprintf("Number of working days per week: %d change?", e.NumWorkingDays),
		}
		survey.AskOne(prompt1, &updateValue)
		if updateValue {
			numDays := ""
			prompt2 := &survey.Select{
				Message: "Number of working days per week:",
				Options: []string{"1", "2", "3", "4", "5", "6", "7"},
			}
			survey.AskOne(prompt2, &numDays)
			fields["NumWorkingDays"], err = strconv.Atoi(numDays)
			if err != nil {
				return err
			}
		}
	}

	em, err := c.ots.UpdateAccount(fields, pkg.Employee{})
	if err != nil {
		if err.Error() == "400 Bad Request" {
			fmt.Println("A account with the new login already exist")
			em, err = c.ots.GetAccount()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)

	fmt.Fprintf(w, "Surname\t: %s\n", em.Surname)
	fmt.Fprintf(w, "Name\t: %s\n", em.Name)
	fmt.Fprintf(w, "Login\t: %s\n", em.Login)
	fmt.Fprintf(w, "NumWorkingDays\t: %d\n", em.NumWorkingDays)
	fmt.Fprintf(w, "WeekWorkingTime\t: %s\n", formatMinutesToHoursAndMinutes(int64(em.WeekWorkingTimeInMinutes)))

	w.Flush()

	return nil
}

func changePassword(fields map[string]interface{}) {
	updatePassword := false
	prompt := &survey.Confirm{
		Message: "Change password?",
	}
	survey.AskOne(prompt, &updatePassword)
	if updatePassword {
		qs := []*survey.Question{
			{
				Name: "password0",
				Prompt: &survey.Password{
					Message: "Please type your password",
				},
				Validate: survey.Required,
			},
			{
				Name: "password1",
				Prompt: &survey.Password{
					Message: "Please type your password again",
				},
				Validate: survey.Required,
			},
		}
		answers := struct {
			Password0 string
			Password1 string
		}{}
		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
		}
		if answers.Password0 == answers.Password1 {
			fields["Password"] = answers.Password0
		} else {
			fmt.Println("The passwords not match")
			changePassword(fields)
		}
	}
}

func updateStringValue(fields map[string]interface{}, currentFieldValue string, msg string, fieldName string) {
	updateValue := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("%s change?", currentFieldValue),
	}
	survey.AskOne(prompt, &updateValue)
	if updateValue {
		value := ""
		prompt := &survey.Input{
			Message: msg,
		}
		survey.AskOne(prompt, &value)
		fields[fieldName] = value
	}
}

func (c *Client) GetAccount() error {
	em, err := c.ots.GetAccount()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)

	fmt.Fprintf(w, "Surname\t: %s\n", em.Surname)
	fmt.Fprintf(w, "Name\t: %s\n", em.Name)
	fmt.Fprintf(w, "Login\t: %s\n", em.Login)
	fmt.Fprintf(w, "NumWorkingDays\t: %d\n", em.NumWorkingDays)
	fmt.Fprintf(w, "WeekWorkingTime\t: %s\n", formatMinutesToHoursAndMinutes(int64(em.WeekWorkingTimeInMinutes)))

	w.Flush()

	return nil
}
