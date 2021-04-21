package conf

import (
	"encoding/base64"
	"fmt"
	"time"

	"git.goasum.de/jasper/overtime-cli/internal/client"
	"git.goasum.de/jasper/overtime-cli/internal/utils"
)

func basicAuth(login string, password string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", login, password))))
}

func InitConf() error {
	url := utils.ReadTextFromStdin("overtime api url:\n")
	token := utils.ReadTextFromStdin("access token (if not exist leaf blank and press return):\n")
	if len(token) == 0 {
		login := utils.ReadTextFromStdin("login:\n")
		pw := utils.ReadTextFromStdin("password:\n")
		c := client.Init(url, basicAuth(login, pw))
		t, err := c.CreateToken(fmt.Sprintf("CLI %s", time.Now()))
		if err != nil {
			return err
		}
		token = t.Token
	}
	defaultActivityDesc := utils.ReadTextFromStdin("default activity description (empty if not needed):\n")
	c := Config{
		Host:                url,
		Token:               fmt.Sprintf("token %s", token),
		DefaultActivityDesc: defaultActivityDesc,
	}
	return WriteConfig(c)
}
