package client

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"git.goasum.de/jasper/overtime-cli/internal/conf"
	"github.com/AlecAivazis/survey/v2"
	"github.com/mdp/qrterminal/v3"
)

func (c *Client) LinkApp() error {
	co, err := conf.LoadConfig()
	if err != nil {
		return err
	}
	var token string
	newToken := true
	prompt := &survey.Confirm{
		Message: "Do you like to create a new token",
	}
	survey.AskOne(prompt, &newToken)
	if newToken {
		now := time.Now()
		t, err := c.CreateToken(fmt.Sprintf("APP %s", now))
		if err != nil {
			return err
		}

		token = fmt.Sprintf("token %s", t.Token)

		println("New token created")
	} else {
		token = co.Token
	}

	qc := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 4,
	}
	payload, err := json.Marshal(map[string]string{"url": co.Host, "authheader": token, "desc": co.DefaultActivityDesc})
	qrterminal.GenerateWithConfig(string(payload), qc)

	return nil
}
