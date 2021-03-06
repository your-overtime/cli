package client

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/your-overtime/cli/internal/conf"
)

func basicAuth(login string, password string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", login, password))))
}

func ConvertWWTStrToMins(wwtStr string) (uint, error) {
	answer := strings.Split(wwtStr, ":")
	var (
		h   int
		m   int
		err error
	)
	for _, p := range answer {
		if strings.HasSuffix(p, "h") {
			h, err = strconv.Atoi(strings.ReplaceAll(p, "h", ""))
			if err != nil {
				return 0, err
			}
		} else if strings.HasSuffix(p, "m") {
			m, err = strconv.Atoi(strings.ReplaceAll(p, "m", ""))
			if err != nil {
				return 0, err
			}
		}
	}
	return uint(h*60 + m), nil
}

func createUser(c *Client, adminToken string) error {

	return nil
}

func InitConf() error {
	var (
		token    string
		login    string
		password string
		url      string
	)

	qs := []*survey.Question{
		{
			Name:     "url",
			Prompt:   &survey.Input{Message: "What is your overtime api url?"},
			Validate: survey.Required,
		},
		{
			Name: "HasToken",
			Prompt: &survey.Confirm{
				Message: "Do you have an access token?",
			},
		},
	}
	answers1 := struct {
		URL      string
		HasToken bool
	}{}
	err := survey.Ask(qs, &answers1)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	url = answers1.URL

	if !answers1.HasToken {
		hasUser := false
		prompt := &survey.Confirm{
			Message: "Do you have login data?",
		}
		survey.AskOne(prompt, &hasUser)
		if !hasUser {
			adminToken := ""
			prompt := &survey.Input{
				Message: "Please insert the API admin token",
			}
			survey.AskOne(prompt, &adminToken)
			qs := []*survey.Question{
				{
					Name:     "name",
					Prompt:   &survey.Input{Message: "Please type your name"},
					Validate: survey.Required,
				},
				{
					Name: "surname",
					Prompt: &survey.Input{
						Message: "Please type your surname",
					},
					Validate: survey.Required,
				},
				{
					Name: "login",
					Prompt: &survey.Input{
						Message: "Please type your login",
					},
					Validate: survey.Required,
				},
				{
					Name: "password",
					Prompt: &survey.Password{
						Message: "Please type your password",
					},
					Validate: survey.Required,
				},
				{
					Name: "weekWorkingTimeInMinutes",
					Prompt: &survey.Input{
						Message: "Please type your working time per week [32h:30m]",
					},
					Validate: func(ans interface{}) error {
						v, err := ConvertWWTStrToMins(ans.(string))
						if err != nil {
							return err
						}
						if v == 0 {
							return strconv.ErrSyntax
						}
						return nil
					},
				},
				{
					Name: "numberWorkingDays",
					Prompt: &survey.Input{
						Message: "Please type your number of working days per week",
					},
					Validate: survey.Required,
				},
			}
			answers3 := struct {
				Name                     string
				Surname                  string
				Login                    string
				Password                 string
				WeekWorkingTimeInMinutes string
				NumWorkingDays           uint
			}{}
			err := survey.Ask(qs, &answers3)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			wwtim, err := ConvertWWTStrToMins(answers3.WeekWorkingTimeInMinutes)
			if err != nil {
				return err
			}
			c := Init(url, fmt.Sprintf("token %s", adminToken))
			err = c.AddEmployee(answers3.Name, answers3.Surname, answers3.Login, answers3.Password, wwtim, answers3.NumWorkingDays, adminToken)
			if err != nil {
				return err
			}
			login = answers3.Login
			password = answers3.Password
		} else {
			qs := []*survey.Question{
				{
					Name:     "login",
					Prompt:   &survey.Input{Message: "Please type your login"},
					Validate: survey.Required,
				},
				{
					Name: "password",
					Prompt: &survey.Password{
						Message: "Please type your password",
					},
					Validate: survey.Required,
				},
			}
			answers2 := struct {
				Login    string
				Password string
			}{}
			err := survey.Ask(qs, &answers2)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			login = answers2.Login
			password = answers2.Password
		}
		c := Init(url, basicAuth(login, password))
		t, err := c.CreateToken(fmt.Sprintf("CLI %s", time.Now()), false)
		if err != nil {
			return err
		}
		token = t.Token
	}
	defaultDesc := true
	prompt := &survey.Confirm{
		Message: "Do you like to set a default activity description?",
	}
	survey.AskOne(prompt, &defaultDesc)
	defaultActivityDesc := ""
	if defaultDesc {
		prompt := &survey.Input{
			Message: "Please type the default description",
		}
		survey.AskOne(prompt, &defaultActivityDesc)
	}
	c := conf.Config{
		Host:                url,
		Token:               fmt.Sprintf("token %s", token),
		DefaultActivityDesc: defaultActivityDesc,
	}
	err = conf.WriteConfig(c)
	if err != nil {
		return err
	}

	fmt.Println("The configuration is finshed and the \"otcli\" can be used now!!")

	return nil
}
